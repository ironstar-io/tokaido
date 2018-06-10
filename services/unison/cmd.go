package unison

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
)

// Sync - Sync once without watching
func Sync() {
	config := conf.GetConfig()

	// fmt.Println(`🚝  Running a one time Unison sync between your system and the Tokaido environment`)

	utils.NoFatalStdoutCmd("unison", config.Project, "-watch=false")
}

// Watch ...
func Watch() {
	config := conf.GetConfig()

	fmt.Println(`Watching your files for changes and synchronising with your container`)

	utils.StdoutCmd("unison", config.Project)
}
