package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/docker/docker/client"
)

var imageMap = []string{
	"tokaido/syslog:1.11.0",
	"tokaido/haproxy:1.11.0",
	"tokaido/varnish:1.11.0",
	"tokaido/nginx:1.11.0",
	"tokaido/php72-fpm:1.11.0",
	"tokaido/admin72-heavy:1.11.0",
	"tokaido/kishu:1.11.0",
	"tokaido/unison:2.51.2",
	"testcafe/testcafe",
	"mysql:5.7",
	"mailhog/mailhog:v1.0.0",
}

// CheckClientVersion - Ensures that the client running on this system is supported
func CheckClientVersion() {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to find Docker installed on your system. Have you run the Docker installer included in this package?")
		log.Fatal(err)
	}

	v, _ := strconv.ParseFloat(cli.ClientVersion(), 32)

	if v < 1.36 {
		panic("The Docker version currently installed on this system does not meet the minimum requirements. It must be v1.36 or higher. This package includes a compatible version for your system.")
	}
}

// LoadTokImages - Import all Tokaido images from this package into the local image cache
func LoadTokImages() {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to find Docker installed on your system. Have you run the Docker installer included in this package?")
		log.Fatal(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	imgPath := filepath.Join(pwd, "images")
	// Find all tar files in the imgPath
	var tarFiles []string
	imgDirectory, err := os.Open(imgPath)
	if err != nil {
		log.Fatal(err)
	}
	defer imgDirectory.Close()

	files, err := imgDirectory.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".tar" {
			tarFiles = append(tarFiles, filepath.Join(imgPath, file.Name()))
		}
	}

	//Load images
	var wg sync.WaitGroup
	wg.Add(len(tarFiles))
	for _, tarFile := range tarFiles {
		go func(tar string) {
			defer wg.Done()
			file, err := os.Open(tar)
			if err != nil {
				log.Fatal(err)
			}

			_, err = cli.ImageLoad(context.Background(), io.Reader(file), false)
			if err != nil {
				log.Fatal(err)
			}
		}(tarFile)
	}

	wg.Wait()
}
