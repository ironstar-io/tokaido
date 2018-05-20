package osx

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

// CheckDeps - Root executable
func CheckDeps() string {
	utils.CheckPathHard("brew")

	CheckBrew()
	CheckAndBrewInstall("unison")
	CheckDockersync()

	return ""
}
