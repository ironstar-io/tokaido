package ssl

import (
	"github.com/ironstar-io/tokaido/constants"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"os"
)

var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}

// GenerateCSR ...
func GenerateCSR() {
	keyBytes, _ := rsa.GenerateKey(rand.Reader, 1024)

	emailAddress := constants.CertEmailAddress
	subj := pkix.Name{
		CommonName:         constants.CommonName,
		Organization:       []string{constants.PkixOrganization},
		Locality:           []string{constants.PkixLocality},
		Province:           []string{constants.PkixProvince},
		Country:            []string{constants.PkixCountry},
		OrganizationalUnit: []string{constants.PkixOrganizationalUnit},
		ExtraNames: []pkix.AttributeTypeAndValue{
			{
				Type: oidEmailAddress,
				Value: asn1.RawValue{
					Tag:   asn1.TagIA5String,
					Bytes: []byte(emailAddress),
				},
			},
		},
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)
	pem.Encode(os.Stdout, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
}
