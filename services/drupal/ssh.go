package drupal

import (
	// "bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/system/ssh"
	// "bitbucket.org/ironstar/tokaido-cli/utils"
)

// DrushSSH SSH into the Drush container
func DrushSSH() {
	//drushPort := docker.LocalPort("drush", "22")
	ssh.Client()
}
