package tok

import (
	"fmt"
	"log"
	"os"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/services/drush"
	"github.com/ironstar-io/tokaido/services/git"
	"github.com/ironstar-io/tokaido/services/proxy"
	"github.com/ironstar-io/tokaido/services/tok/goos"
	"github.com/ironstar-io/tokaido/services/xdebug"
	"github.com/ironstar-io/tokaido/system"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/utils"
)

var profileFlag string

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
	initialize.TokConfig("new")
}

func composerCreateProject() {
	ssh.ConnectCommand([]string{"mkdir", "/tmp/composer"})
	ssh.StreamConnectCommand([]string{"composer", "create-project", "ironstar-io/d8-template:0.3", "/tmp/composer", "--stability", "dev", "--no-interaction"})
	ssh.ConnectCommand([]string{"cp", "-R", "/tmp/composer/*", "/tokaido/site"})
}

func setupProxy() {
	c := conf.GetConfig()
	if c.System.Syncsvc.Enabled && c.System.Proxy.Enabled {
		console.Println("🔐  Setting up HTTPS for your local development environment", "")
		proxy.Setup()
	}
}
func drushSiteInstall(profile string) {
	switch profile {
	case "":
		profile = "standard"
		fallthrough
	case
		"standard",
		"minimal",
		"demo_umami":
		ssh.StreamConnectCommand([]string{"drush", "site-install", profile, "-y"})
		ssh.StreamConnectCommand([]string{"drush", "en", "swiftmailer", "password_policy", "password_policy_character_types", "password_policy_characters", "password_policy_username", "memcache", "health_check"})
		return
	}
	log.Fatalf("Error: The install profile specified was not supported. Possible values are 'standard', 'minimal', and 'demo_umami'")
}

func deduceProjectName(args []string) string {
	if len(args) == 0 {
		return constants.DefaultDrupalProjectName
	}

	return args[0]
}

// New - The core run sheet of `tok new {project}`
func New(args []string, profile string) {
	// Project frame
	pn := deduceProjectName(args)
	buildProjectFrame(pn)

	// System readiness checks
	system.CheckDependencies()

	console.Println("🍚  Creating a brand new Drupal 8 site with Tokaido!", "")

	// Create Tokaido configuration
	conf.SetDrupalConfig("DEFAULT")
	docker.FindOrCreateTokCompose()
	docker.CreateDatabaseVolume()
	docker.CreateComposerCacheVolume()
	ssh.GenerateKeys()

	// Initial directory sync
	wo := console.SpinStart("Performing an initial sync...")
	utils.StdoutStreamCmdDebug("docker", "run", "-e", "AUTO_SYNC=false", "-v", conf.GetConfig().Tokaido.Project.Path+":/tokaido/host-volume", "tokaido/sync:stable")
	console.SpinPersist(wo, "🚛", "Initial sync completed")

	// Lift drush container and configure
	drush.DockerUp()
	drupal.ConfigureSSH()
	xdebug.Configure()

	console.Println("🎼  Using composer to generate the new Drupal project files. This might take some time", "")
	// `composer create-project` inside the drush container
	composerCreateProject()
	// Pull all required docker images
	docker.PullImages()

	// Sync service is async and causes a race condition due to the large number of files changed.
	// Must manually sync until stable
	wo = console.SpinStart("Resyncing Tokaido changes back to your local...")
	utils.StdoutStreamCmdDebug("docker", "run", "-e", "AUTO_SYNC=false", "-v", conf.GetConfig().Tokaido.Project.Path+":/tokaido/host-volume", "tokaido/sync:stable")
	console.SpinPersist(wo, "🚋", "Secondary sync completed")

	// Fire up the full Tokaido environment
	fmt.Println()
	wo = console.SpinStart("Tokaido is starting your containers")
	docker.Up()
	console.SpinPersist(wo, "🚅", "Tokaido started your containers")

	// Create Tokaido configuration for drupal post composer install
	drupal.CheckSettings("FORCE")

	// Setup HTTPS proxy service. Retain if statement to preserve Tokaido level enable/disable defaults
	setupProxy()

	// Drush site install, add additional packages
	console.Println(`💧  Running drush site-install for your new project`, "")
	drushSiteInstall(profile)

	// Generate a new .gitignore file
	git.NewGitignore()
	// Git stage all applicable files and commit
	git.AddAll()
	git.Commit("Initial Tokaido Configuration")

	// TODO: `tok status`
	// TODO: `tok test`

	// Provide the user a success message
	goos.InitMessage()

	// Open the main site at `https://<project-name>.local.tokaido.io:5154`
	system.OpenHaproxySite(false)

}

// NewSuccessMessage ...
func NewSuccessMessage() {
	goos.InitMessage()
}
