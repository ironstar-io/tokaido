package proxy

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/ssl"

	"fmt"
	"path/filepath"
)

const proxy = "proxy"

// Setup ...
func Setup() {
	buildDirectories()

	ssl.Configure(getProxyClientTLSDir())

	GenerateProxyDockerCompose()
	DockerComposeUp()

	ConfigureUnison()
	ConfigureProjectHostsfile()

	ConfigureYamanote()

	ConfigureProjectNginx()
	RestartContainer(proxy)
}

// ConfigureProjectHostsfile ...
func ConfigureProjectHostsfile() {
	pn := conf.GetConfig().Tokaido.Project.Name
	ConfigureHostsfile(pn + "." + constants.ProxyDomain)
}

// ConfigureProjectNginx ...
func ConfigureProjectNginx() {
	h, err := docker.GetContainerIP("haproxy")
	if err != nil {
		fmt.Printf("%s. Skipping HTTPS proxy setup...\n", err)
		return
	}

	if h == "" {
		fmt.Println("The haproxy container doesn't appear to be running. Skipping HTTPS proxy setup...")
		return
	}

	pp := constants.HTTPSProtocol + h + ":" + constants.HaproxyInternalPort

	pn := conf.GetConfig().Tokaido.Project.Name
	do := pn + `.` + constants.ProxyDomain

	nc := GenerateNginxConf(do, pp)

	np := filepath.Join(getProxyClientConfdDir(), pn+".conf")
	fs.Replace(np, nc)
}
