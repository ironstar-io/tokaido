package conf

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

// UnisonSync - Sync once without watching
func UnisonSync() {
	config := GetConfig()

	utils.StdoutCmd("unison", config.Project, "-watch=false")
}
