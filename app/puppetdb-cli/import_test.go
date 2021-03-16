package app

import (
	"context"
	"errors"
	"fmt"
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

func TestRunImportFailsIfNoClient(t *testing.T) {
	assert := assert.New(t)
	filePath := "import.tar.gz"
	errorMessage := "No client"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)

	api.EXPECT().GetClient().Return(nil, errors.New(errorMessage))

	puppetDb := New()
	puppetDb.Token = "my token"
	puppetDb.Client = api

	_, receivedError := puppetDb.PostImportFile(filePath)
	assert.EqualError(receivedError, errorMessage)
}

func TestRunImportFailsIfFileIsAbsent(t *testing.T) {
	assert := assert.New(t)
	filePath := "import.tar.gz"
	errorMessage := fmt.Sprintf("open %s: file does not exist", filePath)

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
	fsmock.EXPECT().Open(filePath).Return(nil, errors.New(errorMessage))

	puppetDb := New()
	puppetDb.Token = "my token"
	puppetDb.Client = api

	_, receivedError := puppetDb.PostImportFile(filePath)
	assert.EqualError(receivedError, errorMessage, "Importing an absent file should fail")
}

func TestRunImportSuccess(t *testing.T) {
	assert := assert.New(t)
	filePath := "import.tar.gz"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetdbCli{
		Operations: operationsMock,
	}
	api.EXPECT().GetClient().Return(client, nil)

	resp := operations.NewPostImportOK().Payload
	result := &operations.PostImportOK{
		Payload: resp,
	}

	fsmock := match.NewMockFs(ctrl)
	appFS = fsmock
	archive := os.NewFile(0, filePath)
	fsmock.EXPECT().Open(filePath).Return(archive, nil)

	postImportParameters := operations.NewPostImportParamsWithContext(context.Background())
	postImportParameters.SetArchive(archive)
	operationsMock.EXPECT().PostImport(postImportParameters, match.XAuthenticationWriter(t, "my token")).Return(result, nil)

	puppetDb := New()
	puppetDb.Token = "my token"
	puppetDb.Client = api

	_, err := puppetDb.PostImportFile(filePath)
	assert.NoError(err)
}

func TestRunImportError(t *testing.T) {
	assert := assert.New(t)
	filePath := "import.tar.gz"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockClient(ctrl)
	operationsMock := mock_operations.NewMockClientService(ctrl)
	client := &client.PuppetdbCli{
		Operations: operationsMock,
	}
	api.EXPECT().GetClient().Return(client, nil)

	result := operations.NewPostImportDefault(404)
	result.Payload = &models.Error{
		Msg:     "error message",
		Details: "details",
	}

	fsmock := match.NewMockFs(ctrl)
	appFS = fsmock

	archive := os.NewFile(0, filePath)
	fsmock.EXPECT().Open(filePath).Return(archive, nil)

	postImportParameters := operations.NewPostImportParamsWithContext(context.Background())
	postImportParameters.SetArchive(archive)
	operationsMock.EXPECT().PostImport(postImportParameters, match.XAuthenticationWriter(t, "my token")).Return(nil, result)

	puppetDb := New()
	puppetDb.Token = "my token"
	puppetDb.Client = api

	res, err := puppetDb.PostImportFile(filePath)
	assert.Equal(false, res)
	assert.EqualError(err, "[POST /pdb/admin/v1/archive][404] postImport default  &{Details:details Kind: Msg:error message}")
}
