package docker

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker/templates"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"fmt"
	"log"
	"os"
)

var tokComposePath = fs.WorkDir() + "/docker-compose.tok.yml"

// HardCheckTokCompose ...
func HardCheckTokCompose() {
	var _, err = os.Stat(tokComposePath)

	// create file if not exists
	if os.IsNotExist(err) {
		fmt.Println(`
🤷‍  No docker-compose.tok.yml file found. Have you run 'tok init'?
		`)
		log.Fatal("Exiting without change")
	}
}

// FindOrCreateTokCompose ...
func FindOrCreateTokCompose() {
	var _, err = os.Stat(tokComposePath)

	// create file if not exists
	if os.IsNotExist(err) {
		fmt.Println(`🏯  Generating a new docker-compose.tok.yml file`)

		fs.TouchByteArray(tokComposePath, dockertmpl.TokComposeYml)
	}
}
