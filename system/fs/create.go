package fs

import (
	"fmt"
	"io/ioutil"
	"os"
)

// TouchByteArray ...
func TouchByteArray(path string, body []byte) error {
	if err := ioutil.WriteFile(path, body, 0644); err != nil {
		fmt.Println("There was an error creating a file: ", err)
		return err
	}

	return nil
}

// TouchEmpty ...
func TouchEmpty(path string) error {
	n, err := os.Create(path)
	if err != nil {
		fmt.Println("There was an error creating a file: ", err)
		return err
	}
	n.Close()

	return nil
}

// TouchOrReplace ...
func TouchOrReplace(path string, body []byte) {
	var _, errf = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(errf) {
		TouchByteArray(path, body)
		return
	}

	Replace(path, body)
}
