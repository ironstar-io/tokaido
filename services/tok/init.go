package tok

import (
	"fmt"
	"os"

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
	"github.com/ironstar-io/tokaido/system/version"
	"github.com/ironstar-io/tokaido/utils"

	. "github.com/logrusorgru/aurora"
)

// Init - The core run sheet of `tok up`
func Init(yes, statuscheck bool) {
	c := conf.GetConfig()
	cs := "ASK"
	if yes {
		cs = "FORCE"
	}

	// System readiness checks
	version.Check()
	docker.CheckClientVersion()
	fmt.Println(Cyan("\n🚀  Tokaido is starting up!"))

	// Create Tokaido configuration
	conf.SetDrupalConfig("CUSTOM")
	drupal.CheckSettings(cs)
	docker.FindOrCreateTokCompose()
	ssh.GenerateKeys()
	docker.CreateDatabaseVolume()
	docker.CreateSiteVolume()
	err := snapshots.Init()
	if err != nil {
		fmt.Println()
		fmt.Println(Red("🙅  Tokaido encountered an unexpected error preparing the database snapshot service"))
		panic(err)
	}

	docker.CreateComposerCacheVolume()

	git.IgnoreDefaults()

	if c.Global.Syncservice == "fusion" {
		utils.DebugString("Using Fusion Sync Service")
		wo := console.SpinStart("Performing an initial sync. This might take a few minutes")
		siteVolName := "tok_" + conf.GetConfig().Tokaido.Project.Name + "_tokaido_site"
		utils.StdoutStreamCmdDebug("docker", "run", "-e", "AUTO_SYNC=false", "-v", conf.GetConfig().Tokaido.Project.Path+":/tokaido/host-volume", "-v", siteVolName+":/tokaido/site", "tokaido/sync:stable")
		console.SpinPersist(wo, "🚛", "Initial sync completed")
	}

	wo := console.SpinStart("Downloading the latest Docker images")
	docker.PullImages()
	console.SpinPersist(wo, "🤖", "Latest Docker images downloaded successfully")

	// wo = console.SpinStart("Starting your containers")
	docker.Up()
	// console.SpinPersist(wo, "🚅", "Tokaido containers are online")

	console.Println("🚅  Configuring your new Drupal environment", "")
	// Perform post-launch configuration
	drupal.ConfigureSSH()
	xdebug.Configure()

	if c.System.Proxy.Enabled {
		// This step can't be in a spinner because the spinner can't ask for user input during the SSL trust stage.
		console.Println(`🔐  Setting up HTTPS access...`, "")
		proxy.Setup()
	}

	// Check docker containers
	ok := docker.StatusCheck("")
	if !ok {
		fmt.Println()
		fmt.Println(Red("😓  Tokaido containers are not working properly"))

		if !docker.StatusCheck("fpm") {
			fmt.Println(Red("    The Tokaido FPM container failed to start up."))
			fmt.Println(Red("    This most likely suggests a problem with your Drupal site that is causing PHP to crash."))
			fmt.Println()
			fmt.Println(Sprintf(Cyan("    You can try running '%s' to see the full PHP startup log"), Blue("tok logs fpm")))
		} else if !docker.StatusCheck("nginx") {
			fmt.Println(Red("    The Tokaido NGINX container failed to start up."))
			fmt.Println(Red("    This is most likely caused by a Tokaido misconfiguration."))
			fmt.Println()
			fmt.Println(Sprintf(Cyan("    You can try running '%s' to see the full startup log"), Blue("tok logs nginx")))
		}

		fmt.Println()
		fmt.Println(Sprintf("    You can view the status of your containers with '%s' and you can see", Bold("tok ps")))
		fmt.Println(Sprintf("    error logs by running '%s', such as '%s'", Bold("tok logs {container-name}"), Bold("tok logs mysql")))
		fmt.Println()
		fmt.Println(Cyan("    You can try to fix this by running 'tok repair'"))
		fmt.Println()

		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(Green(`🙂  All containers are running`))

	// Final startup checks
	if statuscheck {
		ok = ssh.CheckKey()
		if ok {
			fmt.Println(Green("😀  SSH access is configured"))
		} else {
			fmt.Println(Red("😓  SSH access is not configured"))
			fmt.Println("  Your SSH access to the Drush container looks broken.")
			fmt.Println("  You should be able to run 'tok repair' to attempt to fix this automatically")
		}

		if ok {
			fmt.Println(Green(`🍜  Tokaido started up successfully`))
		} else {
			fmt.Println()
			console.Println("🙅  Uh oh! It looks like Tokaido didn't start properly.", "")
			console.Println("    Come find us in #tokaido on the Drupal Slack if you need some help", "")
			fmt.Println()
		}
	}

}

// InitMessage ...
func InitMessage() {
	goos.InitMessage()
}
