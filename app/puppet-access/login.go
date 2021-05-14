package app

import (
	"context"
	"fmt"

	"github.com/puppetlabs/pe-sdk-go/app/puppet-access/api/client/operations"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-access/api/models"
	"github.com/puppetlabs/pe-sdk-go/log"
)

func errorHandling(err error) error {
	if du, ok := err.(*operations.LoginUnauthorized); ok {
		if du.Payload != nil {
			log.Debug(err.Error())
			return fmt.Errorf("Received an error while trying to request the token. The error was %v: %v", du.Payload.Kind, du.Payload.Msg)
		}
	}
	if du, ok := err.(*operations.LoginBadRequest); ok {
		if du.Payload != nil {
			log.Debug(err.Error())
			return fmt.Errorf("Received an error while trying to request the token. The error was %v: %v", du.Payload.Kind, du.Payload.Msg)
		}
	}
	if du, ok := err.(*operations.LoginDefault); ok {
		log.Debug(err.Error())
		return fmt.Errorf("Received an error while trying to request the token. The error with status code %d was %v", du.Code(), du.Payload)
	}

	return fmt.Errorf("Unknown error: %v", err.Error())
}

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
		log.Debug(err.Error())
		err := errorHandling(err)
		return "", err
	}

	if response.Payload.Token == "" {
		return "", fmt.Errorf("The response did not contain a token. Rerun with --debug to see full body")
	}

	return response.Payload.Token, nil
}
