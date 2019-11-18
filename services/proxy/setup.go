package proxy

import (
	"strconv"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"

	"fmt"
	"path/filepath"

	"github.com/logrusorgru/aurora"
)

const proxy = "proxy"

// Setup ...
func Setup() {
	// Remove legacy Tokaido proxy config from 1.11 installations
	UpgradeTok111()

	// Configure the proxy path (if it doesn't exist already)
	buildDirectories()

	// Configure proxy TLS (if it doesn't exist already)
	copyTLSCertificates()

	// (re)-generate the nginx configuration for all active projects
	configureProjectNginx()

	// If the proxy server isn't running, generate a docker-compose.yml file and start it
	bootstrapProxy()

	if conf.GetConfig().Global.Syncservice == "unison" {
		configureUnison()
	}

	// Bump the nginx process with a HUP signal
	restartNginx()
}

// configureProjectNginx ...
func configureProjectNginx() {
	utils.DebugString("starting nginx proxy configuration")
	// Remove all existing nginx configuration files
	fs.EmptyDir(getProxyClientConfdDir())

	// Regenerate Nginx config for all active projects on this system
	for _, v := range conf.GetConfig().Global.Projects {
		h, err := getContainerProxyIP("haproxy", v.Name)
		if err != nil {
			utils.DebugString(err.Error())
			utils.DebugString("Skipping Nginx setup for project [" + v.Name + "]. Project haproxy container is not running")
			continue
		}

		if h == "" {
			fmt.Println(aurora.Red("    There was an error retrieving the IP address for the haproxy container in project [" + v.Name + "]."))
			fmt.Println(aurora.Red("    This is an unexpected error. Could you please create a GitHub Issue to help us fix this?"))
			continue
		}

		pp := constants.HTTPSProtocol + h + ":" + strconv.Itoa(constants.HaproxyInternalPort)
		nc := generateNginxConf(v.Name, constants.ProxyDomain, pp)
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
	utils.DebugString("copying wildcard cert and key")
	certSource := filepath.Join(fs.HomeDir(), constants.TLSRoot, constants.WildcardCertificatePath)
	certDest := filepath.Join(getProxyClientTLSDir(), "wildcard.crt")
	fs.Copy(certSource, certDest)

	keySource := filepath.Join(fs.HomeDir(), constants.TLSRoot, constants.WildcardKeyPath)
	keyDest := filepath.Join(getProxyClientTLSDir(), "wildcard.key")
	fs.Copy(keySource, keyDest)
}
