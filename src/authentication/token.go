package authentication

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/brunoob35/TreeHouse-API/src/config"
	jwt "github.com/golang-jwt/jwt/v5"
)

// CustomClaims represents the custom JWT payload used by the application.
// It stores:
//
//   - Authorized: whether the token was generated as an authenticated token
//   - UserID: the ID of the authenticated user
//   - Permissions: the permission flags assigned to the user
//
// RegisteredClaims is embedded to support standard JWT fields such as
// expiration time, issue time, and not-before time.
type CustomClaims struct {
	Authorized  bool   `json:"authorized"`
	UserID      uint64 `json:"userId"`
	Permissions uint64 `json:"permissions"`
	jwt.RegisteredClaims
}

// GenerateToken creates and signs a JWT token for the authenticated user.
//
// The generated token stores:
//   - authorized = true
//   - userId = authenticated user's ID
//   - permissions = numeric permission flags for route authorization
//   - exp = expiration time
//   - iat = issued at time
//   - nbf = not before time
//
// This token is signed using the application's secret key and the HS256 algorithm.
func GenerateToken(userID uint64, permissions uint64) (string, error) {
	now := time.Now()

	claims := CustomClaims{
		Authorized:  true,
		UserID:      userID,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(config.SecretKey)
}

// ValidateToken validates whether the incoming request contains a valid bearer token.
//
// This function:
//
//   - extracts the token from the Authorization header
//   - parses it using the expected claims structure
//   - validates its signature and expiration
//   - ensures the expected signing method is being used
//   - ensures the token was marked as authorized
//
// It returns nil if the token is valid, otherwise it returns an error.
func ValidateToken(r *http.Request) error {
	tokenString, err := extractToken(r)
	if err != nil {
		return err
	}

	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		retrieveAuthKey,
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	if !claims.Authorized {
		return errors.New("unauthorized token")
	}

	return nil
}

// ExtractTokenData parses the bearer token and returns its typed claims.
//
// This is useful when the application needs to access information stored
// inside the token, such as:
//
//   - authenticated user ID
//   - permission flags
//
// It performs the same security checks as ValidateToken before returning the claims.
func ExtractTokenData(r *http.Request) (*CustomClaims, error) {
	tokenString, err := extractToken(r)
	if err != nil {
		return nil, err
	}

	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		retrieveAuthKey,
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if !claims.Authorized {
		return nil, errors.New("unauthorized token")
	}

	return claims, nil
}

// ExtractUserID retrieves only the authenticated user's ID from the token.
//
// This function is useful when handlers or services only need the user identity
// and do not need the full claims object.
func ExtractUserID(r *http.Request) (uint64, error) {
	claims, err := ExtractTokenData(r)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}

// ExtractPermissions retrieves the user's permission flags from the token.
//
// The returned value is a numeric bitmask that can later be used to validate
// whether the user is allowed to access a given route or perform a given action.
func ExtractPermissions(r *http.Request) (uint64, error) {
	claims, err := ExtractTokenData(r)
	if err != nil {
		return 0, err
	}

	return claims.Permissions, nil
}

// extractToken extracts the JWT token from the Authorization header.
//
// The expected header format is:
//
//	Authorization: Bearer <token>
//
// This function validates:
//   - whether the Authorization header exists
//   - whether it has the expected Bearer structure
//   - whether the token portion is not empty
//
// It is not exported because it is only used internally by the authentication package.
func extractToken(r *http.Request) (string, error) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid authorization header format")
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization type")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("token not found")
	}

	return token, nil
}

// retrieveAuthKey validates the token signing method and returns the secret key.
//
// While retrieving the authentication key, this function also ensures that
// the token was signed using an HMAC signing method, preventing tokens signed
// with unexpected algorithms from being accepted.
func retrieveAuthKey(token *jwt.Token) (interface{}, error) {
	if token.Method == nil {
		return nil, errors.New("invalid signing method")
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected token signing method: %v", token.Header["alg"])
	}

	return config.SecretKey, nil
}
