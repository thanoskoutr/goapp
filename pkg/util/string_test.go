package util

import (
	"fmt"
	"testing"
)

func TestRandString(t *testing.T) {
	tests := []struct {
		length   int
		expected string
	}{
		{0, ""},
		{3, "33e"},
		{5, "de607"},
		{10, "9e1dee6f7a"},
		{10, "20e65b801c"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.length), func(t *testing.T) {
			actual := RandString(test.length)
			if actual != test.expected {
				t.Fatalf("length: %v, expected: %v, actual: %v", test.length, test.expected, actual)
			}
		})
	}
}
