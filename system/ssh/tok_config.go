package ssh

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	// "bufio"
	"bytes"
	"log"
	"os"
	// "strings"
	"text/template"
)

type tokConf struct {
	ProjectName string
	DrushPort   string
}

var tokConfigFile = fs.HomeDir() + "/.ssh/tok_config"
var tokConfigTmpStr = `Host {{.ProjectName}}.tok
    HostName localhost
    Port {{.DrushPort}}
    User tok
    IdentityFile ~/.ssh/tok_ssh.key
    ForwardAgent yes
    StrictHostKeyChecking no
	UserKnownHostsFile /dev/null

`

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
		generateTokConfig()
	} else {
		updateTokConfig()
	}
}

// generateTokConfig - Generate a `tok_config` file
func generateTokConfig() {
	config := conf.GetConfig()
	s := tokConf{DrushPort: docker.LocalPort("drush", "22"), ProjectName: config.Project}

	tmpl := template.New("tok_config")
	tmpl, err := tmpl.Parse(tokConfigTmpStr)

	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, s); err != nil {
		log.Fatal("Parse: ", err)
		return
	}

	createTokConfig(tpl.String(), tokConfigFile)
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
	// TODO

	// f, err := os.Open(sshConfigFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()

	// portLine := "root = socket://localhost:"
	// // Splits on newlines by default.
	// scanner := bufio.NewScanner(f)

	// var prfString []string
	// for scanner.Scan() {
	// 	if strings.Contains(scanner.Text(), portLine) {
	// 		prfString = append(prfString, portLine+LocalPort())
	// 	} else {
	// 		prfString = append(prfString, scanner.Text())
	// 	}
	// }

	// createTokConfig(strings.Join(prfString, "\n"), tokConfigFile+"-tmp")
	// replaceTokConfig()
}

// replaceTokConfig - Replace `~/.ssh/tok_config` with `~/.ssh/tok_config-tmp` file
func replaceTokConfig() {
	tmpTokConfig := tokConfigFile + "-tmp"

	// Remove the original tok_config file
	os.Remove(tokConfigFile)

	// Rename `tok_config-tmp` to be the new `tok_config` file
	os.Rename(tmpTokConfig, tokConfigFile)
}
