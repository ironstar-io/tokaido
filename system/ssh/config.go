package ssh

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ProcessConfig ...
func processConfig() {
	tokInclude := checkTokInclude()
	if tokInclude == false {
		prependTokInclude()
	}
}

func sshConfigFile() string {
	return filepath.Join(fs.HomeDir(), "/.ssh/config")
}

func checkOrCreateSSHConfig() {
	sc := sshConfigFile()
	if fs.CheckExists(sc) == false {
		err := fs.TouchEmpty(sc)
		if err != nil {
			fmt.Println("There was an error creating the config directory:", err)
		}
	}
}

func checkTokInclude() bool {
	checkOrCreateSSHConfig()

	f, err := os.Open(sshConfigFile())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	includeLine := "Include ~/.ssh/tok_config"
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
	checkOrCreateSSHConfig()

	f, err := os.Open(sshConfigFile())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var buffer bytes.Buffer
	scanner := bufio.NewScanner(f)
	buffer.Write([]byte("Include ~/.ssh/tok_config\n"))
	for scanner.Scan() {
		buffer.Write([]byte(scanner.Text() + "\n"))
	}

	fs.Replace(sshConfigFile(), buffer.Bytes())
}
