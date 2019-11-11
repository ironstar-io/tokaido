package goos

import (
	"github.com/ironstar-io/tokaido/utils"
)

// OpenSite - Open a URL using PowerShell
func OpenSite(url string) {
	utils.CommandSubstitution("powershell.exe", "start", url)
}
