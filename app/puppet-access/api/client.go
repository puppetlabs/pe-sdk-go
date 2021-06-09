package api

import "github.com/puppetlabs/pe-sdk-go/app/puppet-access/api/client"

//Client is interface to the api client
type Client interface {
	GetClient() (*client.PuppetAccess, error)
	EnableCN()
}
