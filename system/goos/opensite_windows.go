package goos

// TODO Windows

import (
	"github.com/ironstar-io/tokaido/utils"
)

// OpenSite - OSX Root executable
func OpenSite(url string) {
	utils.CommandSubstitution("powershell.exe", "start", url)
}
