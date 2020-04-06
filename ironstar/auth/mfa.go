package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ironstar-io/tokaido/ironstar/api"
	"github.com/ironstar-io/tokaido/ironstar/utils"

	"github.com/fatih/color"
)

func GetCLIMFAPasscode() (string, error) {
	passcode, err := utils.StdinSecret("MFA Passcode: ")
	if err != nil {
		return "", err
	}

	if len(passcode) != 6 {
		fmt.Println()
		color.Red("Ironstar API authentication failed!")
		fmt.Println()
		fmt.Println("MFA passcode length must be 6")

		return "", errors.New("Passcode length must be 6")
	}

	return passcode, nil
}

func ValidateMFAPasscode(logResBody *AuthLoginBody) (*AuthLoginBody, error) {
	passcode, err := GetCLIMFAPasscode()
	if err != nil {
		return nil, err
	}

	req := &api.Request{
		AuthToken: logResBody.IDToken,
		Method:    "POST",
		Path:      "/auth/mfa/validate",
		MapStringPayload: map[string]string{
			"passcode": passcode,
			"expiry":   time.Now().AddDate(0, 0, 14).UTC().Format(time.RFC3339),
		},
	}

	res, err := req.Send()
	if err != nil {
		return nil, err
	}

	m := &AuthLoginBody{}
	err = json.Unmarshal(res.Body, m)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, res.HandleFailure()
	}

	return m, nil
}
