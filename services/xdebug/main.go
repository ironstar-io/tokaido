package xdebug

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"

	"log"

	"github.com/logrusorgru/aurora"
	yaml "gopkg.in/yaml.v2"
)

// Configure - Amend docker-compose.tok.yml to be compatible with xdebug
func Configure() {
	gprj, err := conf.GetGlobalProjectSettings()
	if err != nil {
		panic(err)
	}

	if gprj.Xdebug.Enabled {
		regenerateGateway(gprj)
	}

}

// regenerateGateway - Regenerate the docker-compose.tok.yml file and change the IP value in
func regenerateGateway(gprj *conf.Project) {
	docker.HardCheckTokCompose()
	fpmPort := "9000"
	if gprj.Xdebug.FpmPort != 0 {
		fpmPort = strconv.Itoa(gprj.Xdebug.FpmPort)
	}

	fmt.Printf("🧐   Enabling PHP FPM XDebug Support on port %s\n", aurora.Cyan(fpmPort))

	var networkGateway string
	if runtime.GOOS == "darwin" {
		// MacOS uses a different address for the remote host
		networkGateway = "host.docker.internal"
	} else {
		// Use the docker gateway IP for the host address
		pn := conf.GetConfig().Tokaido.Project.Name
		networkGateway = docker.GetGateway(pn)
	}

	tokStruct := docker.UnmarshalledDefaults()

	fpmYaml := xdebugFpmEnv(networkGateway, fpmPort)
	err := yaml.Unmarshal(fpmYaml, &tokStruct)
	if err != nil {
		log.Fatalf("Error unmarshalling fpm yaml: %v", err)
	}

	tokComposeYml, err := yaml.Marshal(&tokStruct)
	if err != nil {
		log.Fatalf("Error marhsalling fpm tok yaml: %v", err)
	}

	// Append the volume setting on to the docker-compose setting directly
	mysqlVolName := "tok_" + conf.GetConfig().Tokaido.Project.Name + "_mysql_database"
	siteVolName := "tok_" + conf.GetConfig().Tokaido.Project.Name + "_tokaido_site"
	composeVolumeYml := []byte(`volumes:
  ` + siteVolName + `:
    external: true
  ` + mysqlVolName + `:
    external: true
  tok_composer_cache:
    external: true`)

	tokComposeYml = append(tokComposeYml[:], composeVolumeYml[:]...)

	docker.CreateOrReplaceTokCompose(tokComposeYml)

	docker.KillContainer("fpm")
	docker.UpContainer("fpm")
}

func xdebugFpmEnv(networkGateway string, fpmPort string) []byte {
	return []byte(`
services:
  fpm:
    environment:
      XDEBUG_REMOTE_ENABLE: "ON"
      XDEBUG_REMOTE_HOST: "` + networkGateway + `"
      XDEBUG_REMOTE_PORT: "` + fpmPort + `"`)
}
