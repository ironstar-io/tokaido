// +build !unix

package goos

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"

	"fmt"
)

// InitMessage - Display message post `up` success
func InitMessage() {
	fmt.Println(`
WELCOME TO TOKAIDO
==================

Your Drupal development environment is now up and running
	`)

	fmt.Printf("💻  Run \"ssh %s.tok\" to access the Drush container\n", conf.GetConfig().Project)
	fmt.Println(`🌎  Run "tok open" to open the environment in your browser`)
	fmt.Println(`🤔  Run "tok status" check the status of your environment`)
	fmt.Println(`
Check out https://docs.tokaido.io/environments for tips on managing your Tokaido environment
	`)
}
