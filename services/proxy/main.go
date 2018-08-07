package proxy

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/system/hostsfile"
	"github.com/ironstar-io/tokaido/system/ssl"
)

// Setup ...
func Setup() {
	pn := conf.GetConfig().Tokaido.Project.Name
	buildDirectories()

	ssl.Configure(getProxyClientTLSDir())

	GenerateProxyDockerCompose()
	DockerComposeUp()

	unison.CreateOrUpdatePrf(UnisonPort(), "proxy", getProxyClientDir())
	// unison.StartProxySyncSvc() // If not already started?

	hostsfile.AddEntry(pn + ".tokaido.local")

	RebuildNginxConfigFile()
	RestartContainer("proxy")
}
