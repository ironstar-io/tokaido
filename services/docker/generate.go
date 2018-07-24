package docker

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/docker/templates"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"
	"bitbucket.org/ironstar/tokaido-cli/system/version"

	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var tokComposePath = fs.WorkDir() + "/docker-compose.tok.yml"

// HardCheckTokCompose ...
func HardCheckTokCompose() {
	var _, err = os.Stat(tokComposePath)

	// create file if not exists
	if os.IsNotExist(err) {
		fmt.Println(`🤷‍  No docker-compose.tok.yml file found. Have you run 'tok up'?`)
		log.Fatal("Exiting without change")
	}
}

// FindOrCreateTokCompose ...
func FindOrCreateTokCompose() {
	var _, errf = os.Stat(tokComposePath)

	// create file if not exists
	if os.IsNotExist(errf) {
		fmt.Println(`🏯  Generating a new docker-compose.tok.yml file`)

		CreateOrReplaceTokCompose(MarshalledDefaults())
		return
	}

	if conf.GetConfig().CustomCompose == true {
		StripModWarning()
		return
	}

	CreateOrReplaceTokCompose(MarshalledDefaults())
}

// CreateOrReplaceTokCompose ...
func CreateOrReplaceTokCompose(tokComposeYml []byte) {
	if conf.GetConfig().CustomCompose == false {
		fs.TouchOrReplace(tokComposePath, append(dockertmpl.ModWarning[:], tokComposeYml[:]...))
		return
	}

	fs.TouchOrReplace(tokComposePath, tokComposeYml)
}

// MarshalledDefaults ...
func MarshalledDefaults() []byte {
	emptyTokStruct := UnmarshalledDefaults()

	tokComposeYml, err := yaml.Marshal(&emptyTokStruct)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return tokComposeYml
}

// UnmarshalledDefaults ...
func UnmarshalledDefaults() dockertmpl.ComposeDotTok {
	tokStruct := dockertmpl.ComposeDotTok{}
	unisonVersion := version.GetUnisonVersion()

	err := yaml.Unmarshal(dockertmpl.ComposeTokDefaults, &tokStruct)
	if err != nil {
		log.Fatalf("Error setting Compose file defaults: %v", err)
	}

	err2 := yaml.Unmarshal(getDrupalSettings(), &tokStruct)
	if err2 != nil {
		log.Fatalf("Error adding Drupal settings to Compose file: %v", err)
	}

	errUnison := yaml.Unmarshal(dockertmpl.SetUnisonVersion(unisonVersion), &tokStruct)
	if errUnison != nil {
		log.Fatalf("Error setting Unison version: %v", err)
	}

	if conf.GetConfig().Solr.Enable {
		var v string
		if conf.GetConfig().BetaContainers {
			v = "edge"
		} else {
			v = conf.GetConfig().Solr.Version
		}
		errSolr := yaml.Unmarshal(dockertmpl.EnableSolr(v), &tokStruct)
		if errSolr != nil {
			log.Fatalf("Error enabling Solr in Compose file: %v", err)
		}
	}

	if conf.GetConfig().Memcache.Enable {
		var v string
		if conf.GetConfig().BetaContainers {
			v = "edge"
		} else {
			v = conf.GetConfig().Memcache.Version
		}
		errMemcache := yaml.Unmarshal(dockertmpl.EnableMemcache(v), &tokStruct)
		if errMemcache != nil {
			log.Fatalf("Error enabling Memcache in Compose file: %v", err)
		}
	}

	if conf.GetConfig().BetaContainers {
		err3 := yaml.Unmarshal(dockertmpl.EdgeContainers(), &tokStruct)
		if err3 != nil {
			log.Fatalf("Error enabling edge containers in Compose file: %v", err)
		}
	}

	err4 := yaml.Unmarshal(getCustomTok(), &tokStruct)
	if err4 != nil {
		log.Fatalf("Error enabling custom Compose config: %v", err)
	}

	return tokStruct
}

func getCustomTok() []byte {
	ctp := customTokPath()

	if fs.CheckExists(ctp) == false {
		return nil
	}

	if ctp != "" {
		ct, err := ioutil.ReadFile(ctp)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		return ct
	}

	return nil
}

func getDrupalSettings() []byte {
	return dockertmpl.DrupalSettings(conf.GetRootDir(), conf.GetConfig().Project)
}

// StripModWarning ...
func StripModWarning() {
	f, openErr := os.Open(tokComposePath)
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
		fs.Replace(tokComposePath, buffer.Bytes())
	}
}

func customTokPath() string {
	ct := fs.WorkDir() + "/.tok/compose.tok"

	if fs.CheckExists(ct+".yml") == true {
		return ct + ".yml"
	}

	if fs.CheckExists(ct+".yaml") == true {
		return ct + ".yaml"
	}

	return ""
}
