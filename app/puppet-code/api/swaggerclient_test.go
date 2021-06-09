package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClientFailsIfNoUrl(t *testing.T) {
	assert := assert.New(t)
	errorMessage := "Unsupported protocol scheme: "

	swaggerCfg := SwaggerClientCfg{}
	client := NewClient(swaggerCfg)

	_, receivedError := client.GetClient()
	assert.EqualError(receivedError, errorMessage)
}

func TestGetClientSuccess(t *testing.T) {
	assert := assert.New(t)

	swaggerCfg := SwaggerClientCfg{
		ServiceURL: "https://random3751.com",
	}
	client := NewClient(swaggerCfg)

	_, receivedError := client.GetClient()

	assert.NoError(receivedError)
}

func TestGetClientFailsIfCannotParse(t *testing.T) {
	assert := assert.New(t)

	swaggerCfg := SwaggerClientCfg{
		ServiceURL: "§¶£¡:random.com",
	}

	client := NewClient(swaggerCfg)
	_, receivedError := client.GetClient()
	assert.Error(receivedError)
}
