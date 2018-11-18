package drush

import (
	"log"
	"strings"

	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/ssh"
)

// DockerUp - Lift the drush container
func DockerUp() {
	docker.ComposeStdout("up", "-d", "drush")
}

// GetAdminSignOnPath - Use `drush uli` to generate a one-time login and extract the path
func GetAdminSignOnPath() string {
	u := ssh.ConnectCommand([]string{"drush", "uli"})
	if u == "" {
		log.Fatal("Unable to generate a single use admin URL with `drush uli`. Exiting...")
	}

	s := strings.Split(u, "http://default")
	if len(s) <= 1 {
		log.Fatal("Unable to parse the resulting path from the `drush uli` command. Exiting...")
	}

	return strings.Replace(s[1], "%0A", "", -1)
}

// LocalPort - Return the local port for drush
func LocalPort() string {
	p := docker.LocalPort("drush", "5000")
	if p == "" {
		log.Fatal("Drush container doesn't appear to be running!")
	}

	return p
}
