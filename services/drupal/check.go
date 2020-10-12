package drupal

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/blang/semver"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
)

var rgx = regexp.MustCompile("'(.*?)'")
var validDrupalRange = ">=8.5.x"
var checkFailMsg = "\n⚠️  There were some problems detected during the system checks. This won't stop you from running any Tokaido commands, but they may not behave as you were expecting.\n\n"

// CheckLocal ...
func CheckLocal() {
	d := filepath.Join(conf.GetProjectPath(), conf.CoreDrupalFile())

	if !fs.CheckExists(d) {
		fmt.Println("  ×  A Drupal installation was not found")
		console.Printf(checkFailMsg, "")
		return
	}

	fmt.Println("  √  A Drupal installation was found")

	checkDrupalVersion()
}

// CheckContainer ...
func CheckContainer() error {
	nginxPort := docker.LocalPort("nginx", "8443")

	drupalStatus := utils.BashStringCmd(`curl -sko /dev/null -I -w"%{http_code}" https://localhost:` + nginxPort + ` | grep 200`)
	if drupalStatus == "200" || drupalStatus == "302" || drupalStatus == "301" {
		console.Println("😁  Drupal is ready on HTTPS", "√")
		return nil
	}

	console.Println(`😦  Drupal site is not installed or is not working`, "×")
	fmt.Println(`
Tokaido is running but it looks like your Drupal site isn't ready.

You can visit your Drupal site by using the web at
https://` + conf.GetConfig().Tokaido.Project.Name + `.local.tokaido.io:5154.

Your database credentials are:

Hostname: mysql
Username: tokaido
Password: tokaido
Database name: tokaido

It might be easier to use Drush to install your site, which you can do by
connecting to SSH with 'ssh ` + conf.GetConfig().Tokaido.Project.Name + `.tok' and
running 'drush site-install'`)

	return errors.New("Drupal site is not installed")
}

// checkDrupalVersion ...
func checkDrupalVersion() {
	drupalVersion := GetDrupalVersion()
	if drupalVersion == "" {
		fmt.Printf("  ×  Tokaido was unable to determine the Drupal version, it should be %s\n", validDrupalRange)
		console.Printf(checkFailMsg, "")
		return
	}

	drupalSemver, _ := semver.Parse(drupalVersion)
	validRange, _ := semver.ParseRange(validDrupalRange)

	if validRange(drupalSemver) {
		fmt.Printf("  √  The Drupal version (%s) is supported by Tokaido\n\n", drupalVersion)
	} else {
		fmt.Printf("  ×  The Drupal version (%s) is not supported by Tokaido. Drupal version %s has been tested and is working with Tokaido\n", drupalVersion, validDrupalRange)
		console.Printf(checkFailMsg, "")
	}
}

// GetDrupalVersion ...
func GetDrupalVersion() string {
	drupalVersionString := versionStr()
	if drupalVersionString == "" {
		return ""
	}

	drupalVersion := rgx.FindStringSubmatch(drupalVersionString)[0]

	return strings.Replace(drupalVersion, "'", "", 2)
}

// versionStr ...
func versionStr() string {
	f, err := os.Open(filepath.Join(conf.GetProjectPath(), conf.CoreDrupalFile()))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "const VERSION = '") {
			return scanner.Text()
		}
	}

	return ""
}
