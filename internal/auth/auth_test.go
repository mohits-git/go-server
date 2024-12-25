package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		password string
		error    bool
	}{
		{"password", false},
		{"", true},
		{"too long password with more than 72 characters, random text random text random text.", true},
	}
	for _, test := range tests {
		_, err := HashPassword(test.password)
		if err != nil && !test.error {
			t.Error("HashPassword returned an error, expected nil")
		}
		if err == nil && test.error {
			t.Error("HashPassword did not return an error, expected one")
		}
	}
}

func TestComparePassword(t *testing.T) {
	hash, _ := HashPassword("password")
	tests := []struct {
		hash     string
		password string
		error    bool
	}{
		{hash, "password", false},
		{hash, "wrong password", true},
	}
	for _, test := range tests {
		err := ComparePassword(test.hash, test.password)
		if err != nil && !test.error {
			t.Error("ComparePassword returned an error, expected nil")
		}
		if err == nil && test.error {
			t.Error("ComparePassword did not return an error, expected one")
		}
	}
}

func TestMakeJWT(t *testing.T) {
	tokenSecret := "verysecret"

	tests := []struct {
		userId    uuid.UUID
		expiresIn time.Duration
		error     bool
	}{
		{func() uuid.UUID {
			userId, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
			if err != nil {
				return uuid.Nil
			}
			return userId
		}(), time.Hour, false},
	}

	for _, test := range tests {
		_, err := MakeJwt(test.userId, tokenSecret, test.expiresIn)
		if err != nil && !test.error {
			t.Error("MakeJwt returned an error, expected nil")
		}
		if err == nil && test.error {
			t.Error("MakeJwt did not return an error, expected one")
		}
	}
}

func TestValidateJWT(t *testing.T) {
	tokenSecret := "verysecretpass"

	tests := []struct {
		tokenString string
		error       bool
	}{
		{func() string {
			userId := uuid.New()
			jwtStr, _ := MakeJwt(userId, tokenSecret, time.Hour)
			return jwtStr
		}(), false},
		{func() string {
			userId := uuid.New()
			jwtStr, _ := MakeJwt(userId, "someothersecret", time.Hour)
			return jwtStr
		}(), true},
		{func() string {
			userId := uuid.New()
			jwtStr, _ := MakeJwt(userId, tokenSecret, 0)
			return jwtStr
		}(), true},
	}

	for _, test := range tests {
		_, err := ValidateJWT(test.tokenString, tokenSecret)
		if err != nil && !test.error {
			t.Error("ValidateJWT returned an error, expected nil.\n", err)
		}

		if err == nil && test.error {
			t.Error("ValidateJWT did not return an error, expected one.\n", err)
		}

    // This is a blocking test, it will wait for the jwt to expire
		// timer := time.NewTimer(time.Second)
		// <-timer.C
		// timer.Stop()
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		header   http.Header
		expected string
		error    bool
	}{
		{
			http.Header{"Authorization": []string{"Bearer token"}},
			"token",
			false,
		},
		{
			http.Header{"X-API": []string{"Bearer token"}},
			"",
			true,
		},
		{
			http.Header{"Authorization": []string{"tokens "}},
			"",
			true,
		},
		{
			http.Header{"Authorization": []string{""}},
			"",
			true,
		},
	}

	for _, test := range tests {
		token, err := GetBearerToken(test.header)
		if token != test.expected && !test.error {
			t.Error("GetBearerToken returned unexpected token, expected ", test.expected, " got ", token)
		}
		if err != nil && !test.error {
			t.Error("GetBearerToken returned an error, expected nil")
		}
		if err == nil && test.error {
			t.Error("GetBearerToken did not return an error, expected one")
		}
	}
}
