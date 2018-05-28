package drupal

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/system/ssh"
	"fmt"
)

// ConfigureSSH - Configure `~/.ssh/config` and `~/.ssh/tok_config` to be ready for ssh access
func ConfigureSSH() {
	config := conf.GetConfig()

	ssh.ProcessTokConfig()

	fmt.Printf("\nTo access Drush via SSH run `ssh %s.tok`\n\n", config.Project)
}
