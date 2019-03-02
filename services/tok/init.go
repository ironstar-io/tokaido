package tok

import (
	"fmt"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/services/git"
	"github.com/ironstar-io/tokaido/services/proxy"
	"github.com/ironstar-io/tokaido/services/snapshots"
	"github.com/ironstar-io/tokaido/services/tok/goos"
	"github.com/ironstar-io/tokaido/services/xdebug"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/utils"
)

// Init - The core run sheet of `tok up`
func Init() {
	c := conf.GetConfig()

	// System readiness checks
	console.Println("\n🚀  Tokaido is starting up!", "")
	version.Check()

	// Create Tokaido configuration
	conf.SetDrupalConfig("CUSTOM")
	drupal.CheckSettings("ASK")
	docker.FindOrCreateTokCompose()
	ssh.GenerateKeys()
	docker.CreateDatabaseVolume()
	docker.CreateSiteVolume()
	err := snapshots.Init()
	if err != nil {
		fmt.Println()
		console.Println("🙅  Tokaido encountered an unexpected error preparing the database snapshot service", "")
		panic(err)
	}

	docker.CreateComposerCacheVolume()

	git.IgnoreDefaults()

	// Fusion Sync WIP
	wo := console.SpinStart("Performing an initial sync. This might take a few minutes")
	siteVolName := "tok_" + conf.GetConfig().Tokaido.Project.Name + "_tokaido_site"
	utils.StdoutStreamCmdDebug("docker", "run", "-e", "AUTO_SYNC=false", "-v", conf.GetConfig().Tokaido.Project.Path+":/tokaido/host-volume", "-v", siteVolName+":/tokaido/site", "tokaido/sync:stable")
	console.SpinPersist(wo, "🚛", "Initial sync completed")
	wo = console.SpinStart("Tokaido is starting your containers")

	docker.Up()

	// Perform post-launch configuration
	drupal.ConfigureSSH()
	xdebug.Configure()

	console.SpinPersist(wo, "🚅", "Tokaido containers were started")

	if c.System.Syncsvc.Enabled && c.System.Proxy.Enabled {
		wo = console.SpinStart("Setting up secure HTTPS access")
		proxy.Setup()
		console.SpinPersist(wo, "🔐", "Successfully configured HTTPS")
	}

	err = docker.StatusCheck()
	if err == nil {
		fmt.Println()
		console.Println(`✅  All containers are running`, "√")
	}

	err = ssh.CheckKey()

	err = drupal.CheckContainer()

	if err == nil {
		console.Println(`🍜  Tokaido started up successfully`, "")
	} else {
		fmt.Println()
		console.Println("🙅  Uh oh! It looks like Tokaido didn't start properly.", "")
		console.Println("    Come find us in #tokaido on the Drupal Slack if you need some help", "")
		fmt.Println()
	}
}

// InitMessage ...
func InitMessage() {
	goos.InitMessage()
}
