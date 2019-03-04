package proxy

import (
	"log"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/ssl"

	"fmt"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

const proxy = "proxy"

// Setup ...
func Setup() {
	buildDirectories()

	ssl.Configure(getProxyClientTLSDir())

	GenerateProxyDockerCompose()
	DockerComposeUp()

	ConfigureProjectNginx()

	removeLegacyYamanoteSetup()

	RestartContainer(proxy)
}

// ConfigureProjectNginx ...
func ConfigureProjectNginx() {
	h, err := docker.GetContainerIP("haproxy")
	if err != nil {
		fmt.Printf("%s. Skipping HTTPS proxy setup...\n", err)
		return
	}

	if h == "" {
		fmt.Println("The haproxy container doesn't appear to be running. Skipping HTTPS proxy setup...")
		return
	}

	pp := constants.HTTPSProtocol + h + ":" + constants.HaproxyInternalPort

	pn := conf.GetConfig().Tokaido.Project.Name
	do := pn + `.` + constants.ProxyDomain

	nc := GenerateNginxConf(do, pp)

	np := filepath.Join(getProxyClientConfdDir(), pn+".conf")
	fs.Replace(np, nc)
}

// Yamanote left a 'local.tokaido.io' nginx config file. This needs to be
// remove with the removal of Yamanote in 1.5.0, otherwise the proxy service
// won't start for existing Tokaido users.
func removeLegacyYamanoteSetup() {
	h, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Unable to resolve home directory: %v", err)
	}

	// Remove the yamanote config from when we used DNS auto-resolving "local.tokaido.io"
	p := h + "/.tok/proxy/client/conf.d/local.tokaido.io.conf"

	if fs.CheckExists(p) {
		fs.Remove(p)
	}

	// Remove yamanote config from < 1.2.0, when we used /etc/hosts entries
	p = h + "/.tok/proxy/client/conf.d/tokaido.local.conf"

	if fs.CheckExists(p) {
		fs.Remove(p)
	}

}
