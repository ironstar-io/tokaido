package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ironstar-io/tokaido/ironstar/api"

	"github.com/fatih/color"
)

type AuthLoginBody struct {
	IDToken string    `json:"id_token"`
	Expiry  time.Time `json:"expiry"`
}

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

	if res.StatusCode != 200 {
		f := &api.FailureBody{}
		err = json.Unmarshal(res.Body, f)
		if err != nil {
			return err
		}

		fmt.Println()
		color.Red("Ironstar API authentication failed!")
		fmt.Println()
		fmt.Printf("Status Code: %+v\n", res.StatusCode)
		fmt.Println("Ironstar Code: " + f.Code)
		fmt.Println(f.Message)

		return errors.New("Unsuccessful!")
	}

	b := &AuthLoginBody{}
	err = json.Unmarshal(res.Body, b)
	if err != nil {
		return err
	}

	err = UpdateCredentialsFile(Credentials{
		Login:     email,
		AuthToken: b.IDToken,
		Expiry:    b.Expiry,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	color.Green("Ironstar API authentication succeeded!")
	fmt.Println()
	color.Green("Authentication Token: ")
	fmt.Println(b.IDToken)

	return nil
}
