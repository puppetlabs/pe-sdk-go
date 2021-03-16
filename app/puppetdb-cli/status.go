package app

import (
	"context"

	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/client/operations"

	httptransport "github.com/go-openapi/runtime/client"
)

// GetStatus queries the status endpoint of a puppetdb instance
func (puppetDb *PuppetDb) GetStatus() (interface{}, error) {
	client, err := puppetDb.Client.GetClient()
	if err != nil {
		return nil, err
	}
	apiKeyHeaderAuth := httptransport.APIKeyAuth("X-Authentication", "header", puppetDb.Token)
	getStatusParameters := operations.NewGetStatusParamsWithContext(context.Background())
	response, err := client.Operations.GetStatus(getStatusParameters, apiKeyHeaderAuth)
	if err != nil {
		return nil, err
	}
	return response.Payload, err
}
