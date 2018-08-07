package proxy

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/system/ssl"
)

// Setup ...
func Setup() {
	buildDirectories()

	ssl.Configure(getProxyClientTLSDir())

	GenerateProxyDockerCompose()
	DockerComposeUp()

	unison.CreateOrUpdatePrf(UnisonPort(), "proxy", getProxyClientDir())

	RebuildNginxConfigFile()
	RestartContainer("proxy")

	// AppendHostsfile() // Check then append
	// StartUnisonSyncSvc() // If not already started?
}
