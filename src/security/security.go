package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

//go get golang.org/x/crypto/bcrypt

// Hash receives a string and hashes it
// Its necessary to decode the password to a slice of bytes and to set a cost for the operation.
// Bcrypt has a default value, this is a balanced value it won't consume too much computing resources.
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// ValidatePassword compares a string password from the request to the hash saved in DB, validating the operation
// and returning if equals.
func ValidatePassword(hashedPassword, stringPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(stringPassword))
}

// GenerateSecureToken Gera um novo token para redefinição de senha
func GenerateSecureToken(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashToken Realiza o hash do token antes de salvar no DB
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
