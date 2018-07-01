package utils

import (
	"fmt"
	"log"
)

// if err != nil {
// 		fmt.Println(`
// 🔓  Tokaido requires privilege escalation to register the background sync service.
//     Your system will prompt you for your password so that we can add this service.
// 		`)
// 		utils.GainSudo()
// 	}

// 	_, subSplitErr := utils.CommandSubSplitOutput("sudo", "touch", "/tmp/tokaido-test-sudo")
// 	if subSplitErr != nil {
// 		fmt.Println("error")
// 	}

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
		fmt.Println(`
😢  Tokaido can't gain sudo privileges. Unable to proceed.
		`)
		log.Fatal(err)
	}

}
