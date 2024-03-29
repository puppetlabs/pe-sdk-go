package app

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/client"
	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/client/operations"
	mock_operations "github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/client/operations/testing"
	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/models"
	mock_api "github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/testing"
	match "github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/testing"

	"github.com/stretchr/testify/assert"
)

func TestRunStatusFailsIfNoClient(t *testing.T) {
	assert := assert.New(t)
	errorMessage := "No client"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	api.EXPECT().GetClient().Return(nil, errors.New(errorMessage))

	puppetDb := New()
	puppetDb.Token = "my token"
	puppetDb.Client = api
	_, receivedError := puppetDb.GetStatus()
	assert.EqualError(receivedError, errorMessage)
}

func TestRunStatusSucces(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetdbCli{
		Operations: operationsMock,
	}
	api.EXPECT().GetClient().Return(client, nil)

	result := &operations.GetStatusOK{
		Payload: "ok",
	}

	getStatusParameters := operations.NewGetStatusParamsWithContext(context.Background())
	operationsMock.EXPECT().GetStatus(getStatusParameters, match.XAuthenticationWriter(t, "my token")).Return(result, nil)

	puppetDb := New()
	puppetDb.Token = "my token"
	puppetDb.Client = api
	res, err := puppetDb.GetStatus()

	assert.Equal("ok", res)
	assert.Nil(err)
}

func TestRunStatusError(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetdbCli{
		Operations: operationsMock,
	}
	api.EXPECT().GetClient().Return(client, nil)

	result := operations.NewGetStatusDefault(404)
	result.Payload = &models.Error{
		Msg:     "error message",
		Details: "details",
	}

	getStatusParameters := operations.NewGetStatusParamsWithContext(context.Background())
	operationsMock.EXPECT().GetStatus(getStatusParameters, match.XAuthenticationWriter(t, "my token")).Return(nil, result)

	puppetDb := New()
	puppetDb.Token = "my token"
	puppetDb.Client = api
	res, err := puppetDb.GetStatus()

	assert.Nil(res)
	assert.EqualError(err, "[GET /status/v1/services][404] getStatus default  &{Details:details Kind: Msg:error message}")
}
