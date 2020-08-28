// +build !windows

package goos

import (
	"fmt"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/logrusorgru/aurora"
)

// InitMessage - Display message post `up` success
func InitMessage() {
	n := conf.GetConfig().Tokaido.Project.Name
	fmt.Println()
	fmt.Println(aurora.Green("Your Drupal development environment is now up and running"))
	fmt.Println(aurora.Green(aurora.Sprintf("You can find it at %s", aurora.Bold("https://"+n+".local.tokaido.io:5154/"))))
	fmt.Println()

	fmt.Printf("💻  Run '%s' to ssh into the environment\n", aurora.Bold("ssh "+n+".tok"))
	fmt.Printf("🌎  Run '%s' to open the environment in your browser\n", aurora.Bold("tok open"))
	fmt.Printf("👀  Run '%s' to run one-time commands like '%s'\n", aurora.Bold("tok exec"), aurora.Bold("tok exec drush status"))
	fmt.Printf("🤔  Run '%s' to check the status of your environment\n", aurora.Bold("tok status"))
	fmt.Println()
	fmt.Printf("Come join us in the %s channel in the Drupal Slack community!\n", aurora.Bold("#Tokaido"))
	fmt.Printf("or visit %s to check out the Tokaido Documentation\n\n", aurora.Bold("https://docs.tokaido.io"))
}
