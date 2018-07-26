package goos

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

// OpenSite - OSX Root executable
func OpenSite(url string) {
	utils.CommandSubstitution("open", url)
}
