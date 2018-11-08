package system

import (
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/proxy"
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
		"varnish": "8081",
		"solr":    "8983",
	}

	if len(args) >= 1 {
		for _, arg := range args {
			if len(services[arg]) > 0 {
				p = docker.LocalPort(arg, services[arg])
				if arg == "haproxy" {
					goos.OpenSite("https://localhost:" + p)
					return
				}
				goos.OpenSite("http://localhost:" + p)
				return
			}
		}
	}

	OpenHaproxySite()
}

// OpenHaproxySite ...
func OpenHaproxySite() {
	if proxy.CheckProxyUp() == true {
		u := proxy.GetProxyURL()
		goos.OpenSite(u)
		return
	}

	p := docker.LocalPort("haproxy", "8443")
	goos.OpenSite("https://localhost:" + p)
}
