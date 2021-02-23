package app

import (
	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api"
	"github.com/puppetlabs/pe-sdk-go/token"
	"github.com/puppetlabs/pe-sdk-go/token/filetoken"
)

// PuppetCode FIXME
type PuppetCode struct {
	ServiceURL string
	Cacert     string
	TokenFile  string

	Client api.Client
	Token  token.Token
}

// NewWithConfig creates a puppet code application with configuration
func NewWithConfig(serviceURL, cacert, tokenFile string) *PuppetCode {
	return &PuppetCode{
		ServiceURL: serviceURL,
		Cacert:     cacert,
		TokenFile:  tokenFile,

		Client: api.NewClient(cacert, serviceURL),
		Token:  filetoken.NewFileToken(tokenFile),
	}
}

// New creates an unconfigured puppet code application
func New() *PuppetCode {
	return &PuppetCode{
		Token:  filetoken.NewFileToken(""),
		Client: api.NewClient("", ""),
	}
}
