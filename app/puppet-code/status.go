package app

import (
	"context"
	"fmt"

	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/client/operations"
	"github.com/puppetlabs/pe-sdk-go/json"
	"github.com/puppetlabs/pe-sdk-go/log"

	httptransport "github.com/go-openapi/runtime/client"
)

//GetStatusWithErrorDetails will check code manager api status and add error details to the error text
func (puppetCode *PuppetCode) GetStatusWithErrorDetails() ([]byte, error) {
	resp, err := puppetCode.getStatus()
	if err != nil {
		if du, ok := err.(*operations.GetStatusDefault); ok {
			if du.Payload != nil {
				log.Debug(err.Error())
				err = fmt.Errorf("[GET /status][%v] %v\n%v", du.Code(), du.Payload.Msg, json.PrettyPrintPayload(du.Payload.Details))
			}
		}
		return nil, err
	}
	output, err := json.MarshalIndent(resp.Payload, "", "  ")

	return output, err
}

func (puppetCode *PuppetCode) getStatus() (*operations.GetStatusOK, error) {
	client, err := puppetCode.Client.GetClient()
	if err != nil {
		return nil, err
	}
	aPIKeyHeaderAuth := httptransport.APIKeyAuth("X-Authentication", "header", puppetCode.Token)
	getStatusParameters := operations.NewGetStatusParamsWithContext(context.Background())
	return client.Operations.GetStatus(getStatusParameters, aPIKeyHeaderAuth)
}
