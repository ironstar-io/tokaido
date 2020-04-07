package auth

import (
	"path/filepath"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
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

	globals, err := ReadInGlobals()
	if err != nil {
		return errors.Wrap(err, SetCredentialsErrorMsg)
	}

	var newProjects []conf.Project
	var projMatch bool = false
	for _, proj := range globals.Projects {
		if proj.Path == pr {
			newProjects = append(newProjects, conf.Project{
				Name:  proj.Name,
				Path:  proj.Path,
				Login: email,
			})

			projMatch = true
			continue
		}

		newProjects = append(newProjects, proj)
	}

	if !projMatch {
		prs := strings.Split(pr, "/")
		if len(prs) == 0 {
			return errors.New("An unexpected error occurred, exiting...")
		}

		name := prs[len(prs)-1]
		newProjects = append(newProjects, conf.Project{
			Name:  name,
			Path:  pr,
			Login: email,
		})
	}

	globals.Projects = newProjects

	gp := filepath.Join(fs.HomeDir(), ".tok", "global.yml")
	newMarhsalled, err := yaml.Marshal(globals)
	if err != nil {
		return err
	}

	fs.Replace(gp, newMarhsalled)

	return nil
}
