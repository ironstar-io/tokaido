package unison

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
)

// Sync - Sync once without watching
func Sync() {
	config := conf.GetConfig()

	fmt.Println(`Running the initial synchronisation of your files with the container.`)

	utils.NoFatalStdoutCmd("unison", config.Project, "-watch=false")
}
