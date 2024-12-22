package main

import (
	"testing"
)

func TestCleanBody(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"This is a kerfuffle opinion I need to share with the world", "This is a **** opinion I need to share with the world"},
		{"This is some random text with SharBert Sharbert! in this string", "This is some random text with **** Sharbert! in this string"},
		{"ForNax", "****"},
	}

	for _, test := range tests {
		output := cleanBody(test.input)
		if output != test.expected {
			t.Errorf("Test failed: input: %s, expected: %s, got: %s", test.input, test.expected, output)
		}
	}
}
