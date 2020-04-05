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
	} else {
		color.Yellow("Warning: Supplying the password via the command line is potentially insecure")

		password = passwordFlag
	}

	fmt.Println()

	return password, nil
}
