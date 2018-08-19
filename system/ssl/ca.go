package ssl

import (
	"github.com/ironstar-io/tokaido/constants"

	"github.com/ironstar-io/tokaido/system/ssl/cfssl/cli"
	"github.com/ironstar-io/tokaido/system/ssl/cfssl/csr"
	"github.com/ironstar-io/tokaido/system/ssl/cfssl/initca"
)

// GenerateCA ...
func GenerateCA(c cli.Config) (CertificateGroupBody, error) {
	req := csr.CertificateRequest{
		KeyRequest: &csr.BasicKeyRequest{A: constants.KeyAlgorithm, S: constants.KeySize},
		CN:         constants.CommonName,
		Names: []csr.Name{{
			C:  constants.PkixCountry,
			ST: constants.PkixProvince,
			L:  constants.PkixLocality,
			O:  constants.PkixOrganization,
			OU: constants.PkixOrganizationalUnit,
		}},
	}

	var err error
	var key, csr, cert []byte
	cert, csr, key, err = initca.New(&req)
	if err != nil {
		return CertificateGroupBody{}, err
	}

	return CertificateGroupBody{
		CertificateRequest: csr,
		Certificate:        cert,
		Key:                key,
	}, nil
}
