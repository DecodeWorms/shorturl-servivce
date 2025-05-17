package urlshortener

import (
	"math/rand"
	"time"
)

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateShortCode returns a random 7-character alphanumeric string
func GenerateShortCode() string {
	codeLength := 7
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}
