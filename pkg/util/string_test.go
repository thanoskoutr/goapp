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
		{3, "33E"},
		{5, "DE607"},
		{10, "9E1DEE6F7A"},
		{10, "20E65B801C"},
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

func TestRandHexString(t *testing.T) {
	tests := []struct {
		length   int
		expected string
	}{
		{0, ""},
		{3, "538"},
		{5, "538c7"},
		{10, "538c7f96b1"},
		{10, "538c7f96b1"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.length), func(t *testing.T) {
			actual := RandHexString(test.length)
			if actual != test.expected {
				t.Fatalf("length: %v, expected: %v, actual: %v", test.length, test.expected, actual)
			}
		})
	}
}

func BenchmarkRandString10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandString(10)
	}
}
func BenchmarkRandHexString10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandHexString(10)
	}
}

func BenchmarkRandString100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandString(100)
	}
}
func BenchmarkRandHexString100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandHexString(100)
	}
}

func BenchmarkRandString1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandString(1000)
	}
}
func BenchmarkRandHexString1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandHexString(1000)
	}
}

func BenchmarkRandString10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandString(10000)
	}
}
func BenchmarkRandHexString10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandHexString(10000)
	}
}
