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
	cs, _ := semver.Parse(cv)

	ls, _ := GetLatest()
	lv, _ := semver.Parse(ls.TagName)

	// Checks if the latest version is Greater Than the current version
	if lv.GT(cs) == true {
		confirmUpgrade := utils.ConfirmationPrompt("This will upgrade your Tokaido version to latest ("+lv.String()+").\n\nAre you sure?", "y")
		if confirmUpgrade == false {
			fmt.Println("Exiting...")
			return
		}

		ip := GetInstallPath(lv.String())
		// Empty string if not installed, in which case, download and symlink
		if ip == "" {
			// Download latest release from GH
			p, err := Install(lv.String())
			if err != nil {
				fmt.Println("Tokaido was unable to upgrade you to the latest version.")

				log.Fatal(err)
			}

			ip = p
		}

		// Latest version is installed, just not active. Create a Symlink to finish
		err := CreateSymlink(ip)
		if err != nil {
			fmt.Println("Tokaido was unable to upgrade you to the latest version.")

			log.Fatal(err)
		}

		return
	}

	fmt.Println("No updates available to Tokaido at this time. Version " + lv.String() + " is the latest available.")
}
