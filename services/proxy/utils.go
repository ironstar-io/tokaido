package proxy

import (
	"github.com/ironstar-io/tokaido/system/fs"

	"path/filepath"
)

// `~/.tok/proxy/client/tls`
func getProxyClientTLSDir() string {
	return filepath.Join(getProxyClientDir(), "tls")
}

// `~/.tok/proxy/client/conf.d`
func getProxyClientConfdDir() string {
	return filepath.Join(getProxyClientDir(), "conf.d")
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

func buildDirectories() {
	fs.Mkdir(getTokDir())
	fs.Mkdir(getProxyDir())
	fs.Mkdir(getProxyClientDir())
	fs.Mkdir(getProxyClientTLSDir())
	fs.Mkdir(getProxyClientConfdDir())
}
