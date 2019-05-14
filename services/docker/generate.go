package docker

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
	dockertmpl "github.com/ironstar-io/tokaido/services/docker/templates"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
	. "github.com/logrusorgru/aurora"
	yaml "gopkg.in/yaml.v2"
)

// GetTokComposePath ...
func GetTokComposePath() string {
	return filepath.Join(conf.GetProjectPath(), "/docker-compose.tok.yml")
}

// HardCheckTokCompose ...
func HardCheckTokCompose() {
	var _, err = os.Stat(GetTokComposePath())

	// create file if not exists
	if os.IsNotExist(err) {
		console.Println(`🤷‍  No docker-compose.tok.yml file found. Have you run 'tok up'?`, "×")
		log.Fatal("Exiting without change")
	}
}

// FindOrCreateTokCompose ...
func FindOrCreateTokCompose() {
	var _, errf = os.Stat(GetTokComposePath())

	// create file if not exists
	if os.IsNotExist(errf) {
		console.Println(`🏯  Generating a new docker-compose.tok.yml file`, "")

		CreateOrReplaceTokCompose(MarshalledDefaults())
		return
	}

	if conf.GetConfig().Tokaido.Customcompose == true {
		fmt.Println(Magenta("    You've chosen to use a custom Compose file. Tokaido will not update your docker-compose.tok.yml file"))
		StripModWarning()
		return
	}

	CreateOrReplaceTokCompose(MarshalledDefaults())
}

// CreateOrReplaceTokCompose ...
func CreateOrReplaceTokCompose(tokComposeYml []byte) {
	tc := GetTokComposePath()
	if conf.GetConfig().Tokaido.Customcompose == false {
		fs.TouchOrReplace(tc, append(dockertmpl.ModWarning[:], tokComposeYml[:]...))
		return
	}

	fs.TouchOrReplace(tc, tokComposeYml)
}

