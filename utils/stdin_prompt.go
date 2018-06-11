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

// ConfirmationPrompt - The 'weighting' param should be one of [ "y", "n" ].
func ConfirmationPrompt(prompt string, weighting string) bool {
	response := StdinPrompt(prompt + weightedString(weighting))
	cutResponse := strings.ToLower(string([]rune(response)[0]))

	if weighting == "n" {
		if cutResponse == "y" {
			return true
		}

		return false
	}

	if cutResponse == "n" {
		return false
	}

	return true
}

func weightedString(weighting string) string {
	if weighting == "y" {
		return " (Y/n): "
	}

	return " (y/N): "
}
