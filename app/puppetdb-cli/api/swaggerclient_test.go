package api

import (
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/testdata"
	"github.com/stretchr/testify/assert"
)

func TestGetClientFailsIfNoUrl(t *testing.T) {
	assert := assert.New(t)
	errorMessage := "Invalid scheme for "

	swaggerCfg := SwaggerClientCfg{}
	client := NewClientWithConfig(swaggerCfg)

	_, receivedError := client.GetClient()
	assert.EqualError(receivedError, errorMessage)
}

func TestGetClientSuccessIfHTTP(t *testing.T) {
	assert := assert.New(t)

	swaggerCfg := SwaggerClientCfg{
		URL: "http://random3751.com",
	}
	client := NewClientWithConfig(swaggerCfg)

	_, receivedError := client.GetClient()
	assert.NoError(receivedError)
}

func TestGetClientFailsIfHTTPSNoToken(t *testing.T) {
	assert := assert.New(t)
	errorMessage := "ssl requires a token, please use `puppet access login` to retrieve a token (alternatively use 'cert' and 'key' for whitelist validation)"

	swaggerCfg := SwaggerClientCfg{
		URL: "https://random3751.com",
	}
	client := NewClientWithConfig(swaggerCfg)

	_, receivedError := client.GetClient()
	assert.EqualError(receivedError, errorMessage)
}

func TestGetClientSuccessIfHTTPSWithToken(t *testing.T) {
	assert := assert.New(t)

	swaggerCfg := SwaggerClientCfg{
		URL:   "https://random3751.com",
		Token: filepath.Join(testdata.FixturePath(), "token"),
	}
	client := NewClientWithConfig(swaggerCfg)

	_, receivedError := client.GetClient()
	assert.NoError(receivedError)
}

func TestGetClientSuccessIfHTTPSWithCertAndKey(t *testing.T) {
	assert := assert.New(t)

	swaggerCfg := SwaggerClientCfg{
		Cert: filepath.Join(testdata.FixturePath(), "cert.crt"),
		Key:  filepath.Join(testdata.FixturePath(), "private_key.key"),
		URL:  "http://random3751.com",
	}

	client := NewClientWithConfig(swaggerCfg)

	_, receivedError := client.GetClient()
	assert.NoError(receivedError)
}

func TestGetClientFailsIfCannotParse(t *testing.T) {
	assert := assert.New(t)

	swaggerCfg := SwaggerClientCfg{
		URL: "§¶£¡:random.com",
	}
	client := NewClientWithConfig(swaggerCfg)

	_, receivedError := client.GetClient()
	assert.Error(receivedError)
}
