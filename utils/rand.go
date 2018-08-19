package utils

import (
	"crypto/rand"
	"fmt"
)

// RandomString ...
func RandomString() (string, error) {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%X", b), nil
}
