package docker

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/docker/templates"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

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
		fmt.Println(`
🤷‍  No docker-compose.tok.yml file found. Have you run 'tok init'?
		`)
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

	err := yaml.Unmarshal(dockertmpl.ComposeTokDefaults, &tokStruct)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	appendCustomTok(tokStruct)
	appendDrupalRoot(tokStruct)

	return tokStruct
}

func appendCustomTok(tokStruct dockertmpl.ComposeDotTok) {
	ctp := customTokPath()

	if fs.CheckExists(ctp) == false {
		return
	}

	if ctp != "" {
		customTok, err := ioutil.ReadFile(ctp)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		errTwo := yaml.Unmarshal(customTok, &tokStruct)
		if errTwo != nil {
			log.Fatalf("error: %v", errTwo)
		}
	}
}

func appendDrupalRoot(tokStruct dockertmpl.ComposeDotTok) {
	dr := conf.GetRootDir()
	dry := drupalRootTmpl(dr)

	err := yaml.Unmarshal(dry, &tokStruct)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func drupalRootTmpl(drupalRoot string) []byte {
	return []byte(`services:
  fpm:
    environment:
      DRUPAL_ROOT: ` + drupalRoot + `
  nginx:
    environment:
      DRUPAL_ROOT: ` + drupalRoot + `
  drush:
    environment:
      DRUPAL_ROOT: ` + drupalRoot)
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
