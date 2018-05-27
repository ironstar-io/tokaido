package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadStdin ...
func ReadStdin(prompt string) string {
	fmt.Println(prompt)

	val, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	val = strings.TrimSpace(val)

	return val
}
