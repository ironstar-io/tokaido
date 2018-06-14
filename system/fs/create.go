package fs

import (
	"fmt"
	"io/ioutil"
)

// TouchByteArray ...
func TouchByteArray(path string, body []byte) error {
	if err := ioutil.WriteFile(path, body, 0644); err != nil {
		fmt.Println("There was an error creating a file: ", err)
		return err
	}

	return nil
}
