package system

import (
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drush"
	"github.com/ironstar-io/tokaido/services/proxy"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/goos"
)

// OpenSite - Linux Root executable
func OpenSite(args []string, admin bool) {
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

	OpenTokaidoProxy(admin)
}

// OpenTokaidoProxy ...
func OpenTokaidoProxy(admin bool) {
	var path string
	if admin == true {
		path = drush.GetAdminSignOnPath()
	}

	if proxy.CheckProxyUp() != true {
		console.Println(`
🤔  It looks like your Drupal site might not be working properly.
    Tokaido will open the site in your browser anyway, but you might want to run 'tok open haproxy'
    to bypass the Tokaido proxy service if this is causing problems.
    `, "")
	}

	u := proxy.GetProxyURL()
	goos.OpenSite(u + path)
}
