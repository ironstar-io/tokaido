package ssl

import (
	"github.com/ironstar-io/tokaido/constants"

	"errors"

	"github.com/cloudflare/cfssl/api/generator"
	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/cli/sign"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/log"
	"github.com/cloudflare/cfssl/signer"
)

// GenerateCertificate ...
func GenerateCertificate(c cli.Config) (CertificateGroup, error) {
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
		Hosts: []string{constants.WildcardHost, constants.TopLevelHost},
	}

	c.Profile = "{\"tokaido\":{\"usages\":[\"signing\",\"key encipherment\",\"server auth\",\"client auth\"],\"expiry\":\"87600h\"}}"

	if c.CAFile == "" {
		return CertificateGroup{}, errors.New("CAFile property must be supplied")
	}

	if c.CAKeyFile == "" {
		return CertificateGroup{}, errors.New("CAKeyFile property must be supplied")
	}

	var err error
	var key, csrBytes []byte
	g := &csr.Generator{Validator: genkey.Validator}
	csrBytes, key, err = g.ProcessRequest(&req)
	if err != nil {
		return CertificateGroup{}, err
	}

	s, err := sign.SignerFromConfig(c)
	if err != nil {
		return CertificateGroup{}, err
	}

	var cert []byte
	signReq := signer.SignRequest{
		Request: string(csrBytes),
		Hosts:   signer.SplitHosts(c.Hostname),
		Profile: c.Profile,
		Label:   c.Label,
	}

	cert, err = s.Sign(signReq)
	if err != nil {
		return CertificateGroup{}, err
	}

	// This follows the Baseline Requirements for the Issuance and
	// Management of Publicly-Trusted Certificates, v.1.1.6, from the CA/Browser
	// Forum (https://cabforum.org). Specifically, section 10.2.3 ("Information
	// Requirements"), states:
	//
	// "Applicant information MUST include, but not be limited to, at least one
	// Fully-Qualified Domain Name or IP address to be included in the Certificate’s
	// SubjectAltName extension."
	if len(signReq.Hosts) == 0 && len(req.Hosts) == 0 {
		log.Warning(generator.CSRNoHostMessage)
	}

	return CertificateGroup{
		CertificateRequest: csrBytes,
		Certificate:        cert,
		Key:                key,
	}, nil
}
