package system

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/system/linux"
	"bitbucket.org/ironstar/tokaido-cli/system/osx"

	"log"
)

// OpenSite - Linux Root executable
func OpenSite() {
	httpsPort := docker.LocalPort("haproxy", "8443")
	if httpsPort == "" {
		log.Fatal("Unable to obtain the HTTPS port number. The HAProxy container may be broken")
	}

	url := "https://localhost:" + httpsPort
	if CheckOS() == "osx" {
		osx.OpenSite(url)
	} else {
		linux.OpenSite(url)
	}
}
