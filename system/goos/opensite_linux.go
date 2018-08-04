package goos

import (
	"github.com/ironstar-io/tokaido/utils"
)

// OpenSite - Linux Root executable
func OpenSite(url string) {
	utils.CommandSubstitution("xdg-open", url)
}
