// +build !windows

package goos

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/console"

	"fmt"
)

// InitMessage - Display message post `up` success
func InitMessage() {
	fmt.Println(`
WELCOME TO TOKAIDO
==================

Your Drupal development environment is now up and running
	`)

	console.Println(`💻  Run "ssh `+conf.GetConfig().Tokaido.Project.Name+`.tok" to ssh into the environment`, "-")
	console.Println(`🌎  Run "tok open" to open the environment in your browser`, "-")
	console.Println(`👀  Run "tok exec" to run one-time commands like 'tok exec drush status'`, "-")
	console.Println(`🤔  Run "tok status" to check the status of your environment`, "-")
	fmt.Println(`
Check out https://tokaido.io/docs/using/ for tips to help you get the most out of your Tokaido environment
	`)
}
