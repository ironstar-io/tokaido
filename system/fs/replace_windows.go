// +build !unix

package fs

import (
	"fmt"
	"os"
)

// Replace ...
func Replace(path string, body []byte) {
	if err := TouchByteArray(path+"-copy", body); err != nil {
		escapeHatch(path)
		return
	}

	if err := replaceFile(path); err != nil {
		escapeHatch(path)
		return
	}
}

func replaceFile(path string) error {
	// Copy the original file
	Copy(path, path+"-backup")

	// Copy `-copy` to be the new file
	Copy(path+"-copy", path)

	// Remove the backup file
	if err := os.Remove(path + "-backup"); err != nil {
		fmt.Println("There was an issue removing the backup file: ", err)
		return err
	}

	// Remove the copy file
	if err := os.Remove(path + "-copy"); err != nil {
		fmt.Println("There was an issue removing the copy file: ", err)
		return err
	}

	return nil
}

func escapeHatch(path string) {
	fmt.Println("Reverting...")

	// If the original no longer exists, try to restore the backup
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Copy `-copy` to be the new file
		Copy(path+"-backup", path)

		// Remove the backup file
		if err := os.Remove(path + "-backup"); err != nil {
			fmt.Println("There was an issue removing the backup file: ", err)
		}
	}

	if _, err := os.Stat(path + "-copy"); os.IsNotExist(err) {
		return
	}
	// Remove the copy file
	if err := os.Remove(path + "-copy"); err != nil {
		fmt.Println("There was an issue removing the copy file: ", err)
	}
}
