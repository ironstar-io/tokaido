package version

import (
	"fmt"
	"os"
	"strings"

	"github.com/ironstar-io/tokaido/system/version/goos"
	"github.com/ironstar-io/tokaido/utils"

	"github.com/blang/semver"
)

// Upgrade - Check if update available and auto-upgrade the user
func Upgrade() {
	v := Get().Version
	cv := strings.Replace(v, "v", "", 0)
	cs, err := semver.Parse(cv)
	if err != nil {
		fmt.Println("Tokaido was unable to correctly parse the current version: ", err.Error())
		os.Exit(1)
	}

	ls, err := GetLatest()
	if err != nil {
		fmt.Println("Tokaido was unable to reach GitHub to determine the latest version: ", err.Error())
		os.Exit(1)
	}

	lv, _ := semver.Parse(ls.TagName)
	if err != nil {
		fmt.Println("Tokaido was unable to correctly parse the latest version: ", err.Error())
		os.Exit(1)
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
				fmt.Println("Tokaido wasn't able to upgrade you to the latest version: ", err.Error())
				os.Exit(1)
			}

			ip = p
		}

		// Latest version is installed, just not active. Create a Symlink to finish
		goos.ActivateSavedVersion(ip)

		fmt.Println()
		fmt.Println("Successfully upgraded Tokaido to the latest version (" + lv.String() + ")")

		return
	}

	fmt.Println("No updates available to Tokaido at this time. Version " + lv.String() + " is the latest available.")
}
