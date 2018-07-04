package drupal

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getRootPath() string {
	wd := fs.WorkDir()
	c := conf.GetConfig().Drupal.Path
	if c != "" {
		return wd + c
	}

	var dp string
	dc := "/core/lib/Drupal.php"
	err := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.Contains(path, dc) == true {
			dp = strings.Replace(path, dc, "", -1)
			return io.EOF
		}
		return nil
	})
	if err != io.EOF {
		log.Fatalf("There was an error when searching for your Drupal installation [%v]\n", err)
	}

	if dp == "" {
		log.Fatal("Tokaido was unable to find your Drupal installation. Exiting...")
	}

	return dp
}

// CoreDrupalFile - Return the core drupal file for the users' installation
func CoreDrupalFile() string {
	return getRootPath() + "/core/lib/Drupal.php"
}
