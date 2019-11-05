package utils

import (
	"github.com/ironstar-io/tokaido/conf"

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
	if conf.GetConfig().Tokaido.Force == true || conf.GetConfig().Tokaido.Yes == true {
		return true
	}

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
