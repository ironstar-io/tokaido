package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// StdinPrompt ...
func StdinPrompt(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')

	return text
}

// ConfirmationPrompt ...
func ConfirmationPrompt(prompt string) bool {
	response := StdinPrompt(prompt + " (Y/n): ")
	cutResponse := strings.ToLower(string([]rune(response)[0]))

	if cutResponse == "n" {
		return false
	}

	return true
}
