package osx

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

// OpenSite - Linux Root executable
func OpenSite(url string) {
	utils.CommandSubstitution("open", url)
}
