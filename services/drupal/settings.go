package drupal

import (
	"bitbucket.org/ironstar/tokaido-cli/services/drupal/templates"
	"bitbucket.org/ironstar/tokaido-cli/system"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

type fileMasks struct {
	DocrootDefault  os.FileMode
	DocrootSettings os.FileMode
}

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

	defaultMasks, maskErr := processFilePerimissions()
	if maskErr != nil {
		permissionErrMsg(maskErr.Error())
		return
	}

	if containsTokRef() == false {
		buildTokSettings()
	}

	restoreFilePerimissions(defaultMasks)
}

func permissionErrMsg(path string) {
	fmt.Printf(`
%s
Please make sure that you manually configure your Drupal site to use the following database connection details:

Hostname: mysql
Username: tokaido
Password: tokaido
Database name: tokaido
	`, path)
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
	confirmCreate := utils.ConfirmationPrompt(`
Tokaido needs to create database connection settings for your site. May we add the file 'docroot/sites/default/settings.tok.php' and reference it from 'settings.php'?`)

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

	appendTokRef()
	createSettingsTok()
}

func processFilePerimissions() (fileMasks, error) {
	var emptyStruct fileMasks
	docrootDefaultMask, dirMaskErr := system.GetPermissionsMask(docrootDefault)
	if dirMaskErr != nil {
		return emptyStruct, dirMaskErr
	}

	docrootSettingsMask, settingsMaskErr := system.GetPermissionsMask(docrootSettingsPath)
	if settingsMaskErr != nil {
		return emptyStruct, settingsMaskErr
	}

	defaultMasks := fileMasks{
		DocrootDefault:  docrootDefaultMask,
		DocrootSettings: docrootSettingsMask,
	}

	if fs.Writable(docrootDefault) == false {
		fmt.Println("\nIt looks like Drupal has been installed before, this operation may need elevated privileges to complete. You may be requested to supply your password.")

		os.Chmod(docrootDefault, 0770)
		os.Chmod(docrootSettingsPath, 0660)
	}

	return defaultMasks, nil
}

func restoreFilePerimissions(defaultMasks fileMasks) {
	os.Chmod(docrootDefault, defaultMasks.DocrootDefault)
	os.Chmod(docrootSettingsPath, defaultMasks.DocrootSettings)
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
			buffer.Write([]byte(drupaltmpl.SettingsAppend))
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

	_, _ = file.WriteString(drupaltmpl.SettingsTok)

	defer file.Close()
}
