package auth

import (
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"

	"github.com/pkg/errors"
)

var SetCredentialsErrorMsg = "Unable to set credentials"
var NoProjectFoundErrorMsg = "This command can only be run from a Tokaido project directory"
var NoCredentialMatch = errors.New("There are no credentials available for the supplied email")

func SetProjectCredentials(args []string) error {
	pr := fs.ProjectRoot()
	if pr == constants.ProjectRootNotFound {
		return errors.New(NoProjectFoundErrorMsg)
	}

	email, err := GetCLIEmail(args)
	if err != nil {
		return errors.Wrap(err, SetCredentialsErrorMsg)
	}

	creds, err := ReadInCredentials()
	if err != nil {
		return errors.Wrap(err, SetCredentialsErrorMsg)
	}

	var credMatch Credentials
	for _, cred := range creds {
		if cred.Login == email {
			credMatch = cred
		}
	}

	if (Credentials{}) == credMatch {
		return errors.Wrap(NoCredentialMatch, SetCredentialsErrorMsg)
	}

	return UpdateGlobalProjectLogin(pr, email)
}
