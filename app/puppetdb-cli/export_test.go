package app

import (
	"context"
	"errors"
	"io"
	"os"
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

func TestRunExportFailsIfNoClient(t *testing.T) {
	assert := assert.New(t)
	filePath := "export.tar.gz"
	errorMessage := "No client"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)

	api.EXPECT().GetClient().Return(nil, errors.New(errorMessage))

	puppetDb := New()
	puppetDb.Token = "my token"
	puppetDb.Client = api
	_, receivedError := puppetDb.GetExportFile(filePath, "none")
	assert.EqualError(receivedError, errorMessage)
}

func TestRunExportFailsIfFileCreationFails(t *testing.T) {
	assert := assert.New(t)
	filePath := "export.tar.gz"
	errorMessage := "operation not permitted"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetdbCli{
		Operations: operationsMock,
	}

	api.EXPECT().GetClient().Return(client, nil)

	fsmock := match.NewMockFs(ctrl)
	appFS = fsmock
	fsmock.EXPECT().Create(filePath).Return(nil, errors.New(errorMessage))

	puppetDb := New()
	puppetDb.Token = "my token"
	puppetDb.Client = api
	_, receivedError := puppetDb.GetExportFile(filePath, "none")
	assert.EqualError(receivedError, errorMessage, "Archive file creation should fail")
}

func TestRunExportSucces(t *testing.T) {
	assert := assert.New(t)
	filePath := "export.tar.gz"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetdbCli{
		Operations: operationsMock,
	}
	api.EXPECT().GetClient().Return(client, nil)

	var mockPayload io.Writer
	result := &operations.GetExportOK{
		Payload: mockPayload,
	}

	fsmock := match.NewMockFs(ctrl)
	appFS = fsmock
	archive := os.NewFile(0, filePath)
	fsmock.EXPECT().Create(filePath).Return(archive, nil)

	getExportParameters := operations.NewGetExportParamsWithContext(context.Background())
	anon := "none"
	getExportParameters.SetAnonymizationProfile(&anon)
	operationsMock.EXPECT().GetExport(getExportParameters, match.XAuthenticationWriter(t, "my token"), archive).Return(result, nil)

	puppetDb := New()
	puppetDb.Token = "my token"
	puppetDb.Client = api

	res, err := puppetDb.GetExportFile(filePath, "none")
	assert.Equal("<nil>", res)
	assert.NoError(err)
}

func TestRunExportError(t *testing.T) {
	assert := assert.New(t)
	filePath := "export.tar.gz"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetdbCli{
		Operations: operationsMock,
	}
	api.EXPECT().GetClient().Return(client, nil)

	result := operations.NewGetExportDefault(404)
	result.Payload = &models.Error{
		Msg:     "error message",
		Details: "details",
	}

	fsmock := match.NewMockFs(ctrl)
	appFS = fsmock
	archive := os.NewFile(0, filePath)
	fsmock.EXPECT().Create(filePath).Return(archive, nil)
	fsmock.EXPECT().Remove(filePath).Return(nil)

	getExportParameters := operations.NewGetExportParamsWithContext(context.Background())
	anon := "none"
	getExportParameters.SetAnonymizationProfile(&anon)
	operationsMock.EXPECT().GetExport(getExportParameters, match.XAuthenticationWriter(t, "my token"), gomock.Any()).Return(nil, result)

	puppetDb := New()
	puppetDb.Token = "my token"
	puppetDb.Client = api
	res, err := puppetDb.GetExportFile(filePath, "none")

	assert.Equal("", res)
	assert.EqualError(err, "[GET /pdb/admin/v1/archive][404] getExport default  &{Details:details Kind: Msg:error message}")
}
