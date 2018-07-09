package tok

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/services/drupal"
	"bitbucket.org/ironstar/tokaido-cli/services/git"
	"bitbucket.org/ironstar/tokaido-cli/services/unison"
	"bitbucket.org/ironstar/tokaido-cli/services/xdebug"
	"bitbucket.org/ironstar/tokaido-cli/system"
	"bitbucket.org/ironstar/tokaido-cli/system/version"
	// "bitbucket.org/ironstar/tokaido-cli/system/ssh"

	"fmt"
	"os"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
)

// Init - The core run sheet of `tok init`
func Init() {
	c := conf.GetConfig()

	// System readiness checks
	fmt.Printf("\n🚀  Tokaido is starting up!\n")
	system.CheckDependencies()
	version.GetUnisonVersion()
	git.CheckGitRepo()

	// Create Tokaido configuration
	docker.FindOrCreateTokCompose()
	ssh.GenerateKeys()
	drupal.CheckSettings()
	git.IgnoreDefaults()

	// Run Unison for syncing
	unison.DockerUp()
	unison.CreateOrUpdatePrf()
	s := unison.SyncServiceStatus()
	if s == "stopped" {
		unison.Sync()
	}

	if c.CreateSyncService {
		fmt.Println()
		unison.CreateSyncService()
	}

	fmt.Println()

	// Fire up the Docker environment
	if docker.ImageExists("tokaido/drush-heavy:latest") == false {
		fmt.Printf(`🚡  First time running Tokaido? There's a few images to download, this might take some time.`)
	}
	wo := wow.New(os.Stdout, spin.Get(spin.Dots), `   Tokaido is starting your containers`)
	wo.Start()
	docker.Up()

	// Perform post-launch configuration
	drupal.ConfigureSSH()
	xdebug.Configure()

	wo.PersistWith(spin.Spinner{Frames: []string{"🚅"}}, `  Tokaido started your containers`)
}
