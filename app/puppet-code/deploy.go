package app

import (
	"context"
	"fmt"
	"os"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/client/operations"
	"github.com/puppetlabs/pe-sdk-go/json"
	"github.com/puppetlabs/pe-sdk-go/log"
)

// DeployWithErrorDetails will execute a deploy command and add error details to the error text
func (puppetCode *PuppetCode) DeployWithErrorDetails(args *DeployArgs) ([]byte, error) {
	resp, err := puppetCode.deploy(args)
	if err != nil {
		if du, ok := err.(*operations.DeployDefault); ok {
			if du.Payload != nil {
				log.Debug(err.Error())
				err = fmt.Errorf("[POST /deploys][%d] %v\n%v", du.Code(), du.Payload.Msg, json.PrettyPrintPayload(du.Payload.Details))
			}

		}
		return nil, err
	}
	output := writeDeployResult(args, resp.Payload)
	return output, err
}

func (puppetCode *PuppetCode) deploy(args *DeployArgs) (*operations.DeployOK, error) {
	deployConfig := puppetCode.getDeployConfig(args)

	aPIKeyHeaderAuth := httptransport.APIKeyAuth("X-Authentication", "header", puppetCode.Token)

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

func writeDeployResult(args *DeployArgs, payload []*operations.DeployOKBodyItems0) []byte {
	fmt.Fprintf(os.Stderr, "Found %d environments.\n", len(payload))

	if args.DryRun {
		separator := ""
		environments := ""
		for _, v := range payload {
			environments = environments + separator + v.Environment
			separator = ", "
		}
		log.Info(fmt.Sprintf("Found the following environments: %s", environments))
	}
	resultPayload, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Error(err.Error())
	}
	return resultPayload
}
