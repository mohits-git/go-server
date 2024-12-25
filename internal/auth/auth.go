package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password cannot be empty")
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
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
