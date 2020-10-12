package system

import (
	"fmt"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/drush"
	"github.com/ironstar-io/tokaido/system/goos"
)

// OpenSite - Linux Root executable
func OpenSite(args []string, admin bool) {
	var port string
	services := map[string]string{
		"adminer": "8080",
		"mailhog": "8025",
		"nginx":   "8443",
		"varnish": "8081",
		"solr":    "8983",
	}

	path := ""
	// login as admin by getting a drush uli
	if admin == true {
		path = drush.GetAdminSignOnPath()
	}

	url := getProjectURL()

	if len(args) == 0 || args[0] == "nginx" {
		port = docker.LocalPort("nginx", services["nginx"])
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

// getProjectURL ...
func getProjectURL() string {
	pn := conf.GetConfig().Tokaido.Project.Name

	return pn + "." + constants.ProxyDomain
}
