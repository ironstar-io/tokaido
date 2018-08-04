package drupal

import (
	"github.com/ironstar-io/tokaido/system/ssh"
)

// ConfigureSSH - Configure `~/.ssh/config` and `~/.ssh/tok_config` to be ready for ssh access
func ConfigureSSH() {
	ssh.ProcessTokConfig()
}
