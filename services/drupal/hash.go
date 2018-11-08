package drupal

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"strings"
)

// Hash is run from the command line and prints out a valid Drupal hash salt
func Hash() {
	res, err := GenerateRandomHashSalt()
	if err != nil {
		log.Fatalf("Error generating hash salt: %v", err)
	}
	fmt.Println("Drupal Hash Salt: ", res)
}

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomHashSalt returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomHashSalt() (string, error) {
	assertAvailablePRNG()
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	bytes, err := generateRandomBytes(55)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}

	s := base64.URLEncoding.EncodeToString(bytes)
	s = strings.Replace(s, "+", "-", -1)
	s = strings.Replace(s, "/", "-", -1)
	s = strings.Replace(s, "=", "", -1)

	return s, err
}

// Assert that a cryptographically secure PRNG is available.
// Panic otherwise.
func assertAvailablePRNG() {
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}
