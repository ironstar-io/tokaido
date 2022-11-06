package ssh

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/fs"
	sshtmpl "github.com/ironstar-io/tokaido/system/ssh/templates"

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

var tokConfigFile = filepath.Join(fs.HomeDir(), ".ssh", "tok_config")

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
		fs.TouchByteArray(tokConfigFile, generateTokConfig())
	} else {
		updateTokConfig()
	}
}

// generateTokConfig - Generate a `tok_config` file
func generateTokConfig() []byte {
	s := tokConf{DrushPort: docker.LocalPort("ssh", "22"), ProjectName: conf.GetConfig().Tokaido.Project.Name}

	tmpl := template.New("tok_config")
	tmpl, err := tmpl.Parse(sshtmpl.TokConfTmplStr)

	if err != nil {
		log.Fatal("Parse: ", err)
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, s); err != nil {
		log.Fatal("Parse: ", err)
	}

	return tpl.Bytes()
}

// updateTokConfig - Update a `tok_config` file in `~/.ssh/`
func updateTokConfig() {
	f, err := os.Open(tokConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hostLine := "Host " + conf.GetConfig().Tokaido.Project.Name + ".tok"
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

	fs.Replace(tokConfigFile, buffer.Bytes())
}
