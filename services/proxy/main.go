package proxy

import (
	"fmt"
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/system/hostsfile"
	"github.com/ironstar-io/tokaido/system/ssl"
)

const proxy = "proxy"

// Setup ...
func Setup() {
	pn := conf.GetConfig().Tokaido.Project.Name
	buildDirectories()

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

	RebuildNginxConfigFile()
	RestartContainer(proxy)
}
