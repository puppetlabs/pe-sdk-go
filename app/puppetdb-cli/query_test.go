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

func TestRunQueryFailsIfNoClient(t *testing.T) {
	assert := assert.New(t)
	errorMessage := "No client"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)

	api.EXPECT().GetClient().Return(nil, errors.New(errorMessage))

	puppetDb := New()
	puppetDb.Token = "my token"
	puppetDb.Client = api
	_, receivedError := puppetDb.QueryWithErrorDetails(" ")
	assert.EqualError(receivedError, errorMessage)
}

func TestQueryStatusSucces(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetdbCli{
		Operations: operationsMock,
	}

	api.EXPECT().GetClient().Return(client, nil)

	result := operations.NewGetQueryOK()

	queryParameters := operations.NewGetQueryParamsWithContext(context.Background())
	query := "random.query"

	queryParameters.SetQuery(&query)

	operationsMock.EXPECT().GetQuery(queryParameters, match.XAuthenticationWriter(t, "my token")).Return(result, nil)

	puppetDb := New()
	puppetDb.Token = "my token"
	puppetDb.Client = api
	res, err := puppetDb.QueryWithErrorDetails(query)
	expected := result.Payload
	assert.Equal(expected, res)
	assert.Nil(err)
}

func TestQueryStatusWithError(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetdbCli{
		Operations: operationsMock,
	}

	api.EXPECT().GetClient().Return(client, nil)

	result := operations.NewGetQueryDefault(404)
	result.Payload = &models.Error{
		Msg:     "error message",
		Details: "details",
	}

	queryParameters := operations.NewGetQueryParamsWithContext(context.Background())
	query := "random.query"

	queryParameters.SetQuery(&query)

	operationsMock.EXPECT().GetQuery(queryParameters, match.XAuthenticationWriter(t, "my token")).Return(nil, result)

	puppetCode := New()
	puppetCode.Token = "my token"
	puppetCode.Client = api
	res, err := puppetCode.QueryWithErrorDetails(query)

	assert.Equal(res, "")
	assert.EqualError(err, "[GET /pdb/query/v4][404] getQuery default: error message\n\"details\"")
}
