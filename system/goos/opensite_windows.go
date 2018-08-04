package goos

import (
	"github.com/ironstar-io/tokaido/utils"
)

// OpenSite - Windows
func OpenSite(url string) {
	utils.CommandSubstitution("cmd", "/c", "start", url)
}
