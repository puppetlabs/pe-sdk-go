package app

import (
	"github.com/puppetlabs/pe-sdk-go/app/puppet-access/api"
)

//PuppetAccessCfg config
type PuppetAccessCfg struct {
	Username string
	Password string
	Lifetime string
	URL      string
	Label    string
	Cacert   string

	UseCNVerification bool
}

// PuppetAccess interface
type PuppetAccess struct {
	Username string
	Password string
	Lifetime string
	Label    string
	Client   api.Client
}

// NewWithConfig creates a puppet access application with configuration
func NewWithConfig(cfg PuppetAccessCfg) *PuppetAccess {
	apiCfg := api.SwaggerClientCfg{
		Username:          cfg.Username,
		Password:          cfg.Password,
		Lifetime:          cfg.Lifetime,
		URL:               cfg.URL,
		Label:             cfg.Label,
		Cacert:            cfg.Cacert,
		UseCNVerification: cfg.UseCNVerification,
	}

	return &PuppetAccess{
		Username: cfg.Username,
		Password: cfg.Password,
		Lifetime: cfg.Lifetime,
		Label:    cfg.Label,
		Client:   api.NewClient(apiCfg),
	}
}

// NewWithMinimalConfig creates a puppet access application with minimal configuration
func NewWithMinimalConfig(cfg PuppetAccessCfg) *PuppetAccess {
	apiCfg := api.SwaggerClientCfg{
		Username: cfg.Username,
		Password: cfg.Password,
		URL:      cfg.URL,
		Cacert:   cfg.Cacert,
	}

	return &PuppetAccess{
		Username: cfg.Username,
		Password: cfg.Password,

		Client: api.NewClient(apiCfg),
	}
}

// New creates an unconfigured puppet-db application
func New() *PuppetAccess {
	return &PuppetAccess{
		Username: "",
		Client:   api.NewClient(api.SwaggerClientCfg{}),
	}
}

// EnableCN enables the CN verification
func (pa *PuppetAccess) EnableCN() {
	pa.Client.EnableCN()
}
