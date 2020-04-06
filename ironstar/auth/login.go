package auth

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ironstar-io/tokaido/ironstar/api"

	"github.com/fatih/color"
)

type AuthLoginBody struct {
	IDToken          string    `json:"id_token"`
	RedirectEndpoint string    `json:"redirect_endpoint"`
	Expiry           time.Time `json:"expiry"`
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
		"expiry":   time.Now().AddDate(0, 0, 14).UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return api.HandleFailure(res)
	}

	c, err := mfaCredentialCheck(res.Body, email)
	if err != nil {
		return err
	}

	err = UpdateCredentialsFile(Credentials{
		Login:     email,
		AuthToken: c.IDToken,
		Expiry:    c.Expiry,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	color.Green("Ironstar API authentication succeeded!")
	fmt.Println()
	color.Green("User: ")
	fmt.Println(email)

	if c.Expiry.IsZero() {
		return nil
	}

	fmt.Println()
	color.Green("Expiry: ")

	expDiff := strconv.Itoa(int(math.RoundToEven(c.Expiry.Sub(time.Now().UTC()).Hours() / 24)))
	fmt.Println(c.Expiry.String() + " (" + expDiff + " days)")

	return nil
}

func mfaCredentialCheck(body []byte, email string) (*AuthLoginBody, error) {
	b := &AuthLoginBody{}
	err := json.Unmarshal(body, b)
	if err != nil {
		return nil, err
	}

	// If this is set, user is an MFA user
	if b.RedirectEndpoint != "" {
		c, err := ValidateMFAPasscode(b)
		if err != nil {
			return nil, err
		}

		return c, nil
	}

	return b, nil
}
