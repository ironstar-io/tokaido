package ssl

import (
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/utils"

	"fmt"
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

	if CertIsTrusted(c) == true {
		return
	}

	fmt.Println()
	p := utils.ConfirmationPrompt("Would you like Tokaido to add the generated SSL certificate to your keychain? You may be prompted for elevated access", "n")
	if p == false {
		fmt.Println(`The generated SSL certificates can be manually added to your keychain later. See XXXXXXX for more information.`)
		return
	}

	AddTrustedCertToKeychain(c)
}
