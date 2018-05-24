package utils

import (
	"log"
	"os"
)

// WorkDir - Return the current working directory
func WorkDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return dir
}
