package ssl

import (
	"path"
	"strings"

	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"

	"fmt"
	"path/filepath"

	"github.com/ironstar-io/tokaido/system/ssl/cfssl/cli"
	"github.com/ironstar-io/tokaido/system/ssl/cfssl/config"
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

// check if the certificate that exists is the older 'tokaido.local' style,
// and if so delete it
func removeOldCert(cp CertificateGroupPath) (err error) {
	utils.DebugString("Checking for presence of old-style 'tokaido.local' ssl certificate. A 'not found' error here is OK")
	cert := utils.CommandSubstitution("openssl", "x509", "-noout", "-subject", "-in", cp.Certificate)
	if strings.Contains(cert, "tokaido.local") {
		utils.DebugString("Removing an old-style 'tokaido.local' ssl certificate.")
		utils.DebugString("Removing old certificate file for tokaido.local.")
		err = fs.EmptyDir(path.Dir(cp.Certificate))
	}

	return
}

// FindOrCreateCA ...
func FindOrCreateCA(certPath string) (err error) {
	cp := GetCertificateGroupPath(certPath, constants.CAFilename)

	// Remove previous-generation tokaido.local certificates
	err = removeOldCert(cp)
	if err != nil {
		return err
	}

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
