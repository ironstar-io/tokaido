package osx

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

// OpenSite - Linux Root executable
func OpenSite(url string) {
	utils.NoFatalStdoutCmd("open", url)
}
