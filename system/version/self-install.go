package version

import (
	"fmt"
	"log"
	"strings"

	"github.com/ironstar-io/tokaido/utils"

	"github.com/blang/semver"
)

// Upgrade - Check if update available and auto-upgrade the user
func Upgrade() {
	v := Get().Version
	cv := strings.Replace(v, "v", "", 0)
	cs, err := semver.Parse(cv)
	if err != nil {
		fmt.Println("Tokaido was unable to correctly parse the current version.")

		log.Fatal(err)
	}

	// Checks if the latest version is Greater Than the current version
	confirmUpgrade := utils.ConfirmationPrompt("It looks like this is your first time running Tokaido.\n\nWould you like for us to install `tok` on your PATH?", "y")
	if confirmUpgrade == false {
		fmt.Println("Exiting...")
		return
	}

	ip := GetInstallPath(cs.String())
	// Empty string if not installed, in which case, download and symlink
	if ip == "" {
		// Download latest release from GH
		p, err := DownloadAndInstall(cs.String())
		if err != nil {
			fmt.Println("Tokaido wasn't able to upgrade you to the latest version.")

			log.Fatal(err)
		}

		ip = p
	}

	// Latest version is installed, just not active. Create a Symlink to finish
	err = CreateSymlink(ip)
	if err != nil {
		fmt.Println("Tokaido was unable to create a symlink on your PATH.")

		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("Successfully installed Tokaido (" + cs.String() + ")")

	// fmt.Println("No updates available to Tokaido at this time. Version " + lv.String() + " is the latest available.")
}
