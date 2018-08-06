package ssl

import (
	"log"
	"path/filepath"
)

// Configure ...
func Configure(certPath string) {
	c := filepath.Join(certPath, "tokaido.pem")
	k := filepath.Join(certPath, "tokaido-key.pem")
	h := []string{"*.tokaido.local", "tokaido.local"}

	err := CheckCerts(c, k)
	if err != nil {
		err = GenerateCerts(c, k, h)
		if err != nil {
			log.Fatal("Error: Unable to create https certs.")
		}
	}
}
