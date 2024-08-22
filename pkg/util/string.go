package util

import (
	"encoding/hex"
	"math/rand"
)

const seed int64 = 42

var randx = rand.NewSource(seed)

// RandString returns a random hex string of length n.
func RandString(n int) string {
	const letterBytes = "abcdef0123456789"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, randx.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randx.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// RandHexString is an alternative implementation that returns a random
// hex string of length n.
func RandHexString(n int) string {
	// Create local RNG
	randx := rand.New(rand.NewSource(seed))

	// Calculate number of bytes needed
	fullBytes := n / 2
	extraByte := n % 2

	// Generate random bytes
	bytes := make([]byte, fullBytes+extraByte)
	randx.Read(bytes)

	// Convert bytes to hex string
	hexStr := hex.EncodeToString(bytes)

	// If n is odd, trim last char
	if extraByte == 1 {
		hexStr = hexStr[:n]
	}
	return hexStr
}
