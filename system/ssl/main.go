package ssl

import (
	// "github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"

	// "log"
	"path/filepath"

	"github.com/cloudflare/cfssl/cli"
)

// CertificateGroup ...
type CertificateGroup struct {
	CertificateRequest []byte
	Certificate        []byte
	Key                []byte
}

// Configure ...
func Configure(certPath string) {

	FindOrCreateCA(certPath)

	FindOrCreateClientCert(certPath)
	// err = GenerateCerts(c, k, h)
	// if err != nil {
	// 	log.Fatal("Error: Unable to create https certs.")
	// }
	// }

	// ConfigureTrustedCerts(c)
}

// FindOrCreateClientCert ...
func FindOrCreateClientCert(certPath string) {
	// r := filepath.Join(certPath, "proxy_ca.csr")
	c := filepath.Join(certPath, "proxy_ca.pem")
	k := filepath.Join(certPath, "proxy_ca-key.pem")
	cc, _ := GenerateCertificate(cli.Config{CAFile: c, CAKeyFile: k})

	r2 := filepath.Join(certPath, "tokaido.csr")
	c2 := filepath.Join(certPath, "tokaido.pem")
	k2 := filepath.Join(certPath, "tokaido-key.pem")

	fs.TouchByteArray(r2, cc.CertificateRequest)
	fs.TouchByteArray(c2, cc.Certificate)
	fs.TouchByteArray(k2, cc.Key)
}

// FindOrCreateCA ...
func FindOrCreateCA(certPath string) CertificateGroup {
	r := filepath.Join(certPath, "proxy_ca.csr")
	c := filepath.Join(certPath, "proxy_ca.pem")
	k := filepath.Join(certPath, "proxy_ca-key.pem")

	ca, _ := GenerateCA(cli.Config{})

	fs.TouchByteArray(r, ca.CertificateRequest)
	fs.TouchByteArray(c, ca.Certificate)
	fs.TouchByteArray(k, ca.Key)

	return ca
}

// RemoveTrustedCert ...
func RemoveTrustedCert(certPath string) {
	c := filepath.Join(certPath, "proxy_ca.pem")

	if fs.CheckExists(c) == true {
		RemoveTrustedCertFromKeychain(c)
	}
}
