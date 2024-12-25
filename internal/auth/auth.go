package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password cannot be empty")
	}
	if len(password) > 72 {
		return "", errors.New("password length cannot exceed 72 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(hash, password string) error {
	if len(password) == 0 {
		return errors.New("password cannot be empty")
	}
	if len(password) > 72 {
		return errors.New("password length cannot exceed 72 characters")
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
