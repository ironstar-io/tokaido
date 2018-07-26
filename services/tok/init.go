package tok

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/services/drupal"
	"bitbucket.org/ironstar/tokaido-cli/services/git"
	"bitbucket.org/ironstar/tokaido-cli/services/tok/goos"
	"bitbucket.org/ironstar/tokaido-cli/services/unison"
	"bitbucket.org/ironstar/tokaido-cli/services/xdebug"
	"bitbucket.org/ironstar/tokaido-cli/system"
	"bitbucket.org/ironstar/tokaido-cli/system/console"
	"bitbucket.org/ironstar/tokaido-cli/system/ssh"
	"bitbucket.org/ironstar/tokaido-cli/system/version"

	"fmt"
)

// Init - The core run sheet of `tok up`
func Init() {
	c := conf.GetConfig()

	// System readiness checks
	console.Printf("\n🚀  Tokaido is starting up!\n", "")
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
		fmt.Println()
	}

	fmt.Println()

	// Fire up the Docker environment
	if docker.ImageExists("tokaido/drush-heavy:latest") == false {
		console.Println(`🚡  First time running Tokaido? There's a few images to download, this might take some time.`, "")
		fmt.Println()
	}

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
