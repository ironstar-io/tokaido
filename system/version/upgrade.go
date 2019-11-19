package version

import (
	"fmt"
	"log"
	"strings"

	"github.com/ironstar-io/tokaido/utils"

	"github.com/blang/semver"
)

// SelfInstall - Check if update available and auto-upgrade the user
func SelfInstall() {
	v := Get().Version
	cv := strings.Replace(v, "v", "", 0)
	cs, err := semver.Parse(cv)
	if err != nil {
		fmt.Println("Tokaido was unable to correctly parse the current version.")

		log.Fatal(err)
	}

	ls, err := GetLatest()
	if err != nil {
		fmt.Println("Tokaido was unable to reach GitHub to determine the latest version")

		log.Fatal(err)
	}

	lv, _ := semver.Parse(ls.TagName)
	if err != nil {
		fmt.Println("Tokaido was unable to correctly parse the latest version.")

		log.Fatal(err)
	}

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
			p, err := DownloadAndInstall(lv.String())
			if err != nil {
				fmt.Println("Tokaido wasn't able to upgrade you to the latest version.")

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

		fmt.Println()
		fmt.Println("Successfully upgraded Tokaido to the latest version (" + lv.String() + ")")

		return
	}

	fmt.Println("No updates available to Tokaido at this time. Version " + lv.String() + " is the latest available.")
}
