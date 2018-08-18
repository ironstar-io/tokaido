package ssl

import (
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"

	"fmt"
	"path/filepath"

	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/config"
)

// CertificateGroupBody ...
type CertificateGroupBody struct {
	CertificateRequest []byte
	Certificate        []byte
	Key                []byte
}

// CertificateGroupPath ...
type CertificateGroupPath struct {
	CertificateRequest string
	Certificate        string
	Key                string
}

// GetCertificateGroupPath ...
func GetCertificateGroupPath(certPath, filename string) CertificateGroupPath {
	return CertificateGroupPath{
		CertificateRequest: filepath.Join(certPath, filename+".csr"),
		Certificate:        filepath.Join(certPath, filename+".pem"),
		Key:                filepath.Join(certPath, filename+"-key.pem"),
	}
}

// Configure ...
func Configure(certPath string) {
	if err := FindOrCreateCA(certPath); err != nil {
		fmt.Printf("There was an issue creating a certificate authority: %s", err)
		fmt.Println("Skipping local HTTPS setup...")
		return
	}

	if err := FindOrCreateClientCert(certPath); err != nil {
		fmt.Printf("There was an issue creating client certificate: %s", err)
		fmt.Println("Skipping local HTTPS setup...")
		return
	}

	ConfigureTrustedCerts(filepath.Join(certPath, "proxy_ca.pem"))
}

// FindOrCreateClientCert ...
func FindOrCreateClientCert(certPath string) error {
	ca := GetCertificateGroupPath(certPath, constants.CAFilename)
	cp := GetCertificateGroupPath(certPath, constants.ClientCertFilename)
	if CheckCertsExist(cp) == true {
		return nil
	}

	cc, err := GenerateCertificate(cli.Config{
		CAFile:    ca.Certificate,
		CAKeyFile: ca.Key,
		Profile:   constants.SigningProfileName,
		CFG:       buildCFG(),
	})
	if err != nil {
		return err
	}

	fs.TouchByteArray(cp.CertificateRequest, cc.CertificateRequest)
	fs.TouchByteArray(cp.Certificate, cc.Certificate)
	fs.TouchByteArray(cp.Key, cc.Key)

	return nil
}

// FindOrCreateCA ...
func FindOrCreateCA(certPath string) error {
	cp := GetCertificateGroupPath(certPath, constants.CAFilename)
	if CheckCertsExist(cp) == true {
		return nil
	}

	ca, err := GenerateCA(cli.Config{})
	if err != nil {
		return err
	}

	fs.TouchByteArray(cp.CertificateRequest, ca.CertificateRequest)
	fs.TouchByteArray(cp.Certificate, ca.Certificate)
	fs.TouchByteArray(cp.Key, ca.Key)

	return nil
}

func buildCFG() *config.Config {
	u := []string{"signing", "key encipherment", "server auth", "client auth"}
	return &config.Config{
		Signing: &config.Signing{
			Default: &config.SigningProfile{
				Expiry:       10 * constants.OneYear,
				ExpiryString: constants.TenYearExpiryString,
				Usage:        u,
			},
			Profiles: map[string]*config.SigningProfile{constants.SigningProfileName: {
				Expiry:       10 * constants.OneYear,
				ExpiryString: constants.TenYearExpiryString,
				Usage:        u,
			}},
		},
	}
}

// CheckCertsExist ...
func CheckCertsExist(paths CertificateGroupPath) bool {
	if fs.CheckExists(paths.CertificateRequest) == false {
		return false
	}

	if fs.CheckExists(paths.Certificate) == false {
		return false
	}

	if fs.CheckExists(paths.Key) == false {
		return false
	}

	return true
}

// RemoveTrustedCert ...
func RemoveTrustedCert(certPath string) {
	c := filepath.Join(certPath, "proxy_ca.pem")

	if fs.CheckExists(c) == true {
		RemoveTrustedCertFromKeychain(c)
	}
}
