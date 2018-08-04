package system

import (
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/goos"

	"log"
)

// OpenSite - Linux Root executable
func OpenSite() {
	httpsPort := docker.LocalPort("haproxy", "8443")
	if httpsPort == "" {
		log.Fatal("Unable to obtain the HTTPS port number. The HAProxy container may be broken")
	}

	goos.OpenSite("https://localhost:" + httpsPort)
}
