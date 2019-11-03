package version

import (
	"fmt"
	"log"
	"strings"

	"github.com/blang/semver"
)

// Select - Change the users' Tokaido to their selection
func Select(selection string) {
	v := Get().Version
	cv := strings.Replace(v, "v", "", 0)
	cs, _ := semver.Parse(cv)

	sv, err := semver.Parse(selection)
	if err != nil {
		log.Fatalf("Invalid semver selection supplied. Exiting...")

		return
	}

	// Checks if the current version is Equal to the selected version and exits if so
	if sv.EQ(cs) == true {
		log.Fatalf("Selected version (" + sv.String() + ") is the same as the currently installed version. Exiting")

		return
	}

	fmt.Println("Change to: " + sv.String())
	// Check if version exists on users' machine
	// If exists, Use `ln -s` to symlink this version

	// Check if version is released in GH
	// If exists, download this release from GH
	// Use `ln -s` to symlink this version

	// If both conditions false, can't change to this version warning and exit
}
