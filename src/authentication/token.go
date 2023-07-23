package authentication

import (
	"errors"
	"fmt"
	"github.com/brunoob35/TreeHouse-API/src/config"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

func GenerateToken(userID uint64) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissions["userId"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	return token.SignedString(config.SecretKey)

}

// ValidateToken receives the token alone, not the object containing it.
func ValidateToken(r *http.Request) error {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, retriveAuthKey)

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("invalid Token")
}

// extractToken is not going to be exported.
// It is going to validate if token is inside a bearer object by splitting it first.
// Then it extracts the token for later validation. This process prevents anyone of adding random values as a token.
func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

// retriveAuthKey As we retrieve the authentication ke, we also validate is the receiving token has teh expected signing family.
// Meaning it has not been altered or forged by any means.
func retriveAuthKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Token signing method unexpected! %v", token.Header["alg"])
	}
	return config.SecretKey, nil
}
