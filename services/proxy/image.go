package proxy

import (
	"log"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

// pullImages - Pull all images in compose file
func pullImages() {
	// Don't pull if the proxy config hasn't been created yet
	h, err := homedir.Dir()
	if err != nil {
		log.Fatalln("Unable to resolve home directory so can't initialise Tokaido. Sorry!")
	}
	gc := filepath.Join(h, ".tok", "proxy", "docker-compose.yml")
	_, err = os.Stat(gc)
	if err != nil {
		return
	}

	composeStdout("pull")
}
