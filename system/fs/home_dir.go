package fs

import (
	"fmt"
	"log"
	"os/user"
)

// HomeDir - Return the users' $HOME dir
func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Tokaido encountered a fatal error and had to stop")
		log.Fatal(err)
	}

	return usr.HomeDir
}
