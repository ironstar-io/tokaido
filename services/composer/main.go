package composer

import (
	"bufio"
	"errors"
	"strings"

	"github.com/ironstar-io/tokaido/system/ssh"
)

// FindPackageVersion - Find the version of an installed package
func FindPackageVersion(packageName string) (string, error) {
	d := ssh.ConnectCommand([]string{"cd", "/tokaido/site;", "composer", "show", packageName})
	if d == "" {
		return "", errors.New("Unable to find requested package: " + packageName)
	}

	matchString := "versions : * "
	scanner := bufio.NewScanner(strings.NewReader(d))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), matchString) {
			m := strings.Split(scanner.Text(), matchString)
			return m[1], nil
		}
	}

	return "", errors.New("Unable to find requested package version: " + packageName)
}

// RequirePackage ...
func RequirePackage(args []string) {
	c := append([]string{"cd", "/tokaido/site;", "composer", "require"}, args...)
	ssh.StreamConnectCommand(c)
}

// RemovePackage ...
func RemovePackage(args []string) {
	c := append([]string{"cd", "/tokaido/site;", "composer", "remove"}, args...)
	ssh.StreamConnectCommand(c)
}
