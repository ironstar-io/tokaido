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
	fmt.Printf("\n🚀  Tokaido is starting up!\n")

	c := conf.GetConfig()
	system.CheckDependencies()

	docker.FindOrCreateTokCompose()

	ssh.GenerateKeys()

	drupal.CheckSettings()

	git.IgnoreDefaults()

	unison.DockerUp()
	unison.CreateOrUpdatePrf()

	// Only perform sync if sync service isn't running, otherwise we'll stall
	s := unison.SyncServiceStatus()
	if s == "stopped" {
		unison.Sync()
	}

	if c.CreateSyncService {
		fmt.Println()
		unison.CreateSyncService()
	}

	fmt.Println()

	if docker.ImageExists("tokaido/drush:latest") == false {
		fmt.Printf(`🚡  First time running Tokaido? There's a few images to download, this might take some time.`)
	}

	wo := wow.New(os.Stdout, spin.Get(spin.Dots), `   Tokaido is starting your containers`)
	wo.Start()

	docker.Up()

	drupal.ConfigureSSH()

	xdebug.Configure()

	wo.PersistWith(spin.Spinner{Frames: []string{"🚅"}}, `  Tokaido started your containers`)
}
