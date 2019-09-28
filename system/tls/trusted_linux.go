// +build linux

/*

functions in this file are responsible for installing the local Tokaido CA as a trusted
certificate authority in Linux.

As this functionality is yet to be fully developed, instead it will display messages
guiding the user on how to install the certificate manually.

*/

package tls

import (
	"fmt"
)

// configureTrustedCA ...
func configureTrustedCA(certificate string) {
	fmt.Println("    The generated SSL certificates can be manually added to your keychain for easy HTTPS development. \n    See https://tokaido.io/docs/config/#adding-a-trusted-certificate for more information.")
	return
}

// addTrustedCAToKeychain ...
func addTrustedCAToKeychain(certificate string) {
	return
}

// removeTrustedCAFromKeychain ...
func removeTrustedCAFromKeychain(certificate string) {
	fmt.Println("    The generated SSL certificates need to be manually removed from your keychain if they were added during `tok up`. \n    See https://tokaido.io/docs/config/#adding-a-trusted-certificate for more information.")
	return
}

// caIsTrusted  ...
func caIsTrusted(certificate string) bool {
	return false
}
