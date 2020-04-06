package auth

import (
	"fmt"
)

var DefaultCredentialsErrorMsg = "Unable to set default credentials"

func SetDefaultCredentials(args []string) error {
	fmt.Println(args)
	fmt.Println("TODO - set default credentials")

	return nil
}
