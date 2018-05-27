package drupal

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"
	"bitbucket.org/ironstar/tokaido-cli/system/ssh"

	"log"
)

// DrushSSH SSH into the Drush container
func DrushSSH() {
	drushPort := docker.LocalPort("drush", "22")

	client, err := ssh.DialWithKey("localhost:"+drushPort, "tok", fs.HomeDir()+"/.ssh/tok_ssh.key")
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
}
