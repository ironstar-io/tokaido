package goos

import (
	"github.com/ironstar-io/tokaido/utils"
)

// OpenSite - OSX Root executable
func OpenSite(url string) {
	utils.CommandSubstitution("open", url)
}
