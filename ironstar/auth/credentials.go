package auth

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/fs"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Credentials struct {
	Login     string    `json:"login"`
	AuthToken string    `json:"auth_token"`
	Expiry    time.Time `json:"expiry"`
}

func SafeTouchGlobalConfigYAML(name string) error {
	cp := filepath.Join(fs.HomeDir(), ".tok", name+".yml")

	// Initialise the config file if it doesn't exist
	var _, errf = os.Stat(cp)
	if os.IsNotExist(errf) {
		// The global .tok path requires appropriate permissions
		gp := filepath.Dir(cp)
		if !fs.CheckExists(gp) {
			err := os.MkdirAll(gp, 0700)
			if err != nil {
				return err
			}
		}
		return fs.TouchEmpty(cp)
	}

	return nil
}

func ReadInCredentials() ([]Credentials, error) {
	cp := filepath.Join(fs.HomeDir(), ".tok", "credentials.yml")

	err := SafeTouchGlobalConfigYAML("credentials")
	if err != nil {
		return nil, err
	}

	cBytes, err := ioutil.ReadFile(cp)
	if err != nil {
		return nil, err
	}

	creds := []Credentials{}
	err = yaml.Unmarshal(cBytes, &creds)
	if err != nil {
		return nil, err
	}

	return creds, nil
}

func UpdateCredentialsFile(newCreds Credentials) error {
	cp := filepath.Join(fs.HomeDir(), ".tok", "credentials.yml")

	credSet, err := ReadInCredentials()
	if err != nil {
		return err
	}

	// Pull out matching login if it exists in the map
	var splicedCredSet []Credentials
	for _, cred := range credSet {
		if cred.Login != newCreds.Login {
			splicedCredSet = append(splicedCredSet, cred)
		}
	}

	// Replace/Add the new credentials to the struct slice
	newCredSet := append(splicedCredSet, newCreds)

	newMarhsalled, err := yaml.Marshal(newCredSet)
	if err != nil {
		return err
	}

	fs.Replace(cp, newMarhsalled)

	return nil
}

func UpdateGlobalProjectLogin(projectRoot, login string) error {
	globals, err := ReadInGlobals()
	if err != nil {
		return errors.Wrap(err, SetCredentialsErrorMsg)
	}

	var newProjects []conf.Project
	var projMatch bool = false
	for _, proj := range globals.Projects {
		if proj.Path == projectRoot {
			newProjects = append(newProjects, conf.Project{
				Name:  proj.Name,
				Path:  proj.Path,
				Login: login,
			})

			projMatch = true
			continue
		}

		newProjects = append(newProjects, proj)
	}

	if !projMatch {
		prs := strings.Split(projectRoot, "/")
		if len(prs) == 0 {
			return errors.New("An unexpected error occurred, exiting...")
		}

		name := prs[len(prs)-1]
		newProjects = append(newProjects, conf.Project{
			Name:  name,
			Path:  projectRoot,
			Login: login,
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

func UpdateGlobalDefaultLogin(login string) error {
	globals, err := ReadInGlobals()
	if err != nil {
		return errors.Wrap(err, SetCredentialsErrorMsg)
	}

	globals.DefaultLogin = login

	gp := filepath.Join(fs.HomeDir(), ".tok", "global.yml")
	newMarhsalled, err := yaml.Marshal(globals)
	if err != nil {
		return err
	}

	fs.Replace(gp, newMarhsalled)

	return nil
}
