package ssh

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
)

var sshConfigFile = fs.HomeDir() + "/.ssh/config"

// ProcessConfig ...
func processConfig() {
	tokInclude := checkTokInclude()
	if tokInclude == false {
		prependTokInclude()
	}
}

// checkTokInclude ...
func checkTokInclude() bool {
	f, err := os.Open(sshConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	includeLine := "Include ~/.ssh/tok_config"
	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), includeLine) {
			return true
		}
	}

	return false
}

// prependTokInclude - Prepend `~/.ssh/config` with `Include ~/.ssh/tok_config`
func prependTokInclude() {
	f, err := os.Open(sshConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	var buffer bytes.Buffer
	for scanner.Scan() {
		buffer.Write([]byte(scanner.Text() + "\n"))
	}
	buffer.Write([]byte("\nInclude ~/.ssh/tok_config\n"))

	createTempConfig(buffer.String())
	replaceConfig()
}

// createTempConfig - Write generated temp config file to `~/.ssh/`
func createTempConfig(body string) {
	var file, err = os.Create(sshConfigFile + "-tmp")
	if err != nil {
		log.Fatal("Create: ", err)
	}

	_, _ = file.WriteString(body)

	defer file.Close()
}

// replaceConfig - Replace `~/.ssh/config` with `~/.ssh/config-tmp` file
func replaceConfig() {
	tmpConfig := sshConfigFile + "-tmp"

	// Remove the original `config` file
	os.Remove(sshConfigFile)

	// Rename `config-tmp` to be the new `config` file
	os.Rename(tmpConfig, sshConfigFile)
}
