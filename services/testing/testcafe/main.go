package testcafe

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/database"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/testing/testcafe/goos"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/utils"

	"github.com/logrusorgru/aurora"
)

// RunDrupalTests - Run Drupal TestCafe tests. See: https://github.com/drupal/drupal/blob/8.6.x/core/tests/README.md#testcafe-tests
func RunDrupalTests(useExistingDB bool) error {
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

	db, err := database.ConnectRoot("tokaido")
	if err != nil {
		fmt.Println(aurora.Red("Error: Could not connect to the database container failed. Have you run `tok up`?"))
		utils.DebugString(err.Error())
		return err
	}

	testDBReady := isTestDBAccessible(db)

	if useExistingDB == false || testDBReady == false {
		err = recreateTokaidoTestDB(db)
		if err != nil {
			return err
		}

		clearDrupalCache()
		copyDefaultDB()
		addTestCafeUsers()
	}

	npmCI()
	npmHeadlessTestCafe()

	return nil
}

// Pull the testcafe folder absolute path
func testCafeLocalPath() string {
	pr := conf.GetProjectPath()

	return filepath.Join(pr, ".tok", "testcafe")
}

func testCafeContainerPath() string {
	return filepath.Join("/tokaido", "site", ".tok", "testcafe")
}

func baseRepoExists() bool {
	p := filepath.Join(testCafeLocalPath(), "tests")
	if fs.CheckExists(p) == false {
		return false
	}

	return true
}

// Download the base testcafe repo and copy into the Drupal project
func pullBaseRepo() {
	fmt.Println("🗃️   Fetching the TestCafe base repo and placing in .tok/testcafe")

	utils.CheckCmdHard("wget")

	pr := conf.GetProjectPath()
	tp := filepath.Join(pr, ".tok")
	fs.Mkdir(tp)

	utils.StdoutStreamCmdDebugContext(filepath.Join(tp), "wget", "https://github.com/ironstar-io/tokaido-testcafe/archive/0.0.8.tar.gz")
	utils.StdoutStreamCmdDebugContext(filepath.Join(tp), "tar", "xzf", "0.0.8.tar.gz")
	fs.Remove(filepath.Join(tp, "0.0.8.tar.gz"))
	fs.Remove(filepath.Join(tp, "testcafe"))
	fs.Rename(filepath.Join(tp, "tokaido-testcafe-0.0.8"), filepath.Join(tp, "testcafe"))
}

// isTestDBAccessible returns a boolean marking whether the 'tokaido_test' DB can be accessed or not
func isTestDBAccessible(db *sql.DB) bool {
	// Check to see if the database already exists
	db, err := database.Connect("tokaido_test")
	if err != nil {
		fmt.Println(aurora.Red("Error: Could not connect to the database container failed. Have you run `tok up`?"))
		utils.DebugString(err.Error())

		return false
	}

	err = db.Ping()
	if err != nil {
		// The database connection was opened but failed, so the provided database doesn't exist
		return false
	}

	return true
}

func recreateTokaidoTestDB(db *sql.DB) error {
	db, err := database.ConnectRoot("tokaido")
	if err != nil {
		fmt.Println(aurora.Red("Error: Could not connect to the database container failed. Have you run `tok up`?"))
		utils.DebugString(err.Error())

		return err
	}

	utils.DebugString("Creating tokaido_test database")
	_, err = db.Exec("DROP DATABASE IF EXISTS tokaido_test;")
	if err != nil {
		fmt.Println(aurora.Red("Error: Could not drop the tokaido_test database: "), err.Error())
		utils.DebugString(err.Error())

		return err
	}

	// Create the database
	utils.DebugString("Creating tokaido_test database")
	_, err = db.Exec("CREATE DATABASE tokaido_test;")
	if err != nil {
		fmt.Println(aurora.Red("Error: Could not create the tokaido_test database: "), err.Error())
		utils.DebugString(err.Error())

		return err
	}

	// Grant 'tokaido' user access
	utils.DebugString("Creating full access to 'tokaido_test' database to 'tokaido' user")
	_, err = db.Exec("GRANT ALL ON tokaido_test.* TO 'tokaido'@'%';")
	if err != nil {
		fmt.Println(aurora.Red("Error: Could not grant access to the tokaido_test database: "), err.Error())
		utils.DebugString(err.Error())

		return err
	}

	// Flush
	utils.DebugString("Flushing privileges")
	_, err = db.Exec("FLUSH PRIVILEGES;")
	if err != nil {
		fmt.Println(aurora.Red("Error: Could not successfully flush privileges: "), err.Error())
		utils.DebugString(err.Error())
		return err
	}

	time.Sleep(1 * time.Second)

	// See if the creation was successful
	success := isTestDBAccessible(db)
	if !success {
		fmt.Println(aurora.Red("Error: Tokaido was not able to successfully create the 'tests' database"))
		return err
	}

	return nil
}

func clearDrupalCache() {
	fmt.Println("🗑️   Clearing the Drupal cache")

	ssh.ConnectCommand([]string{"drush", "cache-rebuild", "-d"})
}

func copyDefaultDB() {
	fmt.Println("🖨️   Creating a test copy of the database")
	outputPath := filepath.Join(testCafeContainerPath(), "toktest-dump.sql")

	ssh.ConnectCommand([]string{"drush", "sql-dump", "--result-file=" + outputPath, "-d"})
	ssh.ConnectCommand([]string{"drush", "sql-cli", "--database=test", "<", outputPath, "-d"})
}

func addTestCafeUsers() {
	fmt.Println("👩‍💻  Adding temporary Drupal user and administrator")

	ssh.ConnectCommand([]string{"drush", "user-create", "testcafe_user", `--password="testcafe_user"`, `--mail="testcafe_user@localhost"`, "-d"})
	ssh.ConnectCommand([]string{"drush", "user-create", "testcafe_admin", `--password="testcafe_admin"`, `--mail="testcafe_admin@localhost"`, "-d"})
	ssh.ConnectCommand([]string{"drush", "user-add-role", `"administrator"`, "testcafe_admin", "-d"})
}

func npmCI() {
	nm := filepath.Join(testCafeLocalPath(), "node_modules")
	if fs.CheckExists(nm) == false {
		fmt.Println(aurora.Cyan("🦉  Installing TestCafe dependencies"))
		ssh.StreamConnectCommand([]string{"cd", testCafeContainerPath(), "&&", "npm", "ci"})
	}
}

func npmHeadlessTestCafe() {
	// Restart the testcafe container to ensure the workdir and volume is bound
	docker.ComposeStdout("restart", "testcafe")

	utils.StdoutStreamCmd("docker-compose", "-f", filepath.Join(conf.GetProjectPath(), "/docker-compose.tok.yml"), "exec", "-T", "testcafe", "npm", "run", "test")
}

func checkCompatibility() {
	goos.CheckDependencies()
}
