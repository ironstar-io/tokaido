package unison

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
)

var syncRunningErr = `There is already a sync service for this project running on your system. Exiting...`

// Sync - Sync once without watching
func Sync() {
	s := SyncServiceStatus()
	if s == "running" {
		fmt.Println(syncRunningErr)
		return
	}

	config := conf.GetConfig()

	utils.CommandSubstitution("unison", config.Project, "-watch=false")
}

// Watch ...
func Watch() {
	s := SyncServiceStatus()
	if s == "running" {
		fmt.Println(syncRunningErr)
		return
	}

	config := conf.GetConfig()

	fmt.Println(`Watching your files for changes and synchronising with your container`)

	utils.StdoutStreamCmd("unison", config.Project)
}
