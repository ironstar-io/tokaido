package ssh

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"
	"bitbucket.org/ironstar/tokaido-cli/system/ssh/templates"

	"bufio"
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type tokConf struct {
	ProjectName string
	DrushPort   string
}

var tokConfigFile = filepath.Join(fs.HomeDir(), "/.ssh/tok_config")

// ProcessTokConfig ...
func ProcessTokConfig() {
	createOrUpdate()
	processConfig()
}

// createOrUpdate ...
func createOrUpdate() {
	// detect if file exists
	var _, err = os.Stat(tokConfigFile)

	// create file if not exists
	if os.IsNotExist(err) {
		tpl := generateTokConfig()
		createTokConfig(tpl, tokConfigFile)
	} else {
		updateTokConfig()
	}
}

// generateTokConfig - Generate a `tok_config` file
func generateTokConfig() string {
	config := conf.GetConfig()
	s := tokConf{DrushPort: docker.LocalPort("drush", "22"), ProjectName: config.Project}

	tmpl := template.New("tok_config")
	tmpl, err := tmpl.Parse(sshtmpl.TokConfTmplStr)

	if err != nil {
		log.Fatal("Parse: ", err)
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, s); err != nil {
		log.Fatal("Parse: ", err)
	}

	return tpl.String()
}

// createTokConfig - Write generated config file to `~/.ssh/`
func createTokConfig(body string, path string) {
	var file, err = os.Create(path)
	if err != nil {
		log.Fatal("Create: ", err)
	}

	_, _ = file.WriteString(body)

	defer file.Close()
}

// updateTokConfig - Update a `tok_config` file in `~/.ssh/`
func updateTokConfig() {
	f, err := os.Open(tokConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	config := conf.GetConfig()

	hostLine := "Host " + config.Project + ".tok"
	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	var buffer bytes.Buffer
	var block bytes.Buffer
	replaceBlock := false
	blockWritten := false
	for scanner.Scan() {
		l := scanner.Text()

		if strings.Contains(scanner.Text(), hostLine) {
			replaceBlock = true
			blockWritten = true
		}

		if len(strings.TrimSpace(l)) != 0 {
			block.Write([]byte(l + "\n"))
			continue
		}

		if replaceBlock == true {
			buffer.Write([]byte(generateTokConfig()))
		} else {
			buffer.Write([]byte(block.String() + "\n"))
		}

		block.Reset()
		replaceBlock = false
	}

	if blockWritten == false {
		buffer.Write([]byte(generateTokConfig()))
	}

	createTokConfig(buffer.String(), tokConfigFile+"-tmp")
	replaceTokConfig()
}

// replaceTokConfig - Replace `~/.ssh/tok_config` with `~/.ssh/tok_config-tmp` file
func replaceTokConfig() {
	tmpTokConfig := tokConfigFile + "-tmp"

	// Remove the original tok_config file
	os.Remove(tokConfigFile)

	// Rename `tok_config-tmp` to be the new `tok_config` file
	os.Rename(tmpTokConfig, tokConfigFile)
}
