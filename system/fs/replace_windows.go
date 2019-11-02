// +build windows

package fs

// TODO Windows

import (
	"fmt"
	"os"
	"io/ioutil"
)

// Replace ...
func Replace(path string, body []byte) {
	var _, err = os.Stat(path)

	if os.IsNotExist(err) {
		TouchByteArray(path, body)
		return
	}

	err = ioutil.WriteFile(path, body, 0)
	if err != nil {
		fmt.Println("There was an issue replacing file contents: ", err)
	}
}
