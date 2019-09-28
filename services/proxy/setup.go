package proxy

import (
	"strconv"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"

	"fmt"
	"path/filepath"

	"github.com/logrusorgru/aurora"
)

const proxy = "proxy"

// Setup ...
func Setup() {
	utils.DebugString("setting up proxy directories")
	buildDirectories()

	utils.DebugString("configuring proxy TLS")
	// ssl.Configure(getProxyClientTLSDir())
	copyTLSCertificates()

	// If an existing proxy config exists, remove it to start again.
	if fs.CheckExists(fs.HomeDir() + "/.tok/proxy/docker-compose.yml") {
		DockerComposeRemoveProxy()
	}

	GenerateProxyDockerCompose()
	DockerComposeUp()

	if conf.GetConfig().Global.Syncservice == "unison" {
		ConfigureUnison()
	}

	ConfigureProjectNginx()

	removeLegacyYamanoteSetup()

	utils.DebugString("restarting proxy container")
	PullImages()
	RestartContainer(proxy)
}

// ConfigureProjectNginx ...
func ConfigureProjectNginx() {
	utils.DebugString("starting nginx proxy configuration")
	// Remove all existing nginx configuration files
	fs.EmptyDir(getProxyClientConfdDir())

	// Regenerate Nginx config for all active projects on this system
	for _, v := range conf.GetConfig().Global.Projects {
		h, err := docker.GetContainerIPFromProject("haproxy", v.Name)
		if err != nil {
			utils.DebugString("Skipping Nginx setup for project [" + v.Name + "]. Project haproxy container is not running")
			continue
		}

		if h == "" {
			fmt.Println(aurora.Red("    There was an error retrieving the IP address for the haproxy container in project [" + v.Name + "]."))
			fmt.Println(aurora.Red("    This is an unexpected error. Could you please create a GitHub Issue to help us fix this?"))
			continue
		}

		pp := constants.HTTPSProtocol + h + ":" + strconv.Itoa(constants.HaproxyInternalPort)
		nc := GenerateNginxConf(v.Name, constants.ProxyDomain, pp)
		np := filepath.Join(getProxyClientConfdDir(), v.Name+".conf")
		fs.Replace(np, nc)
	}

}

// Yamanote left a 'local.tokaido.io' nginx config file. This needs to be
// remove with the removal of Yamanote in 1.5.0, otherwise the proxy service
// won't start for existing Tokaido users.
func removeLegacyYamanoteSetup() {
	h := fs.HomeDir()

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

// copyTLSCertificates copies the proxy wildcard certificate from it's official
// location to where the proxy server can mount it
func copyTLSCertificates() {
	certSource := filepath.Join(fs.HomeDir(), constants.TLSRoot, constants.WildcardCertificatePath)
	certDest := filepath.Join(getProxyClientTLSDir(), "wildcard.crt")
	fs.Copy(certSource, certDest)

	keySource := filepath.Join(fs.HomeDir(), constants.TLSRoot, constants.WildcardKeyPath)
	keyDest := filepath.Join(getProxyClientTLSDir(), "wildcard.key")
	fs.Copy(keySource, keyDest)
}
