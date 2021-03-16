package app

import (
	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api"
	"github.com/spf13/afero"
)

var appFS = afero.NewOsFs()

// PuppetDb interface
type PuppetDb struct {
	URL     string
	Token   string
	Cacert  string
	Cert    string
	Key     string
	Version string
	Client  api.Client
}

// NewWithConfig creates a puppet code application with configuration
func NewWithConfig(url, cacert, cert, key, token string) *PuppetDb {
	return &PuppetDb{
		URL:    url,
		Cacert: cacert,
		Cert:   cert,
		Key:    key,
		Token:  token,

		Client: api.NewClient(cacert, cert, key, url, token),
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
		Client: api.NewClient("", "", "", "", ""),
	}
}
