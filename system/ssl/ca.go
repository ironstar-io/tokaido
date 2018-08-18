package ssl

import (
	"github.com/ironstar-io/tokaido/constants"
	// "github.com/ironstar-io/tokaido/system/fs"

	// "errors"
	// "fmt"

	// "github.com/cloudflare/cfssl/api/generator"
	"github.com/cloudflare/cfssl/cli"
	// "github.com/cloudflare/cfssl/cli/genkey"
	// "github.com/cloudflare/cfssl/cli/sign"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	// "github.com/cloudflare/cfssl/log"
	// "github.com/cloudflare/cfssl/signer"
)

// GenerateCA ...
func GenerateCA(c cli.Config) (CertificateGroup, error) {
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

	c.Profile = "{\"tokaido\":{\"usages\":[\"signing\",\"key encipherment\",\"server auth\",\"client auth\"],\"expiry\":\"87600h\"}}"

	var err error
	var key, csr, cert []byte
	cert, csr, key, err = initca.New(&req)
	if err != nil {
		return CertificateGroup{}, err
	}

	return CertificateGroup{
		CertificateRequest: csr,
		Certificate:        cert,
		Key:                key,
	}, nil
}
