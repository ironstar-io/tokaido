package nightwatch

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/database"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/proxy"
	"github.com/ironstar-io/tokaido/services/testing/nightwatch/goos"
	"github.com/ironstar-io/tokaido/services/tok"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/utils"

	. "github.com/logrusorgru/aurora"
)

// RunDrupalTests - Run Drupal Nightwatch tests. See: https://github.com/drupal/drupal/blob/8.6.x/core/tests/README.md#nightwatch-tests
func RunDrupalTests() error {
	fmt.Println(Blue("Please note 'tok test' is a beta feature and under active development"))
	fmt.Println(Blue("We'd to hear what you think - you can use `tok survey` to send us feedback"))
	fmt.Println(Blue("If you'd like to contribute by discussing test frameworks and included tests"))
	fmt.Println(Blue("please find us on the #tokaido channel of the official Drupal Slack"))
	fmt.Println()

	// Create the 'tests' database if it doesn't exist
	dbOk := checkMysqlTestsDatabase()
	if !dbOk {
		createMysqlTestsDatabase()
	}

	err := CheckLocalDrupal()
	if err != nil {
		console.Println("⚠️   "+err.Error(), "")
		return err
	}

	err = checkCompatibility()
	if err != nil {
		return err
	}

	enableChromedriver()

	checkPHPUnit()
	yarnInstall()

	configureEnvFile()

	checkNightwatchConfig()

	yarnTestNightwatch()

	return nil
}

func calcSiteURL() string {
	if proxy.CheckProxyUp() == true {
		return proxy.GetProxyURL()
	}

	return "https://localhost:" + docker.LocalPort("haproxy", "8443")
}

// Copy and amend .env file if it doesn't already exist
func configureEnvFile() {
	pr := conf.GetProjectPath()
	dr := conf.GetConfig().Drupal.Path
	env := filepath.Join(pr, dr, "core", ".env")
	if fs.CheckExists(env) == false {
		fmt.Println(Cyan(fmt.Sprintf("📋  Adding Nightwatch config to %s", env)))
		generateEnvFile(env)
	}
}

func yarnInstall() error {
	cp := conf.CoreDrupal8Path()

	fmt.Println(Cyan("🦉  Ensuring Nightwatch dependencies are installed"))
	ssh.StreamConnectCommand([]string{"cd", cp, "&&", "yarn", "install"})
	return nil
}

func yarnTestNightwatch() {
	cp := conf.CoreDrupal8Path()

	fmt.Println(Cyan("👩‍💻  Tokaido is starting a Nightwatch test run with the command `yarn test:nightwatch`"))
	fmt.Println(Yellow(Sprintf("    If you want to run this command directly, you need run it from %s", Bold("inside the Tokaido SSH container"))))

	ssh.StreamConnectCommand([]string{"cd", cp, "&&", "yarn", "test:nightwatch"})

}

func checkCompatibility() error {
	goos.CheckDependencies()
	return nil
}

func enableChromedriver() {
	if !conf.GetConfig().Services.Chromedriver.Enabled {
		conf.SetConfigValueByArgs([]string{"services", "chromedriver", "enabled", "true"}, "project")
	}

	if !docker.StatusCheck("chromedriver", conf.GetConfig().Tokaido.Project.Name) {
		fmt.Println(Cyan("👾  Restarting Tokaido with Chromedriver enabled"))
		tok.Init(true, false)

	}
}

// checkMysqlTestsDatabase spawns a 'tests' databases in the Tokaido MySQL instance
func checkMysqlTestsDatabase() bool {
	// Check to see if the database already exists
	db, err := database.Connect("tests")
	if err != nil {
		fmt.Println(Red("Error: Could not connect to the database container failed. Have you run `tok up`?"))
		utils.DebugString(err.Error())
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		// The database connection was opened but failed, so the provided database doesn't exist
		return false
	}

	return true
}

func createMysqlTestsDatabase() {
	fmt.Println(Cyan("    Creating a new 'tests' database"))

	db, err := database.ConnectRoot("tokaido")
	if err != nil {
		fmt.Println(Red("Error: Could not connect to the database container failed. Have you run `tok up`?"))
		utils.DebugString(err.Error())
		os.Exit(1)
	}

	// Create the database
	utils.DebugString("Creating test database")
	_, err = db.Exec("CREATE DATABASE tests;")
	if err != nil {
		fmt.Println(Red("Error: Could not create the tests database: "), err.Error())
		utils.DebugString(err.Error())
		os.Exit(1)
	}

	// Grant 'tokaido' user access
	utils.DebugString("Creating full access to 'tests' database to 'tokaido' user")
	_, err = db.Exec("GRANT ALL ON tests.* TO 'tokaido'@'%';")
	if err != nil {
		fmt.Println(Red("Error: Could not grant access to the tests database: "), err.Error())
		utils.DebugString(err.Error())
		os.Exit(1)
	}

	// Flush
	utils.DebugString("Flushing privileges")
	_, err = db.Exec("FLUSH PRIVILEGES;")
	if err != nil {
		fmt.Println(Red("Error: Could not successfully flush privileges: "), err.Error())
		utils.DebugString(err.Error())
		os.Exit(1)
	}

	time.Sleep(1 * time.Second)

	// See if the creation was successful
	success := checkMysqlTestsDatabase()
	if !success {
		fmt.Println(Red("Error: Tokaido was not able to successfully create the 'tests' database"))
		os.Exit(1)
	}
}
