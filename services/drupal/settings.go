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
var settingsPath = docrootDefault + "/settings.php"
var settingsTokPath = docrootDefault + "/settings.tok.php"

// CheckSettings ...
func CheckSettings() {
	if checkExists(settingsPath) == false {
		fmt.Printf(`
Could not find a file located at "` + settingsPath + `", database connection may not work!"
		`)
		return
	}

	tokSettingsReferenced := fs.Contains(settingsPath, "/settings.tok.php")
	tokSettingsExists := checkExists(settingsTokPath)

	if tokSettingsReferenced == false || tokSettingsExists == false {
		if allowBuildSettings() == false {
			return
		}
	} else {
		return
	}

	defaultMasks, err := processFilePerimissions()
	if err != nil {
		permissionErrMsg(err.Error())
		return
	}

	if tokSettingsReferenced == false {
		appendTokSettingsRef()
	}

	if tokSettingsExists == false {
		createSettingsTok()
	}

	restoreFilePerimissions(defaultMasks)
}

func checkExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func checkSettingsExist() {
	_, settingPathErr := os.Stat(settingsPath)
	if settingPathErr != nil {
		permissionErrMsg(settingPathErr.Error())
		return
	}
	if os.IsNotExist(settingPathErr) {
		fmt.Printf(`
Could not find a Drupal settings file located at "` + settingsPath + `", database connection may not work!"
	`)

		return
	}
}

func processFilePerimissions() (fileMasks, error) {
	var emptyStruct fileMasks
	docrootDefaultMask, err := system.GetPermissionsMask(docrootDefault)
	if err != nil {
		return emptyStruct, err
	}

	docrootSettingsMask, err := system.GetPermissionsMask(settingsPath)
	if err != nil {
		return emptyStruct, err
	}

	defaultMasks := fileMasks{
		DocrootDefault:  docrootDefaultMask,
		DocrootSettings: docrootSettingsMask,
	}

	if fs.Writable(docrootDefault) == false {
		fmt.Println("\nIt looks like Drupal has been installed before, this operation may need elevated privileges to complete. You may be requested to supply your password.")

		if err := os.Chmod(docrootDefault, 0770); err != nil {
			return emptyStruct, err
		}

		if err := os.Chmod(settingsPath, 0660); err != nil {
			return emptyStruct, err
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

	settingsChmodErr := os.Chmod(settingsPath, defaultMasks.DocrootSettings)
	if settingsChmodErr != nil {
		fmt.Println(settingsChmodErr)
		return
	}
}

func appendTokSettingsRef() {
	f, openErr := os.Open(settingsPath)
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
			buffer.Write(drupaltmpl.SettingsAppend)
		} else {
			buffer.Write([]byte(scanner.Text() + "\n"))
		}
	}

	if closeTagFound == false {
		buffer.Write(drupaltmpl.SettingsAppend)
	}

	fs.Replace(settingsPath, buffer.Bytes())
}

func createSettingsTok() {
	fs.TouchByteArray(settingsTokPath, drupaltmpl.SettingsTok)
}

func allowBuildSettings() bool {
	confirmation := utils.ConfirmationPrompt(`
Tokaido can now create database connection settings for your site.
Should Tokaido add the file 'docroot/sites/default/settings.tok.php'
and reference it from 'settings.php'?

If you prefer not to do this automatically, we'll show you database connection
settings so that you can configure this manually.`)

	if confirmation == false {
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
	}

	return confirmation
}

func permissionErrMsg(errString string) {
	fmt.Printf(`
%s
Please make sure that you manually configure your Drupal site to use the following database connection details:

Hostname: mysql
Username: tokaido
Password: tokaido
Database name: tokaido
	`, errString)
}
