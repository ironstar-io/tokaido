package auth

import (
	"fmt"

	"github.com/ironstar-io/tokaido/ironstar/utils"

	"github.com/fatih/color"
)

func GetCLIPassword(passwordFlag string) (string, error) {
	var password string
	if passwordFlag == "" {
		input, err := utils.StdinSecret("Password: ")
		if err != nil {
			return "", err
		}
		password = input

		fmt.Println()
	} else {
		color.Yellow("Warning: Supplying a password via the command line flag is potentially insecure")

		password = passwordFlag
	}

	return password, nil
}
