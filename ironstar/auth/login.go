package auth

import (
	"errors"
	"fmt"

	"github.com/ironstar-io/tokaido/ironstar/api"
)

func IronstarAPILogin(args []string, passwordFlag string) error {
	email, err := GetCLIEmail(args)
	if err != nil {
		return err
	}

	password, err := GetCLIPassword(passwordFlag)
	if err != nil {
		return err
	}

	res, err := api.Req("", "POST", "/auth/login", map[string]string{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return err
	}

	fmt.Printf("%+v", res)

	if res.StatusCode != 200 {
		return errors.New("Unsuccessful!")
	}

	return nil
}
