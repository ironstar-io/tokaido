package tls

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"path/filepath"
	"time"

	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/logrusorgru/aurora"
)

// createCA creates a Certificate Authority and writes it's cert and key to disk
func createCA() (err error) {
	fullCertPath := filepath.Join(fs.HomeDir(), constants.TLSRoot, constants.CertificateAuthorityCertificatePath)
	fullKeyPath := filepath.Join(fs.HomeDir(), constants.TLSRoot, constants.CertificateAuthorityKeyPath)

	if fs.CheckExists(fullCertPath) {
		utils.DebugString("skipping creation of CA as one already exists in: " + fullCertPath)
		return nil // CA already exists, don't re-generate
	}

	// Create our CA Request
	req := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization: []string{constants.PkixOrganization},
			Country:      []string{constants.PkixCountry},
			Province:     []string{constants.PkixProvince},
			Locality:     []string{constants.PkixLocality},
			CommonName:   constants.CommonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // 10 years
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// Generate the CA Certificate and Key
	cert, key, err := generateCA(req)
	if err != nil {
		fmt.Println(aurora.Red("😓  Tokaido was not able to generate a trusted SSL certificate because of the following error:"))
		fmt.Println(err.Error())
		fmt.Println("    We'd love to help you fix this. Please visit https://docs.tokaido.io/en/docs/support.")
		return nil // try to carry on even though an error occurred
	}

	// Write the Cert and Key to disk
	fs.Mkdir(filepath.Join(fs.HomeDir(), constants.TLSRoot, "/ca"))
	fs.TouchOrReplace(fullCertPath, cert)
	fs.TouchOrReplace(fullKeyPath, key)

	// Install the Certificate Authority into the OS (if the user consents)
	configureTrustedCA(fullCertPath)

	return nil
}

// generateCA generates a Certificate Authority that can be used to sign any Tokaido environment or component
func generateCA(req *x509.Certificate) (certificate, key []byte, err error) {
	// create our privatekey
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, req, req, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	// pem encode
	caCertPEM := new(bytes.Buffer)
	pem.Encode(caCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	certificate = caCertPEM.Bytes()

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	key = caPrivKeyPEM.Bytes()

	return
}
