// +build windows

package tls

// TODO Windows

import (
	// "github.com/ironstar-io/tokaido/utils"

	// "fmt"
	// "path/filepath"
	// "strings"
)

// configureTrustedCA ...
func configureTrustedCA(certificate string) {
	// if caIsTrusted(certificate) {
	// 	utils.DebugString("TLS certificate is already trusted. Nothing to do.")
	// 	return
	// }

	// utils.DebugString("Adding the Tokaido Certificate Authority to the macOS Keychain. An error here is OK and indicates that the certificate is just not trusted yet")

	// p := utils.ConfirmationPrompt("    Would you like to automatically add Tokaido's SSL authority to your macOS keychain? You may be prompted for elevated access", "y")
	// if p == false {
	// 	fmt.Println(`    The generated TLS certificates can be manually added to your keychain later. \n    See https://tokaido.io/docs/config/#adding-a-trusted-certificate for more information.`)
	// 	return
	// }

	// addTrustedCAToKeychain(certificate)
}

// addTrustedCAToKeychain ...
func addTrustedCAToKeychain(certificate string) {
	// lc := filepath.Join("/Library/Keychains/System.keychain")

	// utils.BashStringCmd("sudo security add-trusted-cert -d -r trustRoot -k " + lc + " " + certificate)
}

// removeTrustedCAFromKeychain ...
func removeTrustedCAFromKeychain(certificate string) {
	// utils.DebugString("removing legacy proxy CA from keychain")
	// utils.BashStringCmd("sudo security remove-trusted-cert -d " + certificate)
	// utils.BashStringCmd("sudo security delete-certificate -c local.tokaido.io")
}

const certSuccess = "certificate verification successful"

// caIsTrusted returns true if the certificate authority is successfully installed in the macOS keychain
func caIsTrusted(certificate string) bool {
	return false
	// o, _ := utils.BashStringSplitOutput("security verify-cert -c " + certificate)

	// return strings.Contains(o, certSuccess)
}
