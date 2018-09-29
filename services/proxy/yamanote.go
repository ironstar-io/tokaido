package proxy

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"

	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ConfigureYamanote ...
func ConfigureYamanote() {
	ConfigureYamanoteNginx()

	CreateOrAppendGatsbyEnvFile()

	RestartContainer("yamanote")
}

// CreateOrAppendGatsbyEnvFile ...
func CreateOrAppendGatsbyEnvFile() {
	if fs.CheckExists(getProxyClientGatsbyEnv()) == true {
		AppendToGatsbyEnvFile()
		return
	}

	CreateGatsbyEnvFile()
	return
}

// AppendToGatsbyEnvFile ...
func AppendToGatsbyEnvFile() {
	f, err := os.Open(getProxyClientGatsbyEnv())
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()
	pn := conf.GetConfig().Tokaido.Project.Name

	replaceNeeded := false
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), constants.GatsbyTokaidoProjectsKey) && !strings.Contains(scanner.Text(), pn) {
			replaceNeeded = true
			ps := []string{pn}
			gk := constants.GatsbyTokaidoProjectsKey + "="
			gs := strings.Split(scanner.Text(), "=")[1]
			for _, v := range strings.Split(gs, ",") {
				ps = append(ps, v)
			}
			buffer.Write([]byte(gk + strings.Join(ps, ",") + "\n"))
		} else {
			buffer.Write([]byte(scanner.Text() + "\n"))
		}
	}

	if replaceNeeded {
		fs.Replace(getProxyClientGatsbyEnv(), buffer.Bytes())
	}
}

// DetachFromGatsbyEnvFile ...
func DetachFromGatsbyEnvFile() {
	if fs.CheckExists(getProxyClientGatsbyEnv()) == false {
		return
	}

	f, err := os.Open(getProxyClientGatsbyEnv())
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()
	pn := conf.GetConfig().Tokaido.Project.Name

	replaceNeeded := false
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), constants.GatsbyTokaidoProjectsKey) && strings.Contains(scanner.Text(), pn) {
			replaceNeeded = true
			var ps []string
			gk := constants.GatsbyTokaidoProjectsKey + "="
			gs := strings.Split(scanner.Text(), "=")[1]
			for _, v := range strings.Split(gs, ",") {
				if strings.Contains(v, pn) {
					continue
				}
				ps = append(ps, v)
			}
			buffer.Write([]byte(gk + strings.Join(ps, ",") + "\n"))
		} else {
			buffer.Write([]byte(scanner.Text() + "\n"))
		}
	}

	if replaceNeeded {
		fs.Replace(getProxyClientGatsbyEnv(), buffer.Bytes())
	}
}

// CreateGatsbyEnvFile ...
func CreateGatsbyEnvFile() {
	pn := conf.GetConfig().Tokaido.Project.Name
	b := []byte(constants.GatsbyTokaidoProjectsKey + "=" + pn)

	fs.TouchByteArray(getProxyClientGatsbyEnv(), b)
}

// ConfigureYamanoteNginx ...
func ConfigureYamanoteNginx() {
	pp := constants.HTTPProtocol + constants.YamanoteInternalDomain + ":" + constants.YamanoteInternalPort

	nc := GenerateNginxConf(constants.ProxyDomain, pp)

	np := filepath.Join(getProxyClientConfdDir(), constants.ProxyDomain+".conf")
	fs.Replace(np, nc)
}
