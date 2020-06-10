package tok

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/utils"
)

// Destroy ...
func Destroy(name string) {
	confirmDestroy := utils.ConfirmationPrompt(`🔥  This will shut down your Tokaido environment for this project. Are you sure?`, "n")
	if confirmDestroy == false {
		console.Println(`🍵  Exiting without change`, "")
		return
	}
	docker.Down()
	conf.DeregisterProject(name)
}
