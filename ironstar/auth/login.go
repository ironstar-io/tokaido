package auth

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/ironstar/api"
	"github.com/ironstar-io/tokaido/system/fs"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

type AuthLoginBody struct {
	IDToken          string    `json:"id_token"`
	RedirectEndpoint string    `json:"redirect_endpoint"` // If this is set, user is MFA registered
	Expiry           time.Time `json:"expiry"`
}

var APILoginErrorMsg = "Ironstar API authentication failed"

func IronstarAPILogin(args []string, passwordFlag string) error {
	email, err := GetCLIEmail(args)
	if err != nil {
		return errors.Wrap(err, APILoginErrorMsg)
	}

	password, err := GetCLIPassword(passwordFlag)
	if err != nil {
		return errors.Wrap(err, APILoginErrorMsg)
	}

	req := &api.Request{
		AuthToken: "",
		Method:    "POST",
		Path:      "/auth/login",
		MapStringPayload: map[string]string{
			"email":    email,
			"password": password,
			"expiry":   time.Now().AddDate(0, 0, 14).UTC().Format(time.RFC3339),
		},
	}

	res, err := req.Send()
	if err != nil {
		return errors.Wrap(err, APILoginErrorMsg)
	}

	if res.StatusCode != 200 {
		return res.HandleFailure()
	}

	c, err := mfaCredentialCheck(res.Body, email)
	if err != nil {
		return errors.Wrap(err, APILoginErrorMsg)
	}

	err = UpdateCredentialsFile(Credentials{
		Login:     email,
		AuthToken: c.IDToken,
		Expiry:    c.Expiry,
	})
	if err != nil {
		return errors.Wrap(err, APILoginErrorMsg)
	}

	pr := fs.ProjectRoot()
	if pr != constants.ProjectRootNotFound {
		err = UpdateGlobalProjectLogin(pr, email)
		if err != nil {
			fmt.Println()
			color.Yellow("Authentication succeeded, but Tokaido was unable to update global credentials: ", err.Error())
		}
	}

	fmt.Println()
	color.Green("Ironstar API authentication succeeded!")
	fmt.Println()
	color.Green("User: ")
	fmt.Println(email)

	if c.Expiry.IsZero() {
		// Should always return expiry, but check anyway for safety
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
