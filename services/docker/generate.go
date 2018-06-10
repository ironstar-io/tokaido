package docker

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker/templates"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"fmt"
	"log"
	"os"
)

var tokComposePath = fs.WorkDir() + "/docker-compose.tok.yml"

// CheckForTokComposeFile ...
func CheckForTokComposeFile() {
	var _, err = os.Stat(tokComposePath)

	// create file if not exists
	if os.IsNotExist(err) {
		fmt.Println(`🏯  Generating a new docker-compose.tok.yml file`)

		createTokCompose()
	}
}

// createTokCompose
func createTokCompose() {
	var file, err = os.Create(tokComposePath)
	if err != nil {
		log.Fatal("Create: ", err)
	}

	_, _ = file.WriteString(dockertmpl.TokComposeYml)

	defer file.Close()
}
