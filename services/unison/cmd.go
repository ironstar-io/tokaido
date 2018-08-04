package unison

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/utils"

	"fmt"
	"runtime"
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

	if runtime.GOOS == "windows" {
		fmt.Println("Synchronizing your files between your local filesystem and the container. This may take some time.")
	}

	utils.CommandSubstitution("unison", config.Tokaido.Project.Name, "-watch=false")
}

// Watch ...
func Watch() {
	s := SyncServiceStatus()
	if s == "running" {
		fmt.Println(syncRunningErr)
		return
	}

	fmt.Println(`Watching your files for changes and synchronising with your container
Please keep this command running in order to retain sync
	`)

	utils.StdoutStreamCmd("unison", conf.GetConfig().Tokaido.Project.Name)
}
