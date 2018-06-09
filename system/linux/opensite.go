package linux

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"	
)

// OpenSite - Linux Root executable
func OpenSite(url string) {
	utils.NoFatalStdoutCmd("xdg-open", url)
}