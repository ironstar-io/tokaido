package drupal

import (
	"bitbucket.org/ironstar/tokaido-cli/services/drupal/templates"
	"bitbucket.org/ironstar/tokaido-cli/system"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"bufio"
	"bytes"
	"fmt"
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
	var _, settingPathErr = os.Stat(docrootSettingsPath)
	if settingPathErr != nil {
		permissionErrMsg(settingPathErr.Error())
		return
	}

	// create file if not exists
	if os.IsNotExist(settingPathErr) {
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
	f, openErr := os.Open(docrootSettingsPath)
	if openErr != nil {
		fmt.Println(openErr)
		return true
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
Tokaido can now create database connection settings for your site.
Should Tokaido add the file 'docroot/sites/default/settings.tok.php'
and reference it from 'settings.php'?

If you prefer not to do this automatically, we'll show you database connection
settings so that you can configure this manually.`)

	if confirmCreate == false {
		fmt.Println(`
No problem! Please make sure that you manually configure your Drupal site to use
the following database connection details:

Hostname: mysql
Username: tokaido
Password: tokaido
Database name: tokaido

Please see the Tokaido environments guide at https://docs.tokaido.io/environments
for more information on setting up your Tokaido environment.

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

		docrootChmodErr := os.Chmod(docrootDefault, 0770)
		if docrootChmodErr != nil {
			return emptyStruct, docrootChmodErr
		}

		settingsChmodErr := os.Chmod(docrootSettingsPath, 0660)
		if settingsChmodErr != nil {
			return emptyStruct, settingsChmodErr
		}
	}

	return defaultMasks, nil
}

func restoreFilePerimissions(defaultMasks fileMasks) {
	docrootChmodErr := os.Chmod(docrootDefault, defaultMasks.DocrootDefault)
	if docrootChmodErr != nil {
		fmt.Println(docrootChmodErr)
		return
	}

	settingsChmodErr := os.Chmod(docrootSettingsPath, defaultMasks.DocrootSettings)
	if settingsChmodErr != nil {
		fmt.Println(settingsChmodErr)
		return
	}
}

func appendTokRef() {
	f, openErr := os.Open(docrootSettingsPath)
	if openErr != nil {
		fmt.Println(openErr)
		return
	}

	defer f.Close()

	closePHP := "?>"
	var closeTagFound = false
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), closePHP) {
			closeTagFound = true
			buffer.Write([]byte(drupaltmpl.SettingsAppend))
		} else {
			buffer.Write([]byte(scanner.Text() + "\n"))
		}
	}

	if closeTagFound == false {
		buffer.Write([]byte(drupaltmpl.SettingsAppend))
	}

	createSettingsCopy(buffer.String())
	replaceSettings()
}

func createSettingsCopy(body string) {
	var file, createErr = os.Create(docrootSettingsPath + "-copy")
	if createErr != nil {
		fmt.Println(createErr)
		return
	}

	_, writeErr := file.WriteString(body)
	if writeErr != nil {
		fmt.Println(writeErr)
		return
	}

	defer file.Close()
}

// replaceSettings - Replace `/docroot/sites/default/settings.php` with `/docroot/sites/default/settings.php-copy`
func replaceSettings() {
	// Remove the original settings file
	removeErr := os.Remove(docrootSettingsPath)
	if removeErr != nil {
		fmt.Println(removeErr)
		return
	}

	// Rename `settings.php-copy` to be the new `settings.php` file
	renameErr := os.Rename(docrootSettingsPath+"-copy", docrootSettingsPath)
	if renameErr != nil {
		fmt.Println(renameErr)
		return
	}
}

func createSettingsTok() {
	var file, createErr = os.Create(docrootDefault + "/settings.tok.php")
	if createErr != nil {
		fmt.Println(createErr)
		return
	}

	_, writeErr := file.WriteString(drupaltmpl.SettingsTok)
	if writeErr != nil {
		fmt.Println(writeErr)
		return
	}

	defer file.Close()
}
