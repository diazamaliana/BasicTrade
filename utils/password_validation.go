
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword verifies the admin's password against a provided plaintext password.
func VerifyPassword(password, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(plainPassword))
}
