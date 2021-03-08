package app

import (
	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api"
)

// PuppetCode FIXME
type PuppetCode struct {
	ServiceURL string
	Cacert     string
	Client     api.Client
	Token      string
}

// NewWithConfig creates a puppet code application with configuration
func NewWithConfig(serviceURL, cacert, token string) *PuppetCode {
	return &PuppetCode{
		ServiceURL: serviceURL,
		Cacert:     cacert,

		Client: api.NewClient(cacert, serviceURL),
		Token:  token,
	}
}

// New creates an unconfigured puppet code application
func New() *PuppetCode {
	return &PuppetCode{
		Client: api.NewClient("", ""),
	}
}
