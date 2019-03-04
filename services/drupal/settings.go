package drupal

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
	drupaltmpl "github.com/ironstar-io/tokaido/services/drupal/templates"
	"github.com/ironstar-io/tokaido/system"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
)

type fileMasks struct {
	DocrootDefault  os.FileMode
	DocrootSettings os.FileMode
}

func sitesDefault() string {
	return filepath.Join(conf.GetRootPath(), "/sites/default")
}

// SettingsPath returns the full system path to the Drupal settings file
func SettingsPath() string {
	return filepath.Join(sitesDefault(), "/settings.php")
}

// SettingsTokPath returns the full system path to the Tokaido settings file
func SettingsTokPath() string {
	return filepath.Join(sitesDefault(), "/settings.tok.php")
}

// CheckSettings checks that Drupal is ready to run in the Tokaido environment
func CheckSettings(checks string) {
	if fs.CheckExists(SettingsPath()) == false {
		fmt.Println(`
Could not find a file located at "` + SettingsPath() + `", database connection may not work!"
		`)
		return
	}

	tokSettingsReferenced := fs.Contains(SettingsPath(), "/settings.tok.php")
	tokSettingsExists := fs.CheckExists(SettingsTokPath())

	if tokSettingsReferenced == false || tokSettingsExists == false {
		if checks != "FORCE" {
			if allowBuildSettings() == false {
				return
			}
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

func checkSettingsExist() {
	_, settingPathErr := os.Stat(SettingsPath())
	if settingPathErr != nil {
		permissionErrMsg(settingPathErr.Error())
		return
	}
	if os.IsNotExist(settingPathErr) {
		fmt.Printf(`
Could not find a Drupal settings file located at "` + SettingsPath() + `", database connection may not work!"
	`)

		return
	}
}

func processFilePerimissions() (fileMasks, error) {
	var emptyStruct fileMasks
	docrootDefaultMask, err := system.GetPermissionsMask(sitesDefault())
	if err != nil {
		return emptyStruct, err
	}

	docrootSettingsMask, err := system.GetPermissionsMask(SettingsPath())
	if err != nil {
		return emptyStruct, err
	}

	defaultMasks := fileMasks{
		DocrootDefault:  docrootDefaultMask,
		DocrootSettings: docrootSettingsMask,
	}

	if fs.Writable(sitesDefault()) == false {
		fmt.Println("\nIt looks like Drupal has been installed before, this operation may need elevated privileges to complete. You may be requested to supply your password.")

		if err := os.Chmod(sitesDefault(), 0770); err != nil {
			return emptyStruct, err
		}

		if err := os.Chmod(SettingsPath(), 0660); err != nil {
			return emptyStruct, err
		}
	}

	return defaultMasks, nil
}

func restoreFilePerimissions(defaultMasks fileMasks) {
	docrootChmodErr := os.Chmod(sitesDefault(), defaultMasks.DocrootDefault)
	if docrootChmodErr != nil {
		fmt.Println(docrootChmodErr)
		return
	}

	settingsChmodErr := os.Chmod(SettingsPath(), defaultMasks.DocrootSettings)
	if settingsChmodErr != nil {
		fmt.Println(settingsChmodErr)
		return
	}
}

func appendTokSettingsRef() {
	f, openErr := os.Open(SettingsPath())
	if openErr != nil {
		fmt.Println(openErr)
		return
	}

	defer f.Close()

	dv := conf.GetConfig().Drupal.Majorversion
	var settingsBody []byte
	if dv == "7" {
		settingsBody = drupaltmpl.SettingsD7Append
	} else if dv == "8" {
		settingsBody = drupaltmpl.SettingsD8Append
	} else {
		log.Fatalf("Could not apply Drupal settings")
	}

	closePHP := "?>"
	var closeTagFound = false
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), closePHP) {
			closeTagFound = true
			buffer.Write(settingsBody)
		} else {
			buffer.Write([]byte(scanner.Text() + "\n"))
		}
	}

	if closeTagFound == false {
		buffer.Write(settingsBody)
	}

	fs.Replace(SettingsPath(), buffer.Bytes())
}

func createSettingsTok() {
	c := conf.GetConfig()
	dv := c.Drupal.Majorversion

	salt, err := GenerateRandomHashSalt()
	if err != nil {
		log.Fatalf("Could not create Drupal hash salt %v", err)
	}

	dt := drupaltmpl.Settings{
		HashSalt:          salt,
		ProjectName:       c.Tokaido.Project.Name,
		FilePublicPath:    c.Drupal.FilePublicPath,
		FilePrivatePath:   c.Drupal.FilePrivatePath,
		FileTemporaryPath: c.Drupal.FileTemporaryPath,
	}

	var settingsTokBody []byte
	if dv == "7" {
		settingsTokBody = dt.SettingsD7Tok()
	} else if dv == "8" {
		settingsTokBody = dt.SettingsD8Tok()
	} else {
		log.Fatalf("Could not add Tokaido settings file")
	}

	fs.TouchByteArray(SettingsTokPath(), settingsTokBody)
}

func allowBuildSettings() bool {
	confirmation := utils.ConfirmationPrompt(`
Tokaido can automatically add database connection settings to your Drupal site.

If you prefer to make these changes yourself, choose no and we'll show you
the database connection settings you'll need.

Let Tokaido automatically configure your database connection settings?`, "y")

	if confirmation == false {
		fmt.Println(`
No problem! Please make sure that you manually configure your Drupal site to use
the following database connection details:

Hostname: mysql
Username: tokaido
Password: tokaido
Database name: tokaido

Please see the Tokaido documentation at https://tokaido.io/docs/
for more information on setting up your Tokaido environment.

		`)
	}

	return confirmation
}

func permissionErrMsg(errString string) {
	fmt.Printf(`
%s
Please make sure that you manually configure your Drupal site to use the
following database connection details:

Hostname: mysql
Username: tokaido
Password: tokaido
Database name: tokaido
	`, errString)
}
