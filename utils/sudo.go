package utils

import (
	"bitbucket.org/ironstar/tokaido-cli/system/console"

	"log"
)

// CheckSudo - Check if a user has permissions to run sudo
// It's up to the calling function to prime the user to input their password,
// and to present an error and next-steps in the event that sudo access fails
func CheckSudo() error {
	_, err := CommandSubSplitOutput("sudo", "-n", "true")
	return err
}

// GainSudo - use `sudo true` to gain sudo privileges. Used if the user
// needs to specify a password to gain sudo, otherwise not necessary.
func GainSudo() {
	_, err := CommandSubSplitOutput("sudo", "true")
	if err != nil {
		console.Println(`
😢  Tokaido can't gain sudo privileges. Unable to proceed.
		`, "×")
		log.Fatal(err)
	}

}
