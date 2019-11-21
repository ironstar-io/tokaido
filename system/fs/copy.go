package fs

import (
	"fmt"
	"io"
	"os"
)

// Copy ...
func Copy(source string, destination string) {
	from, err := os.Open(source)
	if err != nil {
		fmt.Println("Unexpected error copying from source: ", err.Error())
		os.Exit(1)
	}
	defer from.Close()

	to, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Unexpected error opening destination file handler: ", err.Error())
		os.Exit(1)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		fmt.Println("Unexpected error copying ["+source+"] to ["+destination+"]: ", err.Error())
		os.Exit(1)
	}
}
