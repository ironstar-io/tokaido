package osx

import (
	"bitbucket.org/ironstar/tokaido-cli/services/docker"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"log"	
	
)

// OpenSite - Linux Root executable
func OpenSite() {
	httpsPort := docker.LocalPort("haproxy", "8443")
	if httpsPort == "" {
		log.Fatal("Unable to obtain the HTTPS port number. The HAProxy container may be broken")
	}
	utils.NoFatalStdoutCmd("open", "https://localhost:"+httpsPort)
}
