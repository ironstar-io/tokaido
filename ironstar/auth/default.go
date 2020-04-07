package auth

import (
	"github.com/pkg/errors"
)

var DefaultCredentialsErrorMsg = "Unable to set default credentials"

func SetDefaultCredentials(args []string) error {
	email, err := GetCLIEmail(args)
	if err != nil {
		return errors.Wrap(err, SetCredentialsErrorMsg)
	}

	return UpdateGlobalDefaultLogin(email)
}
