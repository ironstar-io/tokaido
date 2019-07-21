package testcafe

import (
	"fmt"
	"path/filepath"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/proxy"
	"github.com/ironstar-io/tokaido/services/testing/testcafe/goos"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/utils"

	"github.com/logrusorgru/aurora"
)

// RunDrupalTests - Run Drupal TestCafe tests. See: https://github.com/drupal/drupal/blob/8.6.x/core/tests/README.md#testcafe-tests
func RunDrupalTests() error {
	fmt.Println(aurora.Blue("Please note 'tok test' is a beta feature and under active development"))
	fmt.Println(aurora.Blue("We'd to hear what you think - you can use `tok survey` to send us feedback"))
	fmt.Println(aurora.Blue("If you'd like to contribute by discussing test frameworks and included tests"))
	fmt.Println(aurora.Blue("please find us on the #tokaido channel of the official Drupal Slack"))
	fmt.Println()

	err := CheckLocalDrupal()
	if err != nil {
		console.Println("⚠️   "+err.Error(), "")
		return err
	}

	checkCompatibility()

	br := baseRepoExists()
	if br == false {
		pullBaseRepo()
		configureDefaultJSON()
		addTestCafeUsers()
	}

	npmCI()
	npmHeadlessTestCafe()

	return nil
}

// Pull the testcafe folder absolute path
func testCafeLocalPath() string {
	pr := conf.GetProjectPath()
	dr := conf.GetConfig().Drupal.Path

	return filepath.Join(pr, dr, ".tok", "testcafe")
}

func testCafeContainerPath() string {
	dr := conf.GetConfig().Drupal.Path

	return filepath.Join("/tokaido", "site", dr, ".tok", "testcafe")
}

// Copy and amend .env file if it doesn't already exist
func baseRepoExists() bool {
	env := filepath.Join(testCafeLocalPath(), "config", "default.json")
	if fs.CheckExists(env) == false {
		return false
	}

	return true
}

// Clone the base testcafe repo and copy into the Drupal project
func pullBaseRepo() {
	utils.CheckCmdHard("curl")

	pr := conf.GetProjectPath()
	dr := conf.GetConfig().Drupal.Path
	tp := filepath.Join(pr, dr, ".tok")
	fs.Mkdir(tp)

	utils.StdoutCmd("wget", "https://github.com/ironstar-io/drupal-testcafe/releases/download/v0.0.1-alpha/drupal-testcafe-0.0.1-alpha.tar.gz")
	utils.StdoutCmd("tar", "xzf", "drupal-testcafe-0.0.1-alpha.tar.gz")
	utils.StdoutCmd("mv", "drupal-testcafe-0.0.1-alpha", filepath.Join(tp, "testcafe"))
	utils.StdoutCmd("rm", "drupal-testcafe-0.0.1-alpha.tar.gz")
}

// Copy and amend the default.json config file if it doesn't already exist
func configureDefaultJSON() {
	env := filepath.Join(testCafeLocalPath(), "config", "default.json")
	if fs.CheckExists(env) == false {
		url := calcSiteURL()
		fmt.Println(aurora.Cyan(fmt.Sprintf("📋  Adding TestCafe config to %s for %s", env, url)))

		generateEnvFile(env, url)
	}
}

func addTestCafeUsers() {
	// TODO - Check if user exists first, or will error

	// ssh.ConnectCommand([]string{"drush", "user-create", "testcafe_user", `--password="testcafe_user"`, `--mail="testcafe_user@localhost"`})
	// ssh.ConnectCommand([]string{"drush", "user-create", "testcafe_admin", `--password="testcafe_admin"`, `--mail="testcafe_admin@localhost"`})
	// ssh.ConnectCommand([]string{"drush", "user-add-role", `"administrator"`, "testcafe_admin"})
	// ssh.ConnectCommand([]string{"drush", "user-create", "testcafe_editor", `--password="testcafe_editor"`, `--mail="testcafe_editor@localhost"`})
	// ssh.ConnectCommand([]string{"drush", "user-add-role", `"editor"`, "testcafe_editor"})
}

func calcSiteURL() string {
	if proxy.CheckProxyUp() == true {
		return proxy.GetProxyURL()
	}

	return "https://localhost:" + docker.LocalPort("haproxy", "8443")
}

func npmCI() {
	fmt.Println(aurora.Cyan("🦉  Ensuring TestCafe dependencies are installed"))
	ssh.StreamConnectCommand([]string{"cd", testCafeContainerPath(), "&&", "npm", "ci"})
}

func npmHeadlessTestCafe() {
	fmt.Println(aurora.Cyan("👩‍💻  Tokaido is starting a TestCafe test run with the command `npm run tests:headless:all`"))
	fmt.Println(aurora.Yellow(aurora.Sprintf("    If you want to run this command directly, you need run it from %s", aurora.Bold("inside the Tokaido SSH container"))))

	ssh.StreamConnectCommand([]string{"cd", testCafeContainerPath(), "&&", "npm", "run", "tests:headless:all"})
}

func checkCompatibility() {
	goos.CheckDependencies()
}
