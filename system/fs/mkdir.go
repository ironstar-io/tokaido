package fs

import (
	"os"
)

// Mkdir - Check existence of folder and create if it doesn't exist
func Mkdir(path string) {
	_, err := os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}
