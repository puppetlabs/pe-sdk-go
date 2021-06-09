package api

import "github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/client"

//Client is interface to the api client
type Client interface {
	GetClient() (*client.PuppetdbCli, error)
}
