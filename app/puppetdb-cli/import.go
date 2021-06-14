package app

import (
	"context"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/client/operations"
)

// PostImportFile uploads a puppetdb archive to the import endpoint of a puppet-db instance
func (puppetDb *PuppetDb) PostImportFile(filePath string) (bool, error) {
	client, err := puppetDb.Client.GetClient()
	if err != nil {
		return false, err
	}

	file, err := appFS.Open(filePath)
	if err != nil {
		return false, err
	}

	apiKeyHeaderAuth := httptransport.APIKeyAuth("X-Authentication", "header", puppetDb.Token)
	postImportParameters := operations.NewPostImportParamsWithContext(context.Background())
	postImportParameters.SetArchive(file)
	resp, err := client.Operations.PostImport(postImportParameters, apiKeyHeaderAuth)
	if err != nil {
		return false, err
	}

	//check if the payload is empty
	if resp.Payload == nil {
		return false, err
	}

	return resp.Payload.Ok, err

}
