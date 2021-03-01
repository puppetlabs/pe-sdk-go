package app

import (
	"context"

	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/client/operations"
	"github.com/puppetlabs/pe-sdk-go/log"

	httptransport "github.com/go-openapi/runtime/client"
)

// PostImportFile uploads a puppetdb archive to the import endpoint of a puppet-db instance
func (puppetDb *PuppetDb) PostImportFile(filePath string) (*operations.PostImportOK, error) {
	stringToken, err := puppetDb.Token.Read()
	if err != nil {
		log.Debug(err.Error())
	}

	client, err := puppetDb.Client.GetClient()
	if err != nil {
		return nil, err
	}

	file, err := appFS.Open(filePath)
	if err != nil {
		return nil, err
	}

	apiKeyHeaderAuth := httptransport.APIKeyAuth("X-Authentication", "header", stringToken)
	postImportParameters := operations.NewPostImportParamsWithContext(context.Background())
	postImportParameters.SetArchive(file)

	return client.Operations.PostImport(postImportParameters, apiKeyHeaderAuth)
}
