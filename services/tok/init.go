package tok

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/services/git"
	"github.com/ironstar-io/tokaido/services/proxy"
	"github.com/ironstar-io/tokaido/services/snapshots"
	"github.com/ironstar-io/tokaido/services/telemetry"
	"github.com/ironstar-io/tokaido/services/tok/goos"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/services/xdebug"
	"github.com/ironstar-io/tokaido/system"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/system/version"
	"github.com/ironstar-io/tokaido/utils"

	. "github.com/logrusorgru/aurora"
)

// Init - The core run sheet of `tok up`
func Init(yes, statuscheck, noPullImagesFlag bool) {
	c := conf.GetConfig()
	cs := "ASK"
	if yes {
		cs = "FORCE"
	}

	startTime := time.Now().UTC()

	// System readiness checks
	version.Check()
	docker.CheckClientVersion()
	checkSyncConfig()

	fmt.Println(Cyan("\n🚀  Tokaido is starting up!"))

	// Add this project to the global configuration
	pr := fs.ProjectRoot()
	name := strings.Replace(filepath.Base(pr), ".", "", -1)
	conf.RegisterProject(name, pr)

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

	// Telemetry actions
	telemetry.GenerateProjectID()
	telemetry.RequestOptIn()
	telemetry.SendGlobal()
	telemetry.SendProject(startTime, 0)

	git.IgnoreDefaults()
	docker.CreateComposerCacheVolume()

	if c.Global.Syncservice == "unison" {
		unison.DockerUp()
		unison.CreateOrUpdatePrf(unison.LocalPort(), c.Tokaido.Project.Name, pr)
		s := unison.SyncServiceStatus(c.Tokaido.Project.Name)
		if s == "stopped" {
			unison.Sync(c.Tokaido.Project.Name)
		}
		fmt.Println()
		console.Println(`🔄  Creating a background process to sync your local repo into the Tokaido environment`, "")
		unison.CreateSyncService(c.Tokaido.Project.Name, pr)
	}

	if c.Global.Syncservice == "fusion" {
		docker.CreateComposerCacheVolume()
		utils.DebugString("Using Fusion Sync Service")
		wo := console.SpinStart("Performing an initial sync. This might take a few minutes")
		siteVolName := "tok_" + conf.GetConfig().Tokaido.Project.Name + "_tokaido_site"
		utils.StdoutStreamCmdDebug("docker", "run", "--rm", "-e", "AUTO_SYNC=false", "-v", conf.GetProjectPath()+":/tokaido/host-volume", "-v", siteVolName+":/tokaido/site", "tokaido/sync:stable")
		console.SpinPersist(wo, "🚛", "Initial sync completed")
	}

	if noPullImagesFlag {
		fmt.Println("    Not downloading the latest Docker images")
	} else {
		fmt.Println("🤖  Downloading the latest Docker images")
		docker.PullImages()
	}

	console.Println("🚅  Starting your Drupal environment", "")
	docker.Up()
	surveyMessage()
	// Perform post-launch configuration
	drupal.ConfigureSSH()
	xdebug.Configure()

	if c.Global.Proxy.Enabled {
		// This step can't be in a spinner because the spinner can't ask for user input during the SSL trust stage.
		console.Println(`🔐  Setting up HTTPS access...`, "")
		proxy.Setup()
	}

	// Check docker containers
	ok := docker.StatusCheck("", conf.GetConfig().Tokaido.Project.Name)
	if !ok {
		fmt.Println()
		fmt.Println(Red("😓  Tokaido containers are not working properly"))

		if !docker.StatusCheck("fpm", conf.GetConfig().Tokaido.Project.Name) {
			fmt.Println(Red("    The Tokaido FPM container failed to start up."))
			fmt.Println(Red("    This most likely suggests a problem with your Drupal site that is causing PHP to crash."))
			fmt.Println()
			fmt.Println(Sprintf(Cyan("    You can try running '%s' to see the full PHP startup log"), Blue("tok logs fpm")))
		} else if !docker.StatusCheck("nginx", conf.GetConfig().Tokaido.Project.Name) {
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

	// Send a final telemetry
	duration := time.Now().Sub(startTime)
	telemetry.SendProject(startTime, int(duration.Seconds()))

}

// 1 in 6 chance of displaying survey message
func surveyMessage() {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(6-1) + 1
	if n == 3 {
		fmt.Println(Sprintf("🤗  How's Tokaido? Run '%s' to share your feedback.", Bold("tok survey")))
	}
}

// InitMessage ...
func InitMessage() {
	goos.InitMessage()
}

// checkSyncConfig verifies if the user's syncservice setting is compatible with their system
func checkSyncConfig() {
	c := conf.GetConfig()

	switch system.CheckOS() {
	case "osx":
		if c.Global.Syncservice == "fusion" {
			fmt.Println(Yellow("Warning: The Fusion Sync Service will be removed in Tokaido 1.10. Please migrate to 'docker' or 'unison'"))
		}
	case "linux":
		if c.Global.Syncservice != "unison" {
			fmt.Println(Sprintf(Yellow("Warning: The syncservice '%s' is not compatible with Linux. Tokaido will automatically be set to use Unison\n\n"), Bold(c.Global.Syncservice)))
			conf.SetGlobalConfigValueByArgs([]string{"syncservice", "unison"})
		}
	}

	if c.Global.Syncservice == "unison" {
		unison.SystemCompatibilityChecks()
	}

}
