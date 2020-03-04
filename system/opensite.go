package system

import (
	"fmt"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drush"
	"github.com/ironstar-io/tokaido/services/proxy"
	"github.com/ironstar-io/tokaido/system/goos"
)

// OpenSite - Linux Root executable
func OpenSite(args []string, admin bool) {
	var port string
	services := map[string]string{
		"haproxy": "8443",
		"adminer": "8080",
		"mailhog": "8025",
		"nginx":   "8082",
		"varnish": "8081",
		"solr":    "8983",
	}

	path := ""
	// login as admin by getting a drush uli
	if admin == true {
		path = drush.GetAdminSignOnPath()
	}

	url := getProjectURL()

	if len(args) == 0 || args[0] == "haproxy" {
		port = docker.LocalPort("haproxy", services["haproxy"])
		goos.OpenSite(fmt.Sprintf("https://%s:%s%s", url, port, path))
		return
	}

	for _, arg := range args {
		if len(services[arg]) > 0 {
			port = docker.LocalPort(arg, services[arg])
			goos.OpenSite(fmt.Sprintf("http://%s:%s", url, port))
			return
		}
	}

}

// OpenTokaidoProxy ...
func OpenTokaidoProxy(admin bool) {
	var path string
	if admin == true {
		path = drush.GetAdminSignOnPath()
	}

	u := proxy.GetProxyURL()
	goos.OpenSite(u + path)
}

// getProjectURL ...
func getProjectURL() string {
	pn := conf.GetConfig().Tokaido.Project.Name

	return pn + "." + constants.ProxyDomain
}
