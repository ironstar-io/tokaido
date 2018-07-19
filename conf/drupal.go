package conf

import (
	"fmt"
	"io/ioutil"

	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetRootPath ...
func GetRootPath() string {
	wd := fs.WorkDir()
	c := GetConfig().Drupal.Path
	if c != "" {
		return filepath.Join(wd, c)
	}

	rootPath, version := scanForCoreDrupal()

	CreateOrReplaceDrupalVars(strings.Replace(rootPath, wd, "", -1), version)

	return rootPath
}

func scanForCoreDrupal() (string, string) {
	wd := fs.WorkDir()
	var dp string
	var dv string
	d7 := filepath.Join("includes", "bootstrap.inc")
	d8 := filepath.Join("core", "lib", "Drupal.php")
	err := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.Contains(path, d7) == true {
			f, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatal("Could not read bootstrap file: ", err)
			}
			s := string(f)
			// bootstrap.inc is a pretty common path, make sure this is _the_ bootstrap.inc from Drupal
			if strings.Contains(s, "'VERSION', '7.") {
				fmt.Printf("🚂  Found a Drupal 7 site")
				dp = strings.Replace(path, d7, "", -1)
				dv = "7"
				return io.EOF
			}
		}
		if strings.Contains(path, d8) == true {
			fmt.Printf("🚇  Found a Drupal 8 site")
			dp = strings.Replace(path, d8, "", -1)
			dv = "8"
			return io.EOF
		}
		return nil
	})
	if err != io.EOF {
		log.Fatalf("Could not find a valid Drupal 7 or 8 installation. You will need to manually specify your Drupal root. \n\nHave you run 'composer install' on your project yet?")
	}

	return dp, dv
}

// CoreDrupalFile - Return the core drupal file for the users' installation
func CoreDrupalFile() string {
	rp := GetRootPath()
	dv := GetConfig().Drupal.MajorVersion

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

// GetRootDir - Return the drupal root folder name without workdir
func GetRootDir() string {
	dr := GetRootPath()

	return filepath.Base(dr)
}
