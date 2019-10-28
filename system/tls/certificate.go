/*

functions in this file are responsible for creating and signing a TLS certificate using
the local Tokaido Certificate Authority (gen 3)

*/

package tls

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"time"

	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/logrusorgru/aurora"
)

// createWildcardCertificate creates the top-level wildcard certificate for use on the
// Tokaido proxy service
func createWildcardCertificate() (err error) {
	// Create our signing request
	req := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization: []string{"Tokaido Proxy Service"},
			Country:      []string{constants.PkixCountry},
			Province:     []string{constants.PkixProvince},
			Locality:     []string{constants.PkixLocality},
			CommonName:   constants.WildcardHost,
		},
		DNSNames:    []string{constants.WildcardHost},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(2, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	// Generate the CA Certificate and Key
	cert, key, err := generateCertificate(req)
	if err != nil {
		fmt.Println(aurora.Red("😓  Tokaido was not able to generate a trusted SSL certificate because of the following error:     "))
		fmt.Println(err.Error())
		fmt.Println("We'd love to help you fix this. Please visit https://docs.tokaido.io/en/docs/support.")
		return nil // try to carry on even though an error occurred
	}

	// Write the Cert and Key to disk
	fs.Mkdir(filepath.Join(GetTLSRootDir(), "/proxy"))
	fs.TouchOrReplace(filepath.Join(GetTLSRootDir(), constants.WildcardCertificatePath), cert)
	fs.TouchOrReplace(filepath.Join(GetTLSRootDir(), constants.WildcardKeyPath), key)

	return nil
}

// createProjectCertificate creates and signs a certificate matching a given common name
// for a single project, and then saves it in that projects .tok/local/tls directory
func createProjectCertificate(projectName, projectPath, commonName string) (err error) {
	// Create our signing request
	req := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization: []string{"Tokaido Proxy Service"},
			Country:      []string{constants.PkixCountry},
			Province:     []string{constants.PkixProvince},
			Locality:     []string{constants.PkixLocality},
			CommonName:   commonName,
		},
		DNSNames:    []string{commonName},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(2, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	// Generate the CA Certificate and Key
	cert, key, err := generateCertificate(req)
	if err != nil {
		fmt.Println(aurora.Red("😓  Tokaido was not able to generate a trusted SSL certificate because of the following error:    "))
		fmt.Println(err.Error())
		fmt.Println("We'd love to help you fix this. Please visit https://docs.tokaido.io/en/docs/support.")
		return nil // try to carry on even though an error occurred
	}

	// Write the Cert and Key to disk
	fs.Mkdir(filepath.Join(projectPath, ".tok/local/tls"))
	fs.TouchOrReplace(filepath.Join(projectPath, ".tok/local/tls/", commonName+".crt"), cert)
	fs.TouchOrReplace(filepath.Join(projectPath, ".tok/local/tls/", commonName+".key"), key)

	return nil
}

// generateCertificate takes a certificate signing request and signs it using the Tokaido CA
func generateCertificate(req *x509.Certificate) (certificate, key []byte, err error) {
	// Open the Tokaido CA cert and key from disk
	caCertPath := filepath.Join(GetTLSRootDir(), constants.CertificateAuthorityCertificatePath)
	caCertBytes, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println(aurora.Red("😓  Unable to generate a new certificate because of the following error while trying to open the CA Certificate at " + caCertPath + ":    "))
		fmt.Println(err.Error())
		fmt.Println("We'd love to help you fix this. Please visit https://docs.tokaido.io/en/docs/support.")
		return nil, nil, err
	}

	caKeyPath := filepath.Join(GetTLSRootDir(), constants.CertificateAuthorityKeyPath)
	caKeyBytes, err := ioutil.ReadFile(caKeyPath)
	if err != nil {
		fmt.Println(aurora.Red("😓  Unable to open generate a new certificate because of the following error while trying to open the CA Key at " + caKeyPath + ":    "))
		fmt.Println(err.Error())
		fmt.Println("We'd love to help you fix this. Please visit https://docs.tokaido.io/en/docs/support.")
		return nil, nil, err
	}

	// Decode our caKey PEM into a usable certificate format
	keyDecode, _ := pem.Decode(caKeyBytes)
	if keyDecode == nil {
		fmt.Println(aurora.Red("😓  Unable to open generate a new certificate because of the following error while trying to decode the CA Key:    "))
		fmt.Println(err.Error())
		fmt.Println("We'd love to help you fix this. Please visit https://docs.tokaido.io/en/docs/support.")
		return nil, nil, err
	}

	caPrivKey, err := x509.ParsePKCS1PrivateKey(keyDecode.Bytes)
	if err != nil {
		return nil, nil, err
	}

	// Decode our caCertificate PEM into a useable cert
	certDecode, _ := pem.Decode(caCertBytes)
	if certDecode == nil {
		fmt.Println(aurora.Red("😓  Unable to open generate a new certificate because of the following error:    "))
		fmt.Println(err.Error())
		fmt.Println("We'd love to help you fix this. Please visit https://docs.tokaido.io/en/docs/support.")
		return nil, nil, err
	}

	caCert, err := x509.ParseCertificate(certDecode.Bytes)
	if err != nil {
		return nil, nil, err
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, req, caCert, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	certificate = certPEM.Bytes()
	key = certPrivKeyPEM.Bytes()

	return certificate, key, nil
}
