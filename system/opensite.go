package system

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
	
	var handler string
	if utils.CheckOS() == "osx" {
		handler = "open"
	} else {
		handler = "xdg-open"
	}

	utils.NoFatalStdoutCmd(handler, "https://localhost:"+httpsPort)
}
