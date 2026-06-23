package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strings"

	"github.com/brunoob35/TreeHouse-API/src/config"
)

func piiKey() ([]byte, error) {
	if len(config.PIIKey) == 0 {
		return nil, errors.New("PII_ENCRYPTION_KEY não configurada")
	}

	sum := sha256.Sum256(config.PIIKey)
	return sum[:], nil
}

func EncryptPII(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", nil
	}

	key, err := piiKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(trimmed), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptPII(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", nil
	}

	key, err := piiKey()
	if err != nil {
		return "", err
	}

	raw, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(raw) < nonceSize {
		return "", errors.New("payload criptografado inválido")
	}

	nonce, ciphertext := raw[:nonceSize], raw[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func HashPII(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", nil
	}

	key, err := piiKey()
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, key)
	if _, err = mac.Write([]byte(trimmed)); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}
