package tok

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/proxy"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/utils"

	"fmt"
)

// Destroy ...
func Destroy() {
	confirmDestroy := utils.ConfirmationPrompt(`🔥  This will also destroy the database inside your Tokaido environment. Are you sure?`, "n")
	if confirmDestroy == false {
		console.Println(`🍵  Exiting without change`, "")
		return
	}
	fmt.Println()

	proxy.DestroyProject()
	fmt.Println()

	docker.Down()

	unison.StopSyncService(conf.GetConfig().Tokaido.Project.Name)

	proxy.DockerComposeUp()
}
