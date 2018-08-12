package ssl

import (
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"

	"log"
	"path/filepath"
)

// Configure ...
func Configure(certPath string) {
	c := filepath.Join(certPath, "tokaido.pem")
	k := filepath.Join(certPath, "tokaido-key.pem")
	h := []string{"*." + constants.ProxyDomain, constants.ProxyDomain}

	err := CheckCerts(c, k)
	if err != nil {
		err = GenerateCerts(c, k, h)
		if err != nil {
			log.Fatal("Error: Unable to create https certs.")
		}
	}

	ConfigureTrustedCerts(c)
}

// RemoveTrustedCert ...
func RemoveTrustedCert(certPath string) {
	c := filepath.Join(certPath, "tokaido.pem")

	if fs.CheckExists(c) == true {
		RemoveTrustedCertFromKeychain(c)
	}
}
