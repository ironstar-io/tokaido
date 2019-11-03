package version

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
)

// Upgrade - Check if update available and auto-upgrade the user
func Upgrade() {
	v := Get().Version
	cv := strings.Replace(v, "v", "", 0)
	cs, _ := semver.Parse(cv)

	ls, _ := GetLatest()

	// Checks if the latest version is Greater Than the current version
	if ls.GT(cs) == true {
		// Check if version exists on users' machine
		// If exists, Use `ln -s` to symlink this version

		// Else - Download latest release from GH
		// Use `ln -s` to symlink this version

		return
	}

	fmt.Println("No updates available to Tokaido at this time. Version " + ls.String() + " is the latest available.")
}
