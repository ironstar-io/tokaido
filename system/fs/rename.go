package fs

import (
	"fmt"
	"os"
)

// Rename - os Rename function
func Rename(orginal, new string) error {
	if err := os.Rename(orginal, new); err != nil {
		fmt.Println("There was an issue renaming the copy file: ", err)
		return err
	}

	return nil
}
