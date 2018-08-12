package system

import (
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/goos"
)

// OpenSite - Linux Root executable
func OpenSite(args []string) {
	var p string
	services := map[string]string{
		"haproxy": "8443",
		"adminer": "8080",
		"mailhog": "8025",
		"nginx":   "8082",
	}

	if len(args) >= 1 {
		for _, arg := range args {
			if len(services[arg]) > 0 {
				p = docker.LocalPort(arg, services[arg])
				proto := "http://"
				if arg == "haproxy" {
					proto = "https://"
				}
				goos.OpenSite(proto + "localhost:" + p)
			}
		}
	} else {
		p = docker.LocalPort("haproxy", "8443")
		goos.OpenSite("https://localhost:" + p)
	}
}
