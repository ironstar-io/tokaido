package proxy

import (
	"path/filepath"

	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
)

// UpgradeTok111 removes previous config for Tokaido versions < 1.11
func UpgradeTok111() {
	fp := filepath.Join(fs.HomeDir(), ".tok/proxy/docker-compose.yml")
	// Only upgrade if we haven't previously upgraded
	if !fs.CheckExists(fp) || !fs.Contains(fp, "tokaido/proxy:latest") {
		return
	}

	utils.DebugString("removing legacy proxy configuration")
	// Stop and remove the proxy docker-compose stack
	dockerComposeRemoveProxy()

	// Delete the proxy docker-compose file
	fs.Remove(fp)
}
