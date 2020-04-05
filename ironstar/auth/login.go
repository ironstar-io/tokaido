package auth

import (
	"errors"
	"fmt"

	"github.com/ironstar-io/tokaido/ironstar/api"
)

func IronstarAPILogin(args []string) error {
	email, err := GetUserEmail(args)
	if err != nil {
		return err
	}

	err = ValidateEmail(email)
	if err != nil {
		return err
	}

	res, err := api.Req("", "POST", "/auth/login", map[string]string{
		"email":    email,
		"password": "hello",
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
