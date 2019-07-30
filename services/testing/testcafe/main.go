package testcafe

import (
	"fmt"
	"path/filepath"

	"github.com/ironstar-io/tokaido/conf"
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
	}

	recreateTokaidoTestDB()
	clearDrupalCache()
	copyDefaultDB()
	// Clone prod DB
	addTestCafeUsers()

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

func baseRepoExists() bool {
	env := filepath.Join(testCafeLocalPath(), "tests")
	if fs.CheckExists(env) == false {
		return false
	}

	return true
}

// Download the base testcafe repo and copy into the Drupal project
func pullBaseRepo() {
	utils.CheckCmdHard("curl")

	pr := conf.GetProjectPath()
	dr := conf.GetConfig().Drupal.Path
	tp := filepath.Join(pr, dr, ".tok")
	fs.Mkdir(tp)

	utils.StdoutCmd("wget", "https://github.com/ironstar-io/tokaido-testcafe/archive/0.0.4.tar.gz")
	utils.StdoutCmd("tar", "xzf", "0.0.4.tar.gz")
	utils.StdoutCmd("mv", "tokaido-testcafe-0.0.4", filepath.Join(tp, "testcafe"))
	utils.StdoutCmd("rm", "0.0.4.tar.gz")
}

func recreateTokaidoTestDB() {
	ssh.ConnectCommand([]string{"drush", "sql:create", "--database=test", "-y", "-d"})
}

func clearDrupalCache() {
	ssh.ConnectCommand([]string{"drush", "cache-rebuild", "-d"})
}

func copyDefaultDB() {
	ssh.ConnectCommand([]string{"drush", "sql:dump", "--result-file=toktest-dump.sql", "-d"})
	ssh.ConnectCommand([]string{"drush", "sql:cli", "--database=test", "<", "toktest-dump.sql", "-d"})
}

func addTestCafeUsers() {
	ssh.ConnectCommand([]string{"drush", "user:create", "testcafe_user", `--password="testcafe_user"`, `--mail="testcafe_user@localhost"`, "-d"})
	ssh.ConnectCommand([]string{"drush", "user:create", "testcafe_admin", `--password="testcafe_admin"`, `--mail="testcafe_admin@localhost"`, "-d"})
	ssh.ConnectCommand([]string{"drush", "user:role:add", `"administrator"`, "testcafe_admin", "-d"})
	ssh.ConnectCommand([]string{"drush", "user:create", "testcafe_editor", `--password="testcafe_editor"`, `--mail="testcafe_editor@localhost"`, "-d"})
	ssh.ConnectCommand([]string{"drush", "user:role:add", `"editor"`, "testcafe_editor", "-d"})
}

func npmCI() {
	nm := filepath.Join(testCafeLocalPath(), "node_modules")
	if fs.CheckExists(nm) == false {
		fmt.Println(aurora.Cyan("🦉  Installing TestCafe dependencies"))
		ssh.StreamConnectCommand([]string{"cd", testCafeContainerPath(), "&&", "npm", "ci"})
	}
}

func npmHeadlessTestCafe() {
	utils.StdoutStreamCmd("docker-compose", "-f", filepath.Join(conf.GetProjectPath(), "/docker-compose.tok.yml"), "exec", "-T", "testcafe", "npm", "run", "test")
}

func checkCompatibility() {
	goos.CheckDependencies()
}
