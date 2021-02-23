package app

import (
	"context"
	"errors"
	"fmt"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/client/operations"
	"github.com/puppetlabs/pe-sdk-go/json"
	"github.com/puppetlabs/pe-sdk-go/log"
)

// DeployWithErrorDetails will execute a deploy command and add error details to the error text
func (puppetCode *PuppetCode) DeployWithErrorDetails(args *DeployArgs) (*operations.DeployOK, error) {
	resp, err := puppetCode.deploy(args)
	if err != nil {
		if du, ok := err.(*operations.DeployDefault); ok {
			if du.Payload != nil {
				log.Debug(err.Error())
				err = fmt.Errorf("[POST /deploys][%d] %v\n%v", du.Code(), du.Payload.Msg, json.PrettyPrintPayload(du.Payload.Details))
			}
		}
	}
	return resp, err
}

func (puppetCode *PuppetCode) deploy(args *DeployArgs) (*operations.DeployOK, error) {
	deployConfig := puppetCode.getDeployConfig(args)

	stringToken, err := puppetCode.Token.Read()
	if err != nil {
		log.Debug(err.Error())
		err = errors.New("Code Manager requires a token, please use `puppet access login` to generate a token")
		return nil, err
	}
	aPIKeyHeaderAuth := httptransport.APIKeyAuth("X-Authentication", "header", stringToken)

	client, err := puppetCode.Client.GetClient()
	if err != nil {
		return nil, err
	}

	return client.Operations.Deploy(deployConfig, aPIKeyHeaderAuth)
}

//getDeployConfig creates a DeployParams based on command line arguments
func (puppetCode *PuppetCode) getDeployConfig(args *DeployArgs) *operations.DeployParams {
	deployParamenters := operations.NewDeployParamsWithContext(context.Background())
	body := operations.DeployBody{
		DryRun:       args.DryRun,
		DeployAll:    args.AllEnvironments,
		Wait:         args.Wait,
		Environments: args.Environments,
	}
	deployParamenters.SetBody(body)

	return deployParamenters
}

// DeployArgs should contain deploy command args
type DeployArgs struct {
	DryRun          bool
	AllEnvironments bool
	Wait            bool
	Environments    []string
}
