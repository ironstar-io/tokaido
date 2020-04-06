package auth

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ironstar-io/tokaido/ironstar/api"
	"github.com/ironstar-io/tokaido/ironstar/utils"

	"github.com/fatih/color"
)

func GetCLIMFAPasscode() (string, error) {
	passcode, err := utils.StdinPrompt("MFA Passcode: ")
	if err != nil {
		return "", err
	}

	return passcode, nil
}

func ValidateMFAPasscode(logResBody *AuthLoginBody) (*AuthLoginBody, error) {
	passcode, err := GetCLIMFAPasscode()
	if err != nil {
		return nil, err
	}

	res, err := api.Req(logResBody.IDToken, "POST", "/auth/mfa/validate", map[string]string{
		"passcode": passcode,
	})
	if err != nil {
		return nil, err
	}

	m := &AuthLoginBody{}
	err = json.Unmarshal(res.Body, m)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		f := &api.FailureBody{}
		err = json.Unmarshal(res.Body, f)
		if err != nil {
			return nil, err
		}

		fmt.Println()
		color.Red("Ironstar API authentication failed!")
		fmt.Println()
		fmt.Printf("Status Code: %+v\n", res.StatusCode)
		fmt.Println("Ironstar Code: " + f.Code)
		fmt.Println(f.Message)

		return nil, errors.New("Unsuccessful!")
	}

	return m, nil
}
