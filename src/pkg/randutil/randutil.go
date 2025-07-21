package randutil

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomStrings returns a cryptographically secure random string of length n.
func GenerateRandomStrings(n int) (string, error) {
	var sb strings.Builder
	sb.Grow(n) // Allocate memory upfront for better performance

	letterBytesLen := big.NewInt(int64(len(letterBytes)))

	for i := 0; i < n; i++ {
		index, err := rand.Int(rand.Reader, letterBytesLen)
		if err != nil {
			return "", fmt.Errorf("failed to generate secure random number: %w", err) // Provide more context for the error
		}
		sb.WriteByte(letterBytes[index.Int64()])
	}

	return sb.String(), nil
}