// MarshalledDefaults ...
func MarshalledDefaults() []byte {
	emptyTokStruct := UnmarshalledDefaults()

	tokComposeYml, err := yaml.Marshal(&emptyTokStruct)
	if err != nil {
		log.Fatalf("error: %v", err)
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

	return tokComposeYml
}

// UnmarshalledDefaults ...
func UnmarshalledDefaults() conf.ComposeDotTok {
	tokStruct := conf.ComposeDotTok{}
	xdebugImageVersion := conf.GetConfig().Tokaido.Stability
	phpVersion := conf.GetConfig().Tokaido.Phpversion

	if conf.GetConfig().Global.Syncservice == "docker" {
		utils.DebugString("attaching repo to /tokaido/site folder using direct Docker mount")

		err := yaml.Unmarshal(dockertmpl.ComposeTokDefaultsDockerVolume, &tokStruct)
		if err != nil {
			log.Fatalf("Error setting Compose file defaults: %v", err)
		}

		err = yaml.Unmarshal(dockertmpl.TokaidoDockerSiteVolumeAttach(conf.GetProjectPath()), &tokStruct)
		if err != nil {
			log.Fatalf("Error attaching persistent /tokaido/site volume: %v", err)
		}

	} else {
		utils.DebugString("attaching repo to /tokaido/site folder using Fusion sync")
		err := yaml.Unmarshal(dockertmpl.ComposeTokDefaultsFusionSync, &tokStruct)
		if err != nil {
			log.Fatalf("Error setting Compose file defaults: %v", err)
		}

		// Configure Fusion Sync volumes
		siteVolName := "tok_" + conf.GetConfig().Tokaido.Project.Name + "_tokaido_site"
		err = yaml.Unmarshal(dockertmpl.ExternalVolumeDeclare(siteVolName), &tokStruct)
		if err != nil {
			log.Fatalf("Error declaring persistent /tokaido/site volume: %v", err)
		}

		err = yaml.Unmarshal(dockertmpl.TokaidoFusionSiteVolumeAttach(conf.GetProjectPath(), siteVolName), &tokStruct)
		if err != nil {
			log.Fatalf("Error attaching persistent /tokaido/site volume: %v", err)
		}
	}

	err := yaml.Unmarshal(getDrupalSettings(), &tokStruct)
	if err != nil {
		log.Fatalf("Error adding Drupal settings to Compose file: %v", err)
	}

	// Create mysql volume declaration and attachment
	mysqlVolName := "tok_" + conf.GetConfig().Tokaido.Project.Name + "_mysql_database"
	err = yaml.Unmarshal(dockertmpl.ExternalVolumeDeclare(mysqlVolName), &tokStruct)
	if err != nil {
		log.Fatalf("Error declaring persistent MySQL volume: %v", err)
	}

	err = yaml.Unmarshal(dockertmpl.MysqlVolumeAttach(mysqlVolName), &tokStruct)
	if err != nil {
		log.Fatalf("Error attaching persistent MySQL volume: %v", err)
	}

	if conf.GetConfig().Services.Solr.Enabled {
		v := "6.6" // Solr only supports 6.6 right now, no need to decide between stable/edge/experimental versions

		err = yaml.Unmarshal(dockertmpl.EnableSolr(v), &tokStruct)
		if err != nil {
			log.Fatalf("Error enabling Solr in Compose file: %v", err)
		}
	}

	if conf.GetConfig().Services.Redis.Enabled {
		v := "4.0.12"
		err = yaml.Unmarshal(dockertmpl.EnableRedis(v), &tokStruct)
		if err != nil {
			log.Fatalf("Error enabling Redis in Compose file: %v", err)
		}
	}

	if conf.GetConfig().Services.Mailhog.Enabled {
		v := "v1.0.0"
		err = yaml.Unmarshal(dockertmpl.EnableMailhog(v), &tokStruct)
		if err != nil {
			log.Fatalf("Error enabling Mailhog in Compose file: %v", err)
		}
	}

	if conf.GetConfig().Services.Adminer.Enabled {
		v := "4-standalone"
		err = yaml.Unmarshal(dockertmpl.EnableAdminer(v), &tokStruct)
		if err != nil {
			log.Fatalf("Error enabling Adminer in Compose file: %v", err)
		}
	}

	if conf.GetConfig().Services.Memcache.Enabled {
		v := "1.5-alpine"

		err = yaml.Unmarshal(dockertmpl.EnableMemcache(v), &tokStruct)
		if err != nil {
			log.Fatalf("Error enabling Memcache in Compose file: %v", err)
		}
	}

	// Set our stability version
	err = yaml.Unmarshal(dockertmpl.StabilityLevel(phpVersion, conf.GetConfig().Tokaido.Stability), &tokStruct)
	if err != nil {
		log.Fatalf("Error updating stability version for containers in Compose file: %v", err)
	}

	if conf.GetConfig().Tokaido.Xdebug {
		err = yaml.Unmarshal(dockertmpl.EnableXdebug(phpVersion, xdebugImageVersion), &tokStruct)
		if err != nil {
			log.Fatalf("Error enabling Xdebug in Compose file: %v", err)
		}
	}

	err = yaml.Unmarshal(getCustomTok(), &tokStruct)
	if err != nil {
		log.Fatalf("Error enabling custom Compose config: %v", err)
	}

	return tokStruct
}

func getCustomTok() []byte {
	dc := &conf.ComposeDotTok{
		Version:  "2",
		Services: conf.GetConfig().Services,
	}

	// Nulify the invalid docker-compose file values
	dc.Services.Memcache.Enabled = false
	dc.Services.Solr.Enabled = false
	dc.Services.Redis.Enabled = false
	dc.Services.Mailhog.Enabled = false
	dc.Services.Adminer.Enabled = false

	cc, err := yaml.Marshal(dc)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return cc
}

func getDrupalSettings() []byte {
	return dockertmpl.DrupalSettings(conf.GetRootDir(), conf.GetConfig().Tokaido.Project.Name)
}

// StripModWarning ...
func StripModWarning() {
	tc := GetTokComposePath()
	f, openErr := os.Open(tc)
	if openErr != nil {
		fmt.Println(openErr)
		return
	}

	defer f.Close()

	warningPresent := false
	warningOne := "# WARNING: THIS FILE IS MANAGED DIRECTLY BY TOKAIDO."
	warningTwo := "# DO NOT MAKE MODIFICATIONS HERE, THEY WILL BE OVERWRITTEN"

	var buffer bytes.Buffer
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), warningOne) || strings.Contains(scanner.Text(), warningTwo) {
			warningPresent = true
			continue
		} else {
			buffer.Write([]byte(scanner.Text() + "\n"))
		}
	}

	if warningPresent == true {
		fs.Replace(tc, buffer.Bytes())
	}
}

func customTokPath() string {
	ct := filepath.Join(conf.GetProjectPath(), "/.tok")

	if fs.CheckExists(ct+".yml") == true {
		return filepath.Join(ct, "compose.tok.yml")
	}

	if fs.CheckExists(ct+".yaml") == true {
		return filepath.Join(ct, "compose.tok.yaml")
	}

	return ""
}
