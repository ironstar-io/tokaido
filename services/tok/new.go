package tok

import (
	"fmt"
	"log"
	"os"
	"time"

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
	ssh.StreamConnectCommand([]string{"composer", "create-project", "ironstar-io/d8-template:0.4", "/tmp/composer", "--stability", "dev", "--no-interaction"})
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

	console.Println("\n🍚  Creating a brand new Drupal 8 site with Tokaido!", "")

	// Create Tokaido configuration
	conf.SetDrupalConfig("DEFAULT")
	docker.FindOrCreateTokCompose()
	docker.CreateDatabaseVolume()
	docker.CreateSiteVolume()
	docker.CreateComposerCacheVolume()
	ssh.GenerateKeys()

	// Initial directory sync
	siteVolName := "tok_" + conf.GetConfig().Tokaido.Project.Name + "_tokaido_site"
	wo := console.SpinStart("Performing an initial sync...")
	utils.StdoutStreamCmdDebug("docker", "run", "-e", "AUTO_SYNC=false", "-v", conf.GetConfig().Tokaido.Project.Path+":/tokaido/host-volume", "-v", siteVolName+":/tokaido/site", "tokaido/sync:stable")
	console.SpinPersist(wo, "🚛", "Initial sync completed")

	// Lift drush container and configure
	drush.DockerUp()
	drupal.ConfigureSSH()
	xdebug.Configure()

	// Lift background sync container
	docker.UpMulti([]string{"sync"})

	console.Println("🎼  Using composer to generate the new Drupal project files. This might take some time", "")
	// `composer create-project` inside the drush container
	composerCreateProject()
	// Pull all required docker images
	docker.PullImages()

	wo = console.SpinStart("Waiting for sync to complete. This might take a few minutes...")
	filepath := drupal.SettingsPath()
	err := fs.WaitForSync(filepath, 180)
	if err != nil {
		console.Println("\n🙅‍  Your new Drupal site failed to sync from the Tokaido environment to your local disk", "")
		panic(err)
	}
	console.SpinPersist(wo, "🚋", "Secondary sync completed successfully")

	// Fire up the full Tokaido environment
	fmt.Println()
	wo = console.SpinStart("Tokaido is starting your containers")
	docker.Up()
	console.SpinPersist(wo, "🚅", "Tokaido started your containers")

	// Create Tokaido configuration for drupal post composer install
	drupal.CheckSettings("FORCE")

	// Setup HTTPS proxy service. Retain if statement to preserve Tokaido level enable/disable defaults
	setupProxy()

	wo = console.SpinStart("Waiting for settings.tok.php to sync. This might take a few minutes...")
	filepath = drupal.SettingsTokPath()
	err = fs.WaitForSync(filepath, 120)
	if err != nil {
		console.Println("\n🙅‍  Your new Drupal site failed to sync from the Tokaido environment to your local disk", "")
		panic(err)
	}

	// Induce a final wait for the sync process to complete
	time.Sleep(120 * time.Second)

	console.SpinPersist(wo, "🚋", "Final sync completed successfully")

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
	system.OpenTokaidoProxy(false)

}

// NewSuccessMessage ...
func NewSuccessMessage() {
	goos.InitMessage()
}
