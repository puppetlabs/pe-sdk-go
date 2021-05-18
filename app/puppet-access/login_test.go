package app

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-access/api/client"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-access/api/client/operations"
	mock_operations "github.com/puppetlabs/pe-sdk-go/app/puppet-access/api/client/operations/testing"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-access/api/models"
	mock_api "github.com/puppetlabs/pe-sdk-go/app/puppet-access/api/testing"
	"github.com/stretchr/testify/assert"
)

func TestRunLoginFailsIfNoClient(t *testing.T) {
	assert := assert.New(t)
	errorMessage := "No client"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	api.EXPECT().GetClient().Return(nil, errors.New(errorMessage))

	puppetAccess := New()
	puppetAccess.Client = api

	_, receivedError := puppetAccess.Login()
	assert.EqualError(receivedError, errorMessage)
}

func TestRunLoginFailure(t *testing.T) {
	assert := assert.New(t)

	login := "username"
	password := "pass"
	lifetime := "10m"
	label := "test_token"
	errorMessage := "Salam"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetAccess{
		Operations: operationsMock,
	}

	api.EXPECT().GetClient().Return(client, nil)

	postLoginParameters := operations.NewLoginParamsWithContext(context.Background())
	body := models.LoginRequest{
		Login:    login,
		Password: password,
		Lifetime: lifetime,
		Label:    label,
	}
	postLoginParameters.SetBody(&body)
	operationsMock.EXPECT().Login(postLoginParameters).Return(nil, errors.New(errorMessage))

	puppetAccess := New()
	puppetAccess.Username = login
	puppetAccess.Password = password
	puppetAccess.Lifetime = lifetime
	puppetAccess.Label = label
	puppetAccess.Client = api

	_, receivedError := puppetAccess.Login()

	assert.EqualError(receivedError, "Unknown error: "+errorMessage)
}

func TestRunLoginSuccess(t *testing.T) {
	assert := assert.New(t)

	login := "username"
	password := "pass"
	lifetime := "10m"
	label := "test_token"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetAccess{
		Operations: operationsMock,
	}

	api.EXPECT().GetClient().Return(client, nil)

	postLoginParameters := operations.NewLoginParamsWithContext(context.Background())
	body := models.LoginRequest{
		Login:    login,
		Password: password,
		Lifetime: lifetime,
		Label:    label,
	}
	postLoginParameters.SetBody(&body)

	responsePayload := operations.NewLoginOK().Payload
	payloadValue := []byte(`{"token":"0IK0epD_I2ejMiXiwct9eXfXYCuiKNVHeWJ9dxomXD1s"}`)

	json.Unmarshal(payloadValue, &responsePayload)

	result := &operations.LoginOK{
		Payload: responsePayload,
	}

	operationsMock.EXPECT().Login(postLoginParameters).Return(result, nil)

	puppetAccess := New()
	puppetAccess.Username = login
	puppetAccess.Password = password
	puppetAccess.Lifetime = lifetime
	puppetAccess.Label = label
	puppetAccess.Client = api

	response, receivedError := puppetAccess.Login()

	assert.NoError(receivedError)
	assert.Equal("0IK0epD_I2ejMiXiwct9eXfXYCuiKNVHeWJ9dxomXD1s", response)
}

func TestRunLoginFailJson(t *testing.T) {
	assert := assert.New(t)

	errorMessage := "The response did not contain a token. Rerun with --debug to see full body"

	login := "username"
	password := "pass"
	lifetime := "10m"
	label := "test_token"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetAccess{
		Operations: operationsMock,
	}

	api.EXPECT().GetClient().Return(client, nil)

	postLoginParameters := operations.NewLoginParamsWithContext(context.Background())
	body := models.LoginRequest{
		Login:    login,
		Password: password,
		Lifetime: lifetime,
		Label:    label,
	}
	postLoginParameters.SetBody(&body)

	responsePayload := operations.NewLoginOK().Payload
	payloadValue := []byte(`{"t":"0IK0epD_I2ejMiXiwct9eXfXYCuiKNVHeWJ9dxomXD1s"}`)

	json.Unmarshal(payloadValue, &responsePayload)

	result := &operations.LoginOK{
		Payload: responsePayload,
	}

	operationsMock.EXPECT().Login(postLoginParameters).Return(result, nil)

	puppetAccess := New()
	puppetAccess.Username = login
	puppetAccess.Password = password
	puppetAccess.Lifetime = lifetime
	puppetAccess.Label = label
	puppetAccess.Client = api

	_, receivedError := puppetAccess.Login()

	assert.Error(receivedError, errorMessage)
}
