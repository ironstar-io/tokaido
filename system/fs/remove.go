package fs

import (
	"fmt"
	"os"
)

// Remove - `rm` on a file/dir
func Remove(path string) {
	_, err := os.Stat(path)

	if !os.IsNotExist(err) {
		if err := os.Remove(path); err != nil {
			fmt.Println("There was an issue removing the file: ", err)
		}
	}
}
