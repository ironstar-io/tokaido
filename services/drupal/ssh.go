package drupal

import (
	"bitbucket.org/ironstar/tokaido-cli/system/ssh"
)

// ConfigureSSH - Configure `~/.ssh/config` and `~/.ssh/tok_config` to be ready for ssh access
func ConfigureSSH() {
	ssh.ProcessTokConfig()
}
