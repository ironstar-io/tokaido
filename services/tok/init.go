package tok

import (
	"fmt"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/services/git"
	"github.com/ironstar-io/tokaido/services/proxy"
	"github.com/ironstar-io/tokaido/services/tok/goos"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/services/xdebug"
	"github.com/ironstar-io/tokaido/system"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/system/version"
)

// Init - The core run sheet of `tok up`
func Init() {
	c := conf.GetConfig()

	// System readiness checks
	console.Println("\n🚀  Tokaido is starting up!", "")
	system.CheckDependencies()
	version.GetUnisonVersion()
	version.Check()

	// Create Tokaido configuration
	conf.SetDrupalConfig("CUSTOM")
	drupal.CheckSettings("ASK")
	docker.FindOrCreateTokCompose()
	ssh.GenerateKeys()
	docker.CreateDatabaseVolume()

	docker.CreateComposerCacheVolume()

	git.IgnoreDefaults()

	// Run Unison for syncing
	unison.DockerUp()
	unison.CreateOrUpdatePrf(unison.LocalPort(), c.Tokaido.Project.Name, c.Tokaido.Project.Path)
	s := unison.SyncServiceStatus(c.Tokaido.Project.Name)
	if s == "stopped" {
		unison.Sync(c.Tokaido.Project.Name)
	}

	if c.System.Syncsvc.Enabled {
		fmt.Println()
		console.Println(`🔄  Creating a background process to sync your local repo into the Tokaido environment`, "")

		unison.CreateSyncService(c.Tokaido.Project.Name, c.Tokaido.Project.Path)
	}

	// Fire up the Docker environment
	if docker.ImageExists("tokaido/unison:2.51.2") == false && docker.ImageExists("tokaido/unison:2.48.4") == false {
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

	if c.System.Syncsvc.Enabled && c.System.Proxy.Enabled {
		console.Println("\n🔐  Setting up HTTPS for your local development environment", "")
		proxy.Setup()
	}
}

// InitMessage ...
func InitMessage() {
	goos.InitMessage()
}
