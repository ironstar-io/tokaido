package fs

import (
	"fmt"
	"os"
)

// RemoveAll - `rm -rf` on a dir
func RemoveAll(path string) {
	_, err := os.Stat(path)

	if !os.IsNotExist(err) {
		if err := os.RemoveAll(path); err != nil {
			fmt.Println("There was an issue removing the directory: ", err)
		}
	}
}
