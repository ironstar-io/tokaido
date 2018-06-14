package fs

import (
	"fmt"
	"os"
)

// Replace ...
func Replace(path string, body string) {
	if err := createCopy(path, body); err != nil {
		escapeHatch(path)
		return
	}

	if err := replaceFile(path); err != nil {
		escapeHatch(path)
		return
	}
}

func createCopy(path string, body string) error {
	file, err := os.Create(path + "-copy")
	if err != nil {
		fmt.Println("There was an issue creating a copy of this file: ", err)
		return err
	}

	if _, err := file.WriteString(body); err != nil {
		fmt.Println("There was an issue writing to a copy of this file: ", err)
		return err
	}

	defer file.Close()

	return nil
}

func replaceFile(path string) error {
	// Rename the original file
	if err := os.Rename(path, path+"-backup"); err != nil {
		fmt.Println("There was an issue renaming the original file: ", err)
		return err
	}

	// Rename `-copy` to be the new file
	if err := os.Rename(path+"-copy", path); err != nil {
		fmt.Println("There was an issue renaming the copy file: ", err)
		return err
	}

	// Remove the backup file
	if err := os.Remove(path + "-backup"); err != nil {
		fmt.Println("There was an issue removing the backup file: ", err)
		return err
	}

	return nil
}

func escapeHatch(path string) {
	fmt.Println("Reverting...")

	// If the original no longer exists, try to restore the backup
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Rename `-backup` to be the new file
		if err := os.Rename(path+"-backup", path); err != nil {
			fmt.Println("There was an issue restoring the original file: ", err)
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
