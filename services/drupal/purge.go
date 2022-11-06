package drupal

import (
	"fmt"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/ssh"
)

// Purge will empty the Varnish cache and run a drush cache-rebuild
func Purge() {
	purgeDrush()
}

func purgeDrush() {
	var cmd []string
	if conf.GetConfig().Drupal.Majorversion == "7" {
		cmd = []string{"drush", "cache-clear", "all", "-y"}
		fmt.Println("    Sending drush cache-clear all request...")
	} else {
		cmd = []string{"drush", "cache-rebuild", "-y"}
		fmt.Println("    Sending drush cache-rebuild request...")
	}

	resp := ssh.ConnectCommand(cmd)
	fmt.Println(resp)
}
