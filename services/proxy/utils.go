package proxy

import (
	"github.com/ironstar-io/tokaido/system/fs"

	"path/filepath"
)

func getProxyClientDir() string {
	return filepath.Join(getProxyDir(), "client")
}

func getProxyDir() string {
	return filepath.Join(fs.HomeDir(), ".tok", "proxy")
}
