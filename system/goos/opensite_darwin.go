package goos

import (
	"github.com/ironstar-io/tokaido/utils"
)

// OpenSite - MacOS Root executable
func OpenSite(url string) {
	utils.CommandSubstitution("open", url)
}
