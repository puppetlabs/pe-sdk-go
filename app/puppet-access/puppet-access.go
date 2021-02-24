package app

import (
	"github.com/puppetlabs/pe-sdk-go/app/puppet-access/api"
)

// PuppetAccess interface
type PuppetAccess struct {
	Username string
	Password string
	Lifetime string
	URL      string
	Label    string
	Cacert   string
	Client   api.Client
}

// NewWithConfig creates a puppet access application with configuration
func NewWithConfig(username, password, lifetime, url, label, cacert string) *PuppetAccess {
	return &PuppetAccess{
		Username: username,
		Password: password,
		Lifetime: lifetime,
		URL:      url,
		Label:    label,
		Cacert:   cacert,
		Client:   api.NewClient(username, password, lifetime, url, label, cacert),
	}
}

// NewWithMinimalConfig creates a puppet access application with minimal configuration
func NewWithMinimalConfig(username, password, url, cacert string) *PuppetAccess {
	return &PuppetAccess{
		Username: username,
		Password: password,
		URL:      url,
		Cacert:   cacert,
		Client:   api.NewClient(username, password, "", url, "", cacert),
	}
}

// New creates an unconfigured puppet-db application
func New() *PuppetAccess {
	return &PuppetAccess{
		Client: api.NewClient("", "", "", "", "", ""),
	}
}
