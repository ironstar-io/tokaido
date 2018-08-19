// Package genkey implements the genkey command.
package genkey

import (
	"encoding/json"
	"errors"

	"github.com/ironstar-io/tokaido/system/ssl/cfssl/cli"
	"github.com/ironstar-io/tokaido/system/ssl/cfssl/csr"
	"github.com/ironstar-io/tokaido/system/ssl/cfssl/initca"
)

func genkeyMain(args []string, c cli.Config) (err error) {
	csrFile, args, err := cli.PopFirstArgument(args)
	if err != nil {
		return
	}
	if len(args) > 0 {
		return errors.New("only one argument is accepted, please check with usage")
	}

	csrFileBytes, err := cli.ReadStdin(csrFile)
	if err != nil {
		return
	}

	req := csr.CertificateRequest{
		KeyRequest: csr.NewBasicKeyRequest(),
	}
	err = json.Unmarshal(csrFileBytes, &req)
	if err != nil {
		return
	}

	if c.IsCA {
		var key, csrPEM, cert []byte
		cert, csrPEM, key, err = initca.New(&req)
		if err != nil {
			return
		}

		cli.PrintCert(key, csrPEM, cert)
	} else {
		if req.CA != nil {
			err = errors.New("ca section only permitted in initca")
			return
		}

		var key, csrPEM []byte
		g := &csr.Generator{Validator: Validator}
		csrPEM, key, err = g.ProcessRequest(&req)
		if err != nil {
			key = nil
			return
		}

		cli.PrintCert(key, csrPEM, nil)
	}
	return nil
}

// Validator does nothing and will never return an error. It exists because creating a
// csr.Generator requires a Validator.
func Validator(req *csr.CertificateRequest) error {
	return nil
}
