package app

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/puppetlabs/pe-sdk-go/json"

	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/client"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/client/operations"
	mock_operations "github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/client/operations/testing"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/models"
	mock_api "github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/testing"
	match "github.com/puppetlabs/pe-sdk-go/app/puppet-code/testing"

	"github.com/stretchr/testify/assert"
)

func TestRunStatusFailsIfNoClient(t *testing.T) {
	assert := assert.New(t)
	errorMessage := "No client"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	api := mock_api.NewMockClient(ctrl)
	api.EXPECT().GetClient().Return(nil, errors.New(errorMessage))

	puppetCode := New()
	puppetCode.Token = "my token"
	puppetCode.Client = api
	_, receivedError := puppetCode.GetStatusWithErrorDetails()
	assert.EqualError(receivedError, errorMessage)
}

func TestRunStatusSucces(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetCode{
		Operations: operationsMock,
	}
	api.EXPECT().GetClient().Return(client, nil)

	result := &operations.GetStatusOK{
		Payload: "ok",
	}

	getStatusParameters := operations.NewGetStatusParamsWithContext(context.Background())
	operationsMock.EXPECT().GetStatus(getStatusParameters, match.XAuthenticationWriter(t, "my token")).Return(result, nil)

	puppetCode := New()
	puppetCode.Token = "my token"
	puppetCode.Client = api
	res, err := puppetCode.GetStatusWithErrorDetails()
	expected, _ := json.MarshalIndent(result.Payload, "", "  ")
	assert.Equal(expected, res)
	assert.Nil(err)
}

func TestRunStatusError(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetCode{
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

	puppetCode := New()
	puppetCode.Token = "my token"
	puppetCode.Client = api
	res, err := puppetCode.GetStatusWithErrorDetails()

	assert.Nil(res)
	assert.EqualError(err, "[GET /status][404] error message\n\"details\"")
}
