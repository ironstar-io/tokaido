package utils

import (
	"fmt"
	"os/user"
)

// HomeDir - Return the users' $HOME dir
func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Tokaido encountered a fatal error and had to stop")
		return FatalError(err)
	}
	return usr.HomeDir
}
