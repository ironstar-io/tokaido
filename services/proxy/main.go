package proxy

import (
	"fmt"
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/system/hostsfile"
	"github.com/ironstar-io/tokaido/system/ssl"
)

const proxy = "proxy"

// Setup ...
func Setup() {
	pn := conf.GetConfig().Tokaido.Project.Name
	buildDirectories()

	hp := docker.LocalPort("haproxy", "8443")
	if hp == "" {
		fmt.Println("The haproxy container doesn't appear to be running. Skipping HTTPS proxy setup...")
		return
	}

	ssl.Configure(getProxyClientTLSDir())

	GenerateProxyDockerCompose()
	DockerComposeUp()

	unison.CreateOrUpdatePrf(UnisonPort(), proxy, getProxyClientDir())
	unison.CreateSyncService(proxy, getProxyClientDir())

	err := hostsfile.AddEntry(pn + ".tokaido.local")
	if err != nil {
		fmt.Println("There was an issue updating your hostsfile. Your hostsfile can be amended manually in order to enable this feature. See XXXXXXX for more information.")
		fmt.Println(err)
	}

	RebuildNginxConfigFile(hp)
	RestartContainer(proxy)
}
