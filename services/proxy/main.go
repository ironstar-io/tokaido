package proxy

import (
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/system/ssl"
)

// Setup ...
func Setup() {
	buildDirectories()

	ssl.GenerateCert() // Check then create
	// GenerateProxyDockerCompose() // Check then create
	unison.CreateOrUpdatePrf("TODO: Proxy Container Local Port", "proxy", getProxyClientDir())

	// RebuildNginxConfigFile() // Every run
	// RestartProxyContainer()  // Every run

	// AppendHostsfile() // Check then append

}
