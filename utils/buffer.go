package utils

import (
	"bufio"
	"bytes"
	"strings"
)

// BufferInsert - Insert content into a buffer at a target string
func BufferInsert(buffer bytes.Buffer, target string, content string) bytes.Buffer {
	var newBuffer bytes.Buffer
	scanner := bufio.NewScanner(&buffer)
	for scanner.Scan() {
		newBuffer.Write([]byte(scanner.Text() + "\n"))
		if strings.Contains(scanner.Text(), target) {
			newBuffer.Write([]byte(content))
		}
	}

	return newBuffer
}
