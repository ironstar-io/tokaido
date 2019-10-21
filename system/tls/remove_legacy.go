/*

functions in this file are responsible for finding and removing legacy SSL
certificate configurations.

Tokaido has had three certificate generations:

- "tokaido.local" was a simple self-signed cert for Tokaido versions before 1.0
  previous logic existed to remove this from Tokaido 1 to 1.11.
- "local.tokaido.io" was a Certificate Authority that can be added to browsers and
  the macOS keychain as a trusted certificate. It signed *.local.tokaido.io but
  relied on a very clunky implentation of CloudFlare's `cfssl` library
- "Tokaido Local Certficiate Authority" is the current implementation that uses
  smarter generation logic to create TLS certificates for both the main
  *.local.tokaido.io wildcard cert used by the Proxy container, but also for
  individual certs for things like 'haproxy' and other future TLS implementations
  like signed MariaDB connections.

*/

package tls

import (
	"github.com/ironstar-io/tokaido/system"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/ironstar-io/tokaido/system/console"

	"fmt"
	"path/filepath"
)

// removeLegacyCertificate will remove the local.tokaido.io CA from the system
func removeLegacyCertificate() {
	utils.DebugString("Checking for presence of legacy 'local.tokaido.io' Certificate Authority. A 'not found' error here is OK")

	legacyPath := filepath.Join(fs.HomeDir(), ".tok/proxy/client/tls/proxy_ca.pem")
	if fs.CheckExists(legacyPath) {
		console.Println("🎁  Tokaido 1.11 improves how Tokaido manages local TLS security, and will now reconfigure your system", "")

		if system.CheckOS() == "osx" {
			utils.DebugString("removing proxy CA from certificate chain")
			fmt.Println("    You may prompted for elevated access to allow Tokaido to refresh your macOS Keychain")
			removeTrustedCAFromKeychain(filepath.Join(legacyPath, "proxy_ca.pem"))
		}

		utils.DebugString("Removing directory [" + legacyPath + "]")
		fs.RemoveAll(legacyPath)
	}
}
