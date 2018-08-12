// +build darwin

package ssl

import (
	"fmt"
)

// ConfigureTrustedCerts ...
func ConfigureTrustedCerts(certificate string) {
	fmt.Println("The generated SSL certificates can be manually added to your keychain for easy HTTPS development. See XXXXXXX for more information.")
	return
}

// AddTrustedCertToKeychain ...
func AddTrustedCertToKeychain(certificate string) {
	return
}

// RemoveTrustedCertFromKeychain ...
func RemoveTrustedCertFromKeychain(certificate string) {
	fmt.Println("The generated SSL certificates need to be manually removed from your keychain if they were added during `tok up`. See XXXXXXX for more information.")
	return
}

// CertIsTrusted  ...
func CertIsTrusted(certificate string) bool {
	return false
}
