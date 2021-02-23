package app

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/client"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/client/operations"
	mock_operations "github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/client/operations/testing"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/models"
	mock_api "github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/testing"
	match "github.com/puppetlabs/pe-sdk-go/app/puppet-code/testing"
	mock_token "github.com/puppetlabs/pe-sdk-go/token/testing"

	"github.com/stretchr/testify/assert"
)

func TestRunDeployFailsIfNoToken(t *testing.T) {
	assert := assert.New(t)
	errorMessage := "Code Manager requires a token, please use `puppet access login` to generate a token"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	token := mock_token.NewMockToken(ctrl)
	api := mock_api.NewMockClient(ctrl)

	args := DeployArgs{
		DryRun:          true,
		AllEnvironments: true,
		Wait:            true,
		Environments:    []string{"environment"},
	}

	token.EXPECT().Read().Return("", errors.New(errorMessage))

	puppetCode := New()
	puppetCode.Token = token
	puppetCode.Client = api
	_, receivedError := puppetCode.DeployWithErrorDetails(&args)
	assert.EqualError(receivedError, errorMessage)
}

func TestRunDeployFailsIfNoClient(t *testing.T) {
	assert := assert.New(t)
	errorMessage := "No client"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	token := mock_token.NewMockToken(ctrl)
	api := mock_api.NewMockClient(ctrl)

	args := DeployArgs{
		DryRun:          true,
		AllEnvironments: true,
		Wait:            true,
		Environments:    []string{"environment"},
	}

	token.EXPECT().Read().Return("my token", nil)
	api.EXPECT().GetClient().Return(nil, errors.New(errorMessage))

	puppetCode := New()
	puppetCode.Token = token
	puppetCode.Client = api
	_, receivedError := puppetCode.DeployWithErrorDetails(&args)
	assert.EqualError(receivedError, errorMessage)
}

func TestDeployStatusSucces(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	token := mock_token.NewMockToken(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetCode{
		Operations: operationsMock,
	}
	args := DeployArgs{
		DryRun:          true,
		AllEnvironments: true,
		Wait:            true,
		Environments:    []string{"environment"},
	}

	api.EXPECT().GetClient().Return(client, nil)
	token.EXPECT().Read().Return("my token", nil)

	result := operations.NewDeployOK()

	deployParamenters := operations.NewDeployParamsWithContext(context.Background())
	body := operations.DeployBody{
		DryRun:       args.DryRun,
		DeployAll:    args.AllEnvironments,
		Wait:         args.Wait,
		Environments: args.Environments,
	}
	deployParamenters.SetBody(body)

	operationsMock.EXPECT().Deploy(deployParamenters, match.XAuthenticationWriter(t, "my token")).Return(result, nil)

	puppetCode := New()
	puppetCode.Token = token
	puppetCode.Client = api
	res, err := puppetCode.DeployWithErrorDetails(&args)

	assert.Equal(operations.NewDeployOK(), res)
	assert.Nil(err)
}

func TestDeployStatusWithError(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	token := mock_token.NewMockToken(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetCode{
		Operations: operationsMock,
	}
	args := DeployArgs{
		DryRun:          true,
		AllEnvironments: true,
		Wait:            true,
		Environments:    []string{"environment"},
	}

	api.EXPECT().GetClient().Return(client, nil)
	token.EXPECT().Read().Return("my token", nil)

	result := operations.NewDeployDefault(404)
	result.Payload = &models.Error{
		Msg:     "error message",
		Details: "details",
	}

	deployParamenters := operations.NewDeployParamsWithContext(context.Background())
	body := operations.DeployBody{
		DryRun:       args.DryRun,
		DeployAll:    args.AllEnvironments,
		Wait:         args.Wait,
		Environments: args.Environments,
	}
	deployParamenters.SetBody(body)

	operationsMock.EXPECT().Deploy(deployParamenters, match.XAuthenticationWriter(t, "my token")).Return(nil, result)

	puppetCode := New()
	puppetCode.Token = token
	puppetCode.Client = api
	res, err := puppetCode.DeployWithErrorDetails(&args)

	assert.Nil(res)
	assert.EqualError(err, "[POST /deploys][404] error message\n\"details\"")
}
