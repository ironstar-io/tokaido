package tok

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/proxy"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/utils"

	"fmt"
)

// Destroy ...
func Destroy(name string) {
	confirmDestroy := utils.ConfirmationPrompt(`🔥  This will shut down your Tokaido environment for this project. Are you sure?`, "n")
	if confirmDestroy == false {
		console.Println(`🍵  Exiting without change`, "")
		return
	}
	fmt.Println()

	proxy.DestroyProject()
	fmt.Println()

	docker.Down()

	proxy.DockerComposeUp()

	conf.DeregisterProject(name)
}
