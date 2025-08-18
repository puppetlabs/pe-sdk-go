package app

import (
	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api"
)

// PuppetCodeCfg config
type PuppetCodeCfg struct {
	ServiceURL string
	Cacert     string
	Token      string

	UseCNVerification bool
}

// PuppetCode FIXME
type PuppetCode struct {
	Token  string
	Client api.Client
}

// NewWithConfig creates a puppet code application with configuration
func NewWithConfig(cfg PuppetCodeCfg) *PuppetCode {
	apiCfg := api.SwaggerClientCfg{
		ServiceURL:        cfg.ServiceURL,
		Cacert:            cfg.Cacert,
		Token:             cfg.Token,
		UseCNVerification: cfg.UseCNVerification,
	}
	return &PuppetCode{
		Token:  cfg.Token,
		Client: api.NewClient(apiCfg),
	}
}

// New creates an unconfigured puppet code application
func New() *PuppetCode {
	return &PuppetCode{
		Token:  "",
		Client: api.NewClient(api.SwaggerClientCfg{}),
	}
}
