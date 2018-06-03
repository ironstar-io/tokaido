package drupal

import (
	"bitbucket.org/ironstar/tokaido-cli/services/drupal/templates"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

var docrootDefault = fs.WorkDir() + "/docroot/sites/default"
var docrootSettingsPath = docrootDefault + "/settings.php"

// CheckSettings ...
func CheckSettings() {
	// detect if file exists
	var _, err = os.Stat(docrootSettingsPath)

	// create file if not exists
	if os.IsNotExist(err) {
		fmt.Printf(`
Could not find a Drupal settings file located at "` + docrootSettingsPath + `", database connection may not work!"
		`)

		return
	}

	if containsTokRef() == false {
		buildTokSettings()
	}
}

func containsTokRef() bool {
	f, err := os.Open(docrootSettingsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	portLine := "/settings.tok.php"
	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	foundRef := false
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), portLine) {
			foundRef = true
		}
	}

	return foundRef
}

func buildTokSettings() {
	confirmCreate := utils.ConfirmationPrompt("Tokaido needs to create database connection settings for your site. May we add the file 'docroot/sites/default/settings.tok.php' and reference it from 'settings.php'?")
	if confirmCreate == false {
		fmt.Println(`
No problem! Please make sure that you manually configure your Drupal site to use the following database connection details:

Hostname: mysql
Username: tokaido
Password: tokaido
Database name: tokaido
		`)
		return
	}

	// appendTokRef()
	createSettingsTok()
}

func appendTokRef() {
	f, err := os.Open(docrootSettingsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	closePHP := "?>"
	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	var buffer bytes.Buffer
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), closePHP) {
			buffer.Write([]byte(drupaltmpl.DocrootSitesDefaultSettings))
		} else {
			buffer.Write([]byte(scanner.Text() + "\n"))
		}
	}

	createSettingsCopy(buffer.String())
	replaceSettings()
}

func createSettingsCopy(body string) {
	var file, err = os.Create(docrootSettingsPath + "-copy")
	if err != nil {
		log.Fatal("Create: ", err)
	}

	_, _ = file.WriteString(body)

	defer file.Close()
}

// replaceSettings - Replace `/docroot/sites/default/settings.php` with `/docroot/sites/default/settings.php-copy`
func replaceSettings() {
	// Remove the original settings file
	os.Remove(docrootSettingsPath)

	// Rename `settings.php-copy` to be the new `settings.php` file
	os.Rename(docrootSettingsPath+"-copy", docrootSettingsPath)
}

func createSettingsTok() {
	var file, err = os.Create(docrootDefault + "/settings.tok.php")
	if err != nil {
		log.Fatal("Create: ", err)
	}

	_, _ = file.WriteString(drupaltmpl.DocrootSitesDefaultSettingsTok)

	defer file.Close()
}
