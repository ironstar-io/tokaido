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
	"github.com/ironstar-io/tokaido/services/xdebug"
	"github.com/ironstar-io/tokaido/system"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/system/tls"
	"github.com/ironstar-io/tokaido/system/version"

	"github.com/logrusorgru/aurora"
)

// Init - The core run sheet of `tok up`
func Init(yes, statuscheck bool) {
	cs := "ASK"
	if yes {
		cs = "FORCE"
	}

	startTime := time.Now().UTC()

	// System readiness checks
	version.Check()
	docker.CheckClientVersion()
	proxy.DecommissionProxy()
	system.CheckDependencies()

	fmt.Println(aurora.Cyan("\n🚀  Tokaido is starting up!"))

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
	docker.CreateLogsVolume()
	err := snapshots.Init()
	if err != nil {
		fmt.Println()
		fmt.Println(aurora.Red("🙅  Tokaido encountered an unexpected error preparing the database snapshot service    "))
		panic(err)
	}

	// Telemetry actions
	telemetry.GenerateProjectID()
	telemetry.RequestOptIn()
	telemetry.SendGlobal()
	telemetry.SendProject(startTime, 0)

	git.IgnoreDefaults()
	docker.CreateComposerCacheVolume()

	c := conf.GetConfig()
	if c.Global.Syncservice == "unison" || c.Global.Syncservice == "fusion" {
		fmt.Println(aurora.Red("    Sorry, the sync service you selected is not supported."))
	}

	// Configure TLS
	console.Println("🔐  Configuring TLS Certificates", "")
	tls.ConfigureTLS()

	console.Println("🚅  Starting your Drupal environment", "")
	docker.Up()
	surveyMessage()
	// Perform post-launch configuration
	drupal.ConfigureSSH()
	xdebug.Configure()

	// Check docker containers
	ok := docker.StatusCheck("", conf.GetConfig().Tokaido.Project.Name)
	if !ok {
		fmt.Println()
		fmt.Println(aurora.Red("😓  Tokaido containers are not working properly    "))

		if !docker.StatusCheck("fpm", conf.GetConfig().Tokaido.Project.Name) {
			fmt.Println(aurora.Red("    The Tokaido FPM container failed to start up."))
			fmt.Println(aurora.Red("    This most likely suggests a problem with your Drupal site that is causing PHP to crash."))
			fmt.Println()
			fmt.Println(aurora.Sprintf(aurora.Cyan("    You can try running '%s' to see the full PHP startup log"), aurora.Blue("tok logs fpm")))
		} else if !docker.StatusCheck("nginx", conf.GetConfig().Tokaido.Project.Name) {
			fmt.Println(aurora.Red("    The Tokaido NGINX container failed to start up."))
			fmt.Println(aurora.Red("    This is most likely caused by a Tokaido misconfiguration."))
			fmt.Println()
			fmt.Println(aurora.Sprintf(aurora.Cyan("    You can try running '%s' to see the full startup log"), aurora.Blue("tok logs nginx")))
		}

		fmt.Println()
		fmt.Println(aurora.Sprintf("    You can view the status of your containers with '%s' and you can see", aurora.Bold("tok ps")))
		fmt.Println(aurora.Sprintf("    error logs by running '%s', such as '%s'", aurora.Bold("tok logs {container-name}"), aurora.Bold("tok logs mysql")))
		fmt.Println()
		fmt.Println(aurora.Cyan("    You can try to fix this by running 'tok repair'"))
		fmt.Println()

		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(aurora.Green("🙂  All containers are running    "))

	// Final startup checks
	if statuscheck {
		ok = ssh.CheckKey()
		if ok {
			fmt.Println(aurora.Green("😀  SSH access is configured    "))
		} else {
			fmt.Println(aurora.Red("😓  SSH access is not configured    "))
			fmt.Println("  Your SSH access to the Drush container looks broken.")
			fmt.Println("  You should be able to run 'tok repair' to attempt to fix this automatically")
		}

		if ok {
			fmt.Println(aurora.Green(`🍜  Tokaido started up successfully    `))
		} else {
			fmt.Println()
			console.Println("🙅  Uh oh! It looks like Tokaido didn't start properly.    ", "")
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
		console.Println(aurora.Sprintf("🤗  How's Tokaido? Run '%s' to share your feedback.", aurora.Bold("tok survey")), "")
	}
}

// InitMessage ...
func InitMessage() {
	goos.InitMessage()
}
