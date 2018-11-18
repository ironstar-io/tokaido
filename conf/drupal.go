package conf

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ironstar-io/tokaido/system/console"
	"github.com/manifoldco/promptui"
)

// GetRootPath ...
func GetRootPath() string {
	c := GetConfig()
	dp := c.Drupal.Path
	if dp != "" {
		return filepath.Join(c.Tokaido.Project.Path, dp)
	}

	log.Fatalf("Drupal path setting is missing.")
	return ""
}

// SetDrupalConfig if there is no config already applied
func SetDrupalConfig(drupalType string) {
	if drupalType == "DEFAULT" {
		CreateOrReplaceDrupalConfig("/web", "8")
		return
	}

	p := GetConfig().Drupal.Path
	v := GetConfig().Drupal.Majorversion

	if (p == "") || (v == "") {
		p, v = DetectDrupalSettings(GetConfig().Tokaido.Project.Path)
	}

	CreateOrReplaceDrupalConfig(p, v)
}

// DetectDrupalSettings - Return drupal root and version
func DetectDrupalSettings(projectRoot string) (string, string) {
	var dp string
	var dv string
	d7 := filepath.Join("includes", "bootstrap.inc")
	d8 := filepath.Join("core", "lib", "Drupal.php")
	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.Contains(path, d7) == true {
			f, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatal("Could not read bootstrap file: ", err)
			}
			s := string(f)
			// there will be bootstrap.inc files in a Drupal site, make sure this is _the_ bootstrap.inc from Drupal core
			if strings.Contains(s, "'VERSION', '7.") {
				console.Println("🚂  Found a Drupal 7 site", "")
				// Strip the Drupal component from the full path
				dp = strings.Replace(path, d7, "", -1)
				// Strip the work dir from the remainder
				dp = strings.Replace(dp, projectRoot, "", -1)
				// Strip slashes
				dp = strings.Replace(dp, "/", "", -1)
				dv = "7"
				return io.EOF
			}
		}
		if strings.Contains(path, d8) == true {
			console.Println("🚇  Found a Drupal 8 site", "")
			// Strip the Drupal component from the full path
			dp = strings.Replace(path, d8, "", -1)
			// Strip the work dir from the remainder
			dp = strings.Replace(dp, projectRoot, "", -1)
			// Strip slashes
			dp = strings.Replace(dp, "/", "", -1)
			dv = "8"
			return io.EOF
		}
		return nil
	})
	if err != io.EOF {
		fmt.Println("\n🤷‍  Tokaido could not auto-detect your Drupal installation. You'll need to tell us about it.")
		dp, dv = manualDrupalSettings()
	}

	return dp, dv
}

func manualDrupalSettings() (string, string) {
	vPrompt := promptui.Select{
		Label: "What major version of Drupal are you running?",
		Items: []string{"Drupal 8", "Drupal 7"},
	}

	_, dv, vErr := vPrompt.Run()

	if vErr != nil {
		log.Fatalf("Prompt failed %v\n", vErr)

	}

	pPrompt := promptui.Select{
		Label: "Where is your Drupal root?",
		Items: []string{"/docroot", "/app", "/web", "other"},
	}

	_, dp, pErr := pPrompt.Run()

	if dp == "other" {
		var err error
		prompt := promptui.Prompt{
			Label: "Please enter the name of your Drupal root directory",
		}

		dp, err = prompt.Run()

		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}
	}

	if pErr != nil {
		log.Fatalf("Prompt failed %v\n", pErr)
	}

	// Convert the human-friendly values to something Tokaido can use
	var version string
	switch dv {
	case "Drupal 8":
		version = "8"
	case "Drupal 7":
		version = "7"
	}

	path := strings.Replace(dp, "/", "", -1)

	return path, version
}

// CoreDrupalFile - Return the core drupal file for the users' installation
func CoreDrupalFile() string {
	rp := GetConfig().Drupal.Path
	dv := GetConfig().Drupal.Majorversion

	var path string
	switch dv {
	case "7":
		path = filepath.Join(rp, "includes", "bootstrap.inc")
	case "8":
		path = filepath.Join(rp, "core", "lib", "Drupal.php")
	default:
		log.Fatal("Unknown or unspecified Drupal majorVersion in config file")
	}

	return path
}

// CoreDrupal8Path - Return the core Drupal 8 path for the users' installation
func CoreDrupal8Path() string {
	c := GetConfig()
	tp := c.Tokaido.Project.Path
	rp := c.Drupal.Path

	return filepath.Join(tp, rp, "core")
}

// GetRootDir - Return the drupal root folder name without workdir
func GetRootDir() string {
	dr := GetRootPath()

	return filepath.Base(dr)
}
