// +build darwin

package ssl

import (
	"github.com/ironstar-io/tokaido/utils"

	"fmt"
	"path/filepath"
	"strings"
)

// ConfigureTrustedCerts ...
func ConfigureTrustedCerts(certificate string) {
	if CertIsTrusted(certificate) == true {
		return
	}

	fmt.Println()
	p := utils.ConfirmationPrompt("Would you like Tokaido to add the generated SSL certificate to your keychain? You may be prompted for elevated access", "y")
	if p == false {
		fmt.Println(`    The generated SSL certificates can be manually added to your keychain later. \n    See https://tokaido.io/docs/config/#adding-a-trusted-certificate for more information.`)
		return
	}

	AddTrustedCertToKeychain(certificate)
}

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
