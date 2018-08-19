package ssl

import (
	"github.com/ironstar-io/tokaido/constants"

	"errors"

	"github.com/ironstar-io/tokaido/system/ssl/cfssl/api/generator"
	"github.com/ironstar-io/tokaido/system/ssl/cfssl/cli"
	"github.com/ironstar-io/tokaido/system/ssl/cfssl/cli/genkey"
	"github.com/ironstar-io/tokaido/system/ssl/cfssl/cli/sign"
	"github.com/ironstar-io/tokaido/system/ssl/cfssl/csr"
	"github.com/ironstar-io/tokaido/system/ssl/cfssl/log"
	"github.com/ironstar-io/tokaido/system/ssl/cfssl/signer"
)

// GenerateCertificate ...
func GenerateCertificate(c cli.Config) (CertificateGroupBody, error) {
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

	if c.CAFile == "" {
		return CertificateGroupBody{}, errors.New("CAFile property must be supplied")
	}

	if c.CAKeyFile == "" {
		return CertificateGroupBody{}, errors.New("CAKeyFile property must be supplied")
	}

	var err error
	var key, csrBytes []byte
	g := &csr.Generator{Validator: genkey.Validator}
	csrBytes, key, err = g.ProcessRequest(&req)
	if err != nil {
		return CertificateGroupBody{}, err
	}

	s, err := sign.SignerFromConfig(c)
	if err != nil {
		return CertificateGroupBody{}, err
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
		return CertificateGroupBody{}, err
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

	return CertificateGroupBody{
		CertificateRequest: csrBytes,
		Certificate:        cert,
		Key:                key,
	}, nil
}
