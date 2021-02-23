package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/client/operations"
	"github.com/puppetlabs/pe-sdk-go/json"
	"github.com/puppetlabs/pe-sdk-go/log"

	httptransport "github.com/go-openapi/runtime/client"
)

//GetStatusWithErrorDetails will check code manager api status and add error details to the error text
func (puppetCode *PuppetCode) GetStatusWithErrorDetails() (*operations.GetStatusOK, error) {
	resp, err := puppetCode.getStatus()
	if err != nil {
		if du, ok := err.(*operations.GetStatusDefault); ok {
			if du.Payload != nil {
				log.Debug(err.Error())
				err = fmt.Errorf("[GET /status][%v] %v\n%v", du.Code(), du.Payload.Msg, json.PrettyPrintPayload(du.Payload.Details))
			}
		}
	}
	return resp, err
}

func (puppetCode *PuppetCode) getStatus() (*operations.GetStatusOK, error) {
	stringToken, err := puppetCode.Token.Read()
	if err != nil {
		log.Debug(err.Error())
		err = errors.New("Code Manager requires a token, please use `puppet access login` to generate a token")
		return nil, err
	}
	client, err := puppetCode.Client.GetClient()
	if err != nil {
		return nil, err
	}
	aPIKeyHeaderAuth := httptransport.APIKeyAuth("X-Authentication", "header", stringToken)
	getStatusParameters := operations.NewGetStatusParamsWithContext(context.Background())
	return client.Operations.GetStatus(getStatusParameters, aPIKeyHeaderAuth)
}
