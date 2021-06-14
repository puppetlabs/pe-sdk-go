package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClientSuccess(t *testing.T) {
	assert := assert.New(t)

	swaggerCfg := SwaggerClientCfg{
		URL:      "https://random1234.com",
		Password: "pass",
		Lifetime: "10m",
		Label:    "test_token",
		Login:    "username",
	}

	client := NewClient(swaggerCfg)
	_, receivedError := client.GetClient()
	assert.NoError(receivedError)
}

func TestGetClientFailsIfCannotParse(t *testing.T) {
	assert := assert.New(t)

	swaggerCfg := SwaggerClientCfg{
		URL:      "§¶£¡:random.com",
		Password: "pass",
		Lifetime: "10m",
		Label:    "test_token",
		Login:    "username",
	}
	client := NewClient(swaggerCfg)

	_, receivedError := client.GetClient()
	assert.Error(receivedError)
}
