package app

import (
	"context"

	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/client/operations"
	"github.com/puppetlabs/pe-sdk-go/log"

	httptransport "github.com/go-openapi/runtime/client"
)

// GetStatus queries the status endpoint of a puppetdb instance
func (puppetDb *PuppetDb) GetStatus() (*operations.GetStatusOK, error) {
	stringToken, err := puppetDb.Token.Read()
	if err != nil {
		log.Debug(err.Error())
	}

	client, err := puppetDb.Client.GetClient()
	if err != nil {
		return nil, err
	}
	apiKeyHeaderAuth := httptransport.APIKeyAuth("X-Authentication", "header", stringToken)
	getStatusParameters := operations.NewGetStatusParamsWithContext(context.Background())
	return client.Operations.GetStatus(getStatusParameters, apiKeyHeaderAuth)
}
