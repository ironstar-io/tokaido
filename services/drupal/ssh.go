package drupal

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/system/ssh"
	"fmt"
)

// DrushSSH SSH into the Drush container
func DrushSSH() {
	config := conf.GetConfig()

	ssh.ProcessTokConfig()

	fmt.Printf("\nTo access Drush via SSH run `ssh %s.tok`\n\n", config.Project)
}
