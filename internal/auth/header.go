package auth

import (
	"errors"
	"net/http"
)

func GetBearerToken(h http.Header) (string, error) {
  authHeader := h.Get("Authorization")

  if authHeader == "" {
    return "", errors.New("Authorization header is missing")
  }

  if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
    return "", errors.New("Authorization header is not a bearer token")
  }

  return authHeader[7:], nil
}

func GetAPIKey(h http.Header) (string, error) {
  authHeader := h.Get("Authorization")

  if authHeader == "" {
    return "", errors.New("Authorization header is missing")
  }

  if len(authHeader) < 7 || authHeader[:7] != "ApiKey " {
    return "", errors.New("Authorization header is not a ApiKey")
  }

  return authHeader[7:], nil
}
