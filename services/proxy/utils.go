package proxy

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"

	"path/filepath"
)

// CheckProxyUp ...
func CheckProxyUp() bool {
	_, err := utils.BashStringSplitOutput("curl --insecure -sSf " + GetProxyURL() + " > /dev/null")
	if err != nil {
		return false
	}

	return true
}

// GetProxyURL ...
func GetProxyURL() string {
	pn := conf.GetConfig().Tokaido.Project.Name

	return "https://" + pn + "." + constants.ProxyDomain + ":" + constants.ProxyPort
}

// `~/.tok/proxy/client/tls`
func getProxyClientTLSDir() string {
	return filepath.Join(getProxyClientDir(), "tls")
}

// `~/.tok/proxy/client/conf.d`
func getProxyClientConfdDir() string {
	return filepath.Join(getProxyClientDir(), "conf.d")
}

// `~/.tok/proxy/client/.env.production`
func getProxyClientGatsbyEnv() string {
	return filepath.Join(getProxyClientDir(), ".env.production")
}

// `~/.tok/proxy/client`
func getProxyClientDir() string {
	return filepath.Join(getProxyDir(), "client")
}

// `~/.tok/proxy`
func getProxyDir() string {
	return filepath.Join(getTokDir(), "proxy")
}

// `~/.tok`
func getTokDir() string {
	return filepath.Join(fs.HomeDir(), ".tok")
}

func getComposePath() string {
	return filepath.Join(getProxyDir(), "docker-compose.yml")
}

func buildDirectories() {
	utils.DebugString("setting up proxy directories")
	fs.Mkdir(getTokDir())
	fs.Mkdir(getProxyDir())
	fs.Mkdir(getProxyClientDir())
	fs.Mkdir(getProxyClientTLSDir())
	fs.Mkdir(getProxyClientConfdDir())
}
