package auth

import "testing"

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
