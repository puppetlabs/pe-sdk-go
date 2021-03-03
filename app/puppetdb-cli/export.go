package app

import (
	"context"
	"fmt"

	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/client/operations"
	"github.com/puppetlabs/pe-sdk-go/log"
	"github.com/spf13/afero"

	httptransport "github.com/go-openapi/runtime/client"
)

// createExportFile creates a file object from a provided path
func (puppetDb *PuppetDb) createExportFile(filePath string) (afero.File, error) {
	file, err := appFS.Create(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// GetExportFile queries the export endpoint of a puppet-db instance and saves the result in a file
func (puppetDb *PuppetDb) GetExportFile(filePath, anonymizationProfile string) (string, error) {
	stringToken, err := puppetDb.Token.Read()
	if err != nil {
		log.Debug(err.Error())
	}

	client, err := puppetDb.Client.GetClient()
	if err != nil {
		return "", err
	}

	file, err := puppetDb.createExportFile(filePath)
	if err != nil {
		return "", err
	}

	apiKeyHeaderAuth := httptransport.APIKeyAuth("X-Authentication", "header", stringToken)
	getExportParameters := operations.NewGetExportParamsWithContext(context.Background())
	getExportParameters.SetAnonymizationProfile(&anonymizationProfile)
	result, err := client.Operations.GetExport(getExportParameters, apiKeyHeaderAuth, file)
	if err != nil {
		appFS.Remove(filePath)
		return "", err
	}

	output := fmt.Sprint(result.Payload)
	return output, err
}
