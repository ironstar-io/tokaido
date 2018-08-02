// +build !unix

package goos

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/system/console"

	"fmt"
)

// InitMessage - Display message post `up` success
func InitMessage() {
	fmt.Println(`
WELCOME TO TOKAIDO
==================

Your Drupal development environment is now up and running
	`)

	console.Println(`🔁  Run "tok watch" in a new terminal window to keep your local and container files in sync`, "-")
	console.Println(`💻  Run "ssh `+conf.GetConfig().Tokaido.Project.Name+`.tok" to access the Drush container`, "-")
	console.Println(`🌎  Run "tok open" to open the environment in your browser`, "-")
	console.Println(`🚀  Run "tok status" check the status of your environment`, "-")
	fmt.Println(`
Check out https://docs.tokaido.io/environments for tips on managing your Tokaido environment
	`)
}
