package tok

import (
	"fmt"
	"os"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/services/drush"
	"github.com/ironstar-io/tokaido/services/git"
	"github.com/ironstar-io/tokaido/services/proxy"
	"github.com/ironstar-io/tokaido/services/tok/goos"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/services/xdebug"
	"github.com/ironstar-io/tokaido/system"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/system/version"
)

// New - The core run sheet of `tok new {project}`
func New(args []string) {
	console.Println("\n🍚  Creating a brand new Drupal 8 site with Tokaido!", "")

	// System readiness checks
	system.CheckDependencies()
	version.GetUnisonVersion()

	var projectName string
	if len(args) == 0 {
		projectName = "temp-test-drupal-site"
	} else {
		projectName = args[0]
	}

	// TODO: Error if folder name already exists
	fs.Mkdir(projectName)
	err := os.Chdir(projectName)
	if err != nil {
		panic(err)
	}

	fmt.Println(os.Getwd())

	git.Init()
	initialize.TokConfig()

	c := conf.GetConfig()

	// Create Tokaido configuration
	// TODO: Allow application of Drupal defaults
	conf.SetDrupalConfig()
	drupal.CheckSettings()
	docker.FindOrCreateTokCompose()
	ssh.GenerateKeys()
	docker.CreateDatabaseVolume()

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

	drush.DockerUp()

	// Perform post-launch configuration
	drupal.ConfigureSSH()
	xdebug.Configure()

	// TODO: Move composer setup block to drush package
	a := ssh.ConnectCommand([]string{"mkdir", "/tmp/composer"})
	fmt.Println(a)
	x := ssh.ConnectCommand([]string{"composer", "create-project", "ironstar-io/d8-template:0.3", "/tmp/composer", "--stability", "dev", "--no-interaction"})
	fmt.Println(x)
	y := ssh.ConnectCommand([]string{"cp", "-R", "/tmp/composer", "/tokaido/site"})
	fmt.Println(y)

	fmt.Println("OK")

	return

	// Fire up the Docker environment
	if docker.ImageExists("tokaido/drush-heavy:latest") == false {
		console.Println(`🚡  First time running Tokaido? There's a few images to download, this might take some time.`, "")
		fmt.Println()
	}

	fmt.Println()

	wo := console.SpinStart("Tokaido is starting your containers")

	docker.Up()

	// // Perform post-launch configuration
	// drupal.ConfigureSSH()
	// xdebug.Configure()

	console.SpinPersist(wo, "🚅", "Tokaido started your containers")

	if c.System.Syncsvc.Enabled && c.System.Proxy.Enabled {
		console.Println("\n🔐  Setting up HTTPS for your local development environment", "")
		proxy.Setup()
	}
}

// NewSuccessMessage ...
func NewSuccessMessage() {
	goos.InitMessage()
}
