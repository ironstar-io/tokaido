package xdebug

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

var tokComposePath = fs.WorkDir() + "/docker-compose.tok.yml"

// Configure - Amend docker-compose.tok.yml to be compatible with xdebug
func Configure() {
	networkUp := docker.CheckNetworkUp()

	regenerateGateway(networkUp)
}

// regenerateGateway - Regenerate the docker-compose.tok.yml file and change the IP value in
func regenerateGateway(networkUp bool) {
	docker.HardCheckTokCompose()

	if networkUp == false {
		docker.Up()
	}

	networkGateway := docker.GetGateway()

	f, openErr := os.Open(tokComposePath)
	if openErr != nil {
		fmt.Println(openErr)
		return
	}

	defer f.Close()

	closePHP := "XDEBUG_REMOTE_HOST: "
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), closePHP) {
			buffer.Write([]byte("      XDEBUG_REMOTE_HOST: " + networkGateway + "\n"))
		} else {
			buffer.Write([]byte(scanner.Text() + "\n"))
		}
	}

	fs.Replace(tokComposePath, buffer.Bytes())

	docker.KillContainer("fpm")
	docker.UpContainer("fpm")
}
