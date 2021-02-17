package app

import (
	"context"
	"fmt"

	"github.com/puppetlabs/pe-sdk-go/app/puppet-access/api/client/operations"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-access/api/models"
)

// Login method requests a token from the server
func (puppetAccess *PuppetAccess) Login() (string, error) {
	client, err := puppetAccess.Client.GetClient()
	if err != nil {
		return "", err
	}

	loginParameters := operations.NewLoginParamsWithContext(context.Background())
	body := models.LoginRequest{
		Login:    puppetAccess.Username,
		Password: puppetAccess.Password,
		Lifetime: puppetAccess.Lifetime,
		Label:    puppetAccess.Label,
	}
	loginParameters.SetBody(&body)
	response, err := client.Operations.Login(loginParameters)

	if err != nil {
		return "", err
	}

	if response.Payload.Token == "" {
		return "", fmt.Errorf("The response did not contain a token. Rerun with --debug to see full body")
	}

	return response.Payload.Token, nil
}
