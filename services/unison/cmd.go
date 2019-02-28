package unison

import (
	"fmt"

	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/utils"
)

var syncRunningErr = `There is already a sync service for this project running on your system. Exiting...`

// Sync - Sync once without watching
func Sync(syncName string) {
	s := SyncServiceStatus(syncName)
	if s == "running" {
		fmt.Println(syncRunningErr)
		return
	}

	fmt.Println()
	console.Println("⚡  Syncing files between your system and the Tokaido environment", "")
	utils.StdoutStreamCmdDebug("unison", syncName, "-watch=false")
}

// Watch ...
func Watch(syncName string) {
	s := SyncServiceStatus(syncName)
	if s == "running" {
		fmt.Println(syncRunningErr)
		return
	}

	fmt.Println(`Watching your files for changes and synchronising with your container
Please keep this command running in order to retain sync
	`)

	utils.StdoutStreamCmd("unison", syncName)
}
