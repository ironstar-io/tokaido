package docker

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker/templates"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/version"
	"gopkg.in/yaml.v2"
)

// GetTokComposePath ...
func GetTokComposePath() string {
	return filepath.Join(conf.GetConfig().Tokaido.Project.Path, "/docker-compose.tok.yml")
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
	composeVolumeYml := []byte(`volumes:
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
	unisonVersion := version.GetUnisonVersion()

	err := yaml.Unmarshal(dockertmpl.ComposeTokDefaults, &tokStruct)
	if err != nil {
		log.Fatalf("Error setting Compose file defaults: %v", err)
	}

	errDrupal := yaml.Unmarshal(getDrupalSettings(), &tokStruct)
	if errDrupal != nil {
		log.Fatalf("Error adding Drupal settings to Compose file: %v", err)
	}

	errUnison := yaml.Unmarshal(dockertmpl.SetUnisonVersion(unisonVersion), &tokStruct)
	if errUnison != nil {
		log.Fatalf("Error setting Unison version: %v", err)
	}

	// Create mysql volume declaration and attachment
	mysqlVolName := "tok_" + conf.GetConfig().Tokaido.Project.Name + "_mysql_database"
	errDbDec := yaml.Unmarshal(dockertmpl.ExternalVolumeDeclare(mysqlVolName), &tokStruct)
	if errDbDec != nil {
		log.Fatalf("Error declaring persistent MySQL volume: %v", err)
	}

	errDbAtt := yaml.Unmarshal(dockertmpl.MysqlVolumeAttach(mysqlVolName), &tokStruct)
	if errDbAtt != nil {
		log.Fatalf("Error attaching persistent MySQL volume: %v", err)
	}

	errCoAtt := yaml.Unmarshal(dockertmpl.ComposerCacheVolumeAttach(), &tokStruct)
	if errCoAtt != nil {
		log.Fatalf("Error attaching persistent Composer cache volume: %v", err)
	}

	if conf.GetConfig().Services.Solr.Enabled {
		var v string
		if conf.GetConfig().Tokaido.Betacontainers {
			v = "edge"
		} else {
			v = "6.6"
		}
		errSolr := yaml.Unmarshal(dockertmpl.EnableSolr(v), &tokStruct)
		if errSolr != nil {
			log.Fatalf("Error enabling Solr in Compose file: %v", err)
		}
	}

	if conf.GetConfig().Services.Redis.Enabled {
		var v string
		if conf.GetConfig().Tokaido.Betacontainers {
			v = "edge"
		} else {
			v = "4.0.11"
		}
		errRedis := yaml.Unmarshal(dockertmpl.EnableRedis(v), &tokStruct)
		if errRedis != nil {
			log.Fatalf("Error enabling Redis in Compose file: %v", err)
		}
	}

	if conf.GetConfig().Services.Mailhog.Enabled {
		v := "v1.0.0"
		errMailhog := yaml.Unmarshal(dockertmpl.EnableMailhog(v), &tokStruct)
		if errMailhog != nil {
			log.Fatalf("Error enabling Mailhog in Compose file: %v", err)
		}
	}

	if conf.GetConfig().Services.Adminer.Enabled {
		v := "4-standalone"
		errAdminer := yaml.Unmarshal(dockertmpl.EnableAdminer(v), &tokStruct)
		if errAdminer != nil {
			log.Fatalf("Error enabling Adminer in Compose file: %v", err)
		}
	}

	if conf.GetConfig().Services.Memcache.Enabled {
		var v string
		if conf.GetConfig().Tokaido.Betacontainers {
			v = "1.5-alpine" // Temporary, there is currently no edge memcache container
		} else {
			v = "1.5-alpine"
		}
		errMemcache := yaml.Unmarshal(dockertmpl.EnableMemcache(v), &tokStruct)
		if errMemcache != nil {
			log.Fatalf("Error enabling Memcache in Compose file: %v", err)
		}
	}

	if conf.GetConfig().Tokaido.Betacontainers {
		errEdge := yaml.Unmarshal(dockertmpl.EdgeContainers(), &tokStruct)
		if errEdge != nil {
			log.Fatalf("Error enabling edge containers in Compose file: %v", err)
		}
	}

	if conf.GetConfig().System.Xdebug.Enabled {
		var v string
		if conf.GetConfig().Tokaido.Betacontainers {
			v = "edge"
		} else {
			v = "latest"
		}
		errSolr := yaml.Unmarshal(dockertmpl.EnableXdebug(v), &tokStruct)
		if errSolr != nil {
			log.Fatalf("Error enabling Xdebug in Compose file: %v", err)
		}
	}

	if runtime.GOOS == "windows" {
		errOs := yaml.Unmarshal(dockertmpl.WindowsAjustments(), &tokStruct)
		if errOs != nil {
			log.Fatalf("Error enabling Windows containers in Compose file: %v", err)
		}
	}

	errCt := yaml.Unmarshal(getCustomTok(), &tokStruct)
	if errCt != nil {
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
	ct := filepath.Join(conf.GetConfig().Tokaido.Project.Path, "/.tok")

	if fs.CheckExists(ct+".yml") == true {
		return filepath.Join(ct, "compose.tok.yml")
	}

	if fs.CheckExists(ct+".yaml") == true {
		return filepath.Join(ct, "compose.tok.yaml")
	}

	return ""
}
