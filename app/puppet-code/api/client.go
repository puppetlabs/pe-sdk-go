package api

import "github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/client"

//Client is interface to the api client
type Client interface {
	GetClient() (*client.PuppetCode, error)
}
