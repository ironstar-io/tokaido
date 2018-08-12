// +build darwin

package ssl

import (
	"github.com/ironstar-io/tokaido/utils"

	"path/filepath"
	"strings"
)

// AddTrustedCertToKeychain ...
func AddTrustedCertToKeychain(certificate string) {
	lc := filepath.Join("/Library/Keychains/System.keychain")

	utils.BashStringCmd("sudo security add-trusted-cert -d -r trustRoot -k " + lc + " " + certificate)
}

// RemoveTrustedCertFromKeychain ...
func RemoveTrustedCertFromKeychain(certificate string) {
	utils.BashStringCmd("sudo security remove-trusted-cert -d " + certificate)
}

const certSuccess = "certificate verification successful"

// CertIsTrusted  ...
func CertIsTrusted(certificate string) bool {
	o, _ := utils.BashStringSplitOutput("security verify-cert -c " + certificate)

	return strings.Contains(o, certSuccess)
}
