package goos

import (
	"github.com/ironstar-io/tokaido/utils"
	"github.com/ironstar-io/tokaido/system/wsl"
)

// OpenSite - Linux Root executable
func OpenSite(url string) {
	if wsl.IsWSL() {
		utils.CommandSubstitution("powershell.exe", "/c", "start", url)
	}

	utils.CommandSubstitution("xdg-open", url)
}
