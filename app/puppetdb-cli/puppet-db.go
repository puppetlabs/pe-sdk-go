package app

import (
	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api"
	"github.com/spf13/afero"
)

var appFS = afero.NewOsFs()

// PuppetDbCfg is a pupperDB client configuration
type PuppetDbCfg struct {
	URL    string
	Token  string
	Cacert string
	Cert   string
	Key    string

	UseCNVerification bool
}

// PuppetDb interface
type PuppetDb struct {
	Token   string
	Version string

	Client api.Client
}

// NewWithConfig creates a puppet code application with configuration
func NewWithConfig(cfg PuppetDbCfg) *PuppetDb {
	apiCfg := api.SwaggerClientCfg{
		Cacert:            cfg.Cacert,
		Cert:              cfg.Cert,
		Key:               cfg.Key,
		URL:               cfg.URL,
		Token:             cfg.Token,
		UseCNVerification: cfg.UseCNVerification,
	}

	return &PuppetDb{
		Token:  cfg.Token,
		Client: api.NewClientWithConfig(apiCfg),
	}
}

// NewPuppetDbApp FIXME
func NewPuppetDbApp(version string) *PuppetDb {
	return &PuppetDb{
		Version: version,
	}
}

// New creates an unconfigured puppet-db application
func New() *PuppetDb {
	return &PuppetDb{
		Token:  "",
		Client: api.NewClientWithConfig(api.SwaggerClientCfg{}),
	}
}

// EnableCN enables the CN verification
func (pdb *PuppetDb) EnableCN() {
	pdb.Client.EnableCN()
}
