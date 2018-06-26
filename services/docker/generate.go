package docker

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker/templates"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
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

		tokStruct := dockertmpl.ComposeDotTok{}

		err := yaml.Unmarshal(dockertmpl.ComposeTokDefaults, &tokStruct)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		tokComposeYml, err := yaml.Marshal(&tokStruct)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		fs.TouchByteArray(tokComposePath, tokComposeYml)
	}
}
