package app

import (
	"context"
	"fmt"

	"strconv"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/client/operations"

	"github.com/puppetlabs/pe-sdk-go/json"
	"github.com/puppetlabs/pe-sdk-go/log"
)

// PostImportFile uploads a puppetdb archive to the import endpoint of a puppet-db instance
func (puppetDb *PuppetDb) PostImportFile(filePath string) (bool, error) {
	//var newOutput string
	stringToken, err := puppetDb.Token.Read()
	if err != nil {
		log.Debug(err.Error())
		return false, err
	}

	client, err := puppetDb.Client.GetClient()
	if err != nil {
		return false, err
	}

	file, err := appFS.Open(filePath)
	if err != nil {
		return false, err
	}
	fmt.Println(file.Name())
	apiKeyHeaderAuth := httptransport.APIKeyAuth("X-Authentication", "header", stringToken)
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

	output, err := json.MarshalIndent(resp.Payload, "", "")
	if err != nil {
		return false, err
	}

	finalOutput, err := strconv.ParseBool(string(output)[:len(output)-1])
	if err != nil {
		return false, err
	}
	return finalOutput, err

}
