package version

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/ironstar-io/tokaido/system/version/goos"
	"github.com/ironstar-io/tokaido/utils"

	"github.com/blang/semver"
)

// SelfInstall - Check state then install runing Tokaido binary to PATH
func SelfInstall(forceInstall bool) {
	if forceInstall == true {
		// Tokaido not in PATH, confirm install of current bin
		confirmUpgrade := utils.ConfirmationPrompt("This command will install Tokaido "+Get().Version+". Would you like to continue?", "y")
		if confirmUpgrade == false {
			fmt.Println("Exiting...")
			return
		}

		installRunningBin()

		return
	}

	_, err := exec.LookPath("tok")
	if err != nil {
		// Tokaido not in PATH, confirm install of current bin
		confirmUpgrade := utils.ConfirmationPrompt("It looks like this is your first time running Tokaido.\n\nWould you like for us to install it?"+Get().Version, "y")
		if confirmUpgrade == false {
			fmt.Println("Exiting...")
			return
		}

		installRunningBin()

		return
	}

	// Tokaido already in PATH, display help message and exit
	fmt.Println("For help with Tokaido run `tok help` or take a look at our documentation at https://docs.tokaido.io")
}

// installRunningBin - Install runing Tokaido binary to PATH
func installRunningBin() {
	v := Get().Version
	cv := strings.Replace(v, "v", "", 0)
	cs, err := semver.Parse(cv)
	if err != nil {
		fmt.Println("Tokaido was unable to correctly parse the current version.")

		log.Fatal(err)
	}
	bv := cs.String()

	ip := GetInstallPath(bv)
	// Empty string if not installed, in which case, download and symlink
	if ip == "" {
		p, err := goos.CopyTokBinary(bv)
		if err != nil {
			fmt.Println("Tokaido wasn't able to install this version correctly.")

			log.Fatal(err)
		}

		ip = p
	}

	// Version is installed, just not active. Create a Symlink to finish
	err = CreateSymlink(ip)
	if err != nil {
		fmt.Printf("Tokaido was unable to create a symlink at %s", ip)

		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("Success! Tokaido version " + bv + " should now be avaliable as 'tok'")
}
