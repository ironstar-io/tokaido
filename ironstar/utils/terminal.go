package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// StdinPrompt ...
func StdinPrompt(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(text), nil
}

// StdinSecret ...
func StdinSecret(prompt string) (string, error) {
	fmt.Print(prompt)
	byteSecret, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	secret := string(byteSecret)

	return strings.TrimSpace(secret), nil
}
