package service

import (
	"golang.org/x/crypto/bcrypt"
)

// hash the password
func HashPassword(password string) (string, error) {

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
