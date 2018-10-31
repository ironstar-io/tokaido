package tok

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
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

var wg = sync.WaitGroup{}

func buildProjectFrame(projectName string) {
	if fs.CheckExists(projectName) == true {
		log.Fatal("The project '" + projectName + "' already exists in this directory. Exiting...")
	}

	fs.Mkdir(projectName)
	err := os.Chdir(projectName)
	if err != nil {
		panic(err)
	}

	git.Init()
	initialize.TokConfig()
}

func composerCreateProject() {
	cs := console.SpinStart("Composer is generating a new Drupal project. This might take some time")

	ssh.ConnectCommand([]string{"mkdir", "/tmp/composer"})
	ssh.ConnectCommand([]string{"composer", "create-project", "ironstar-io/d8-template:0.3", "/tmp/composer", "--stability", "dev", "--no-interaction"})
	ssh.ConnectCommand([]string{"cp", "-R", "/tmp/composer/*", "/tokaido/site"})

	console.SpinPersist(cs, "🎼", "Composer completed generation of a new Drupal project")

	wg.Done()
}
func createSyncService(c *conf.Config) {
	console.Println(`🔄  Creating a background process to sync your local repo into the Tokaido environment`, "")
	unison.CreateSyncService(c.Tokaido.Project.Name, c.Tokaido.Project.Path)

	wg.Done()
}
func dockerPullImages() {
	docker.PullImages()

	wg.Done()
}

func setupProxy(c *conf.Config) {
	if c.System.Syncsvc.Enabled && c.System.Proxy.Enabled {
		console.Println("\n🔐  Setting up HTTPS for your local development environment", "")
		proxy.Setup()
	}

	wg.Done()
}
func drushSiteInstall() {
	ssh.ConnectCommand([]string{"cd", "/tokaido/site/web", "&&", "drush site-install", "-y"})
	ssh.ConnectCommand([]string{"drush", "en", "swiftmailer", "password_policy", "password_policy_character_types", "password_policy_characters", "password_policy_username", "memcache", "health_check"})

	wg.Done()
}

func deduceProjectName(args []string) string {
	if len(args) == 0 {
		return constants.DefaultDrupalProjectName
	}

	return args[0]
}

// New - The core run sheet of `tok new {project}`
func New(args []string) {
	// Project frame
	pn := deduceProjectName(args)
	buildProjectFrame(pn)

	// System readiness checks
	system.CheckDependencies()
	version.GetUnisonVersion()

	console.Println("🍚  Creating a brand new Drupal 8 site with Tokaido!", "")

	c := conf.GetConfig()

	// Create Tokaido configuration
	conf.SetDrupalConfig("DEFAULT")
	drupal.CheckSettings()
	docker.FindOrCreateTokCompose()
	ssh.GenerateKeys()
	docker.CreateDatabaseVolume()
	git.IgnoreDefaults()

	// Lift the unison container and sync
	unison.DockerUp()
	unison.CreateOrUpdatePrf(unison.LocalPort(), c.Tokaido.Project.Name, c.Tokaido.Project.Path)
	unison.Sync(c.Tokaido.Project.Name)

	// Lift drush container and configure
	drush.DockerUp()
	drupal.ConfigureSSH()
	xdebug.Configure()

	// Batch composer create project and create sync service
	wg.Add(3)
	// `composer create-project` inside the drush container
	go composerCreateProject()
	// Create a unison sync service
	go createSyncService(c)
	// Pull all required docker images
	go dockerPullImages()
	// Wait until all processes complete
	wg.Wait()

	// Create Tokaido configuration for drupal post composer install
	drupal.CheckSettings()
	docker.CreateDatabaseVolume()

	// TODO: Memcache, Mailhog, and Adminer should be included in this install.

	// Fire up the full Tokaido environment
	fmt.Println()
	wo := console.SpinStart("Tokaido is starting your containers")
	docker.Up()
	console.SpinPersist(wo, "🚅", "Tokaido started your containers")

	// Batch proxy setup and drush site-install
	wg.Add(2)
	// Setup HTTPS proxy service. Retain if statement to preserve Tokaido level enable/disable defaults
	go setupProxy(c)
	// Drush site install, add additional packages
	go drushSiteInstall()
	wg.Wait()

	// Git stage all applicable files and commit
	git.AddAll()
	git.Commit("Initial Tokaido Configuration")

	// TODO: `tok status`
	// TODO: `tok test`

	// Provide the user a success message
	goos.InitMessage()

	// Open the main site at `https://<project-name>.local.tokaido.io:5154`
	system.OpenHaproxySite()
}

// NewSuccessMessage ...
func NewSuccessMessage() {
	goos.InitMessage()
}
