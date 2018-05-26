package drupal

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

// DrushSSH SSH into the Drush container
func DrushSSH() {
	drushPort := docker.LocalPort("drush", "22")

	utils.StdoutCmd("ssh", "tok@localhost", "-p", drushPort, "-i", fs.HomeDir()+"/.ssh/tok_ssh.key", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-q")
}
