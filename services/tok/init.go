package tok

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/services/git"
	"github.com/ironstar-io/tokaido/services/tok/goos"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/services/xdebug"
	"github.com/ironstar-io/tokaido/system"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/system/version"

	"fmt"
)

// Init - The core run sheet of `tok up`
func Init() {
	c := conf.GetConfig()

	// System readiness checks
	console.Println("\n🚀  Tokaido is starting up!", "")
	system.CheckDependencies()
	version.GetUnisonVersion()
	git.CheckGitRepo()

	// Create Tokaido configuration
	docker.FindOrCreateTokCompose()
	ssh.GenerateKeys()
	conf.SetDefaultDrupalVars()
	drupal.CheckSettings()
	git.IgnoreDefaults()

	// Run Unison for syncing
	unison.DockerUp()
	unison.CreateOrUpdatePrf()
	s := unison.SyncServiceStatus()
	if s == "stopped" {
		unison.Sync()
	}

	if c.System.Syncsvc.Enabled {
		unison.CreateSyncService()
		fmt.Println()
	}

	// Fire up the Docker environment
	if docker.ImageExists("tokaido/drush-heavy:latest") == false {
		console.Println(`🚡  First time running Tokaido? There's a few images to download, this might take some time.`, "")
		fmt.Println()
	}

	fmt.Println()

	wo := console.SpinStart("Tokaido is starting your containers")

	docker.Up()

	// Perform post-launch configuration
	drupal.ConfigureSSH()
	xdebug.Configure()

	console.SpinPersist(wo, "🚅", "Tokaido started your containers")
}

// InitMessage ...
func InitMessage() {
	goos.InitMessage()
}
