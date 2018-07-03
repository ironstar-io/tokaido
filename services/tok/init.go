package tok

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/services/drupal"
	"bitbucket.org/ironstar/tokaido-cli/services/git"
	"bitbucket.org/ironstar/tokaido-cli/services/unison"
	"bitbucket.org/ironstar/tokaido-cli/services/xdebug"
	"bitbucket.org/ironstar/tokaido-cli/system"
	"bitbucket.org/ironstar/tokaido-cli/system/ssh"

	"fmt"
	"os"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
)

// Init - The core run sheet of `tok init`
func Init() {
	w := wow.New(os.Stdout, spin.Get(spin.Dots), `   Tokaido is starting up!`)
	w.Start()

	c := conf.GetConfig()
	system.CheckDependencies()

	docker.FindOrCreateTokCompose()

	ssh.GenerateKeys()

	drupal.CheckSettings()

	git.IgnoreDefaults()

	unison.DockerUp()
	unison.CreateOrUpdatePrf()

	// Only perform sync if sync service isn't running, otherwise we'll stall
	err := unison.CheckSyncServiceSilent()
	if err != nil {
		unison.Sync()
	}

	w.PersistWith(spin.Spinner{}, "")

	if c.CreateSyncService {
		unison.CreateSyncService()
		fmt.Println()
	}

	wo := wow.New(os.Stdout, spin.Get(spin.Dots), `   Tokaido is pulling up your containers!`)
	wo.Start()

	docker.Up()

	drupal.ConfigureSSH()

	xdebug.Configure()

	wo.PersistWith(spin.Spinner{Frames: []string{"🚅"}}, `  Tokaido lifted your containers!`)
}
