package proxy

import (
	"github.com/ironstar-io/tokaido/system/ssl"
)

// Setup ...
func Setup() {
	ssl.GenerateCert()           // Check then create
	GenerateProxyUnison()        // Check then create
	GenerateProxyDockerCompose() // Check then create

	RebuildNginxConfigFile() // Every run
	RestartProxyContainer()  // Every run

	AppendHostsfile() // Check then append

}
