package fs

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Contains ...
func Contains(path string, keyword string) bool {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return true
	}

	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)
	foundKeyword := false
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), keyword) {
			foundKeyword = true
		}
	}

	return foundKeyword
}
