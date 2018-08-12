package proxy

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/system/hostsfile"
	"github.com/ironstar-io/tokaido/system/ssl"

	"fmt"
)

const proxy = "proxy"

// Setup ...
func Setup() {
	buildDirectories()

	ssl.Configure(getProxyClientTLSDir())

	ConfigureUnison()
	ConfigureHostsfile()
	ConfigureNginx()
}

// ConfigureUnison ...
func ConfigureUnison() {
	GenerateProxyDockerCompose()
	DockerComposeUp()

	unison.CreateOrUpdatePrf(UnisonPort(), proxy, getProxyClientDir())
	unison.CreateSyncService(proxy, getProxyClientDir())
}

// ConfigureHostsfile ...
func ConfigureHostsfile() {
	pn := conf.GetConfig().Tokaido.Project.Name
	err := hostsfile.AddEntry(pn + "." + constants.ProxyDomain)
	if err != nil {
		fmt.Println("There was an issue updating your hostsfile. Your hostsfile can be amended manually in order to enable this feature. See XXXXXXX for more information.")
		fmt.Println(err)
	}
}

// ConfigureNginx ...
func ConfigureNginx() {
	hp := docker.LocalPort("haproxy", "8443")
	if hp == "" {
		fmt.Println("The haproxy container doesn't appear to be running. Skipping HTTPS proxy setup...")
		return
	}

	RebuildNginxConfigFile(hp)
	RestartContainer(proxy)
}
