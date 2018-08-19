package cli

import (
	"time"

	"github.com/ironstar-io/tokaido/system/ssl/cfssl/config"
	"github.com/ironstar-io/tokaido/system/ssl/cfssl/signer/universal"
)

// Config is a type to hold flag values used by cfssl commands.
type Config struct {
	Hostname          string
	CertFile          string
	CSRFile           string
	CAFile            string
	CAKeyFile         string
	TLSCertFile       string
	TLSKeyFile        string
	MutualTLSCAFile   string
	MutualTLSCNRegex  string
	TLSRemoteCAs      string
	MutualTLSCertFile string
	MutualTLSKeyFile  string
	KeyFile           string
	IntermediatesFile string
	CABundleFile      string
	IntBundleFile     string
	Address           string
	Port              int
	MinTLSVersion     string
	Password          string
	ConfigFile        string
	CFG               *config.Config
	Profile           string
	IsCA              bool
	RenewCA           bool
	IntDir            string
	Flavor            string
	Metadata          string
	Domain            string
	IP                string
	Remote            string
	Label             string
	AuthKey           string
	ResponderFile     string
	ResponderKeyFile  string
	Status            string
	Reason            string
	RevokedAt         string
	Interval          time.Duration
	List              bool
	Family            string
	Timeout           time.Duration
	Scanner           string
	CSVFile           string
	NumWorkers        int
	MaxHosts          int
	Responses         string
	Path              string
	CRL               string
	Usage             string
	PGPPrivate        string
	PGPName           string
	Serial            string
	CNOverride        string
	AKI               string
	DBConfigFile      string
	CRLExpiration     time.Duration
	Disable           string
}

// RootFromConfig returns a universal signer Root structure that can
// be used to produce a signer.
func RootFromConfig(c *Config) universal.Root {
	return universal.Root{
		Config: map[string]string{
			"cert-file": c.CAFile,
			"key-file":  c.CAKeyFile,
		},
		ForceRemote: c.Remote != "",
	}
}
