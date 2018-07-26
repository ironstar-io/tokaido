package goos

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

// OpenSite - Windows
func OpenSite(url string) {
	utils.CommandSubstitution("cmd", "/c", "start", url)
}
