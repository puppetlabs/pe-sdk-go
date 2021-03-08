package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClientFailsIfNoUrl(t *testing.T) {
	assert := assert.New(t)
	errorMessage := "Unsupported protocol scheme: "

	cacert := ""
	serviceURL := ""
	client := NewClient(cacert, serviceURL)
	_, receivedError := client.GetClient()
	assert.EqualError(receivedError, errorMessage)
}

func TestGetClientSuccess(t *testing.T) {
	assert := assert.New(t)

	cacert := ""
	serviceURL := "https://random3751.com"
	client := NewClient(cacert, serviceURL)

	_, receivedError := client.GetClient()

	assert.NoError(receivedError)
}

func TestGetClientFailsIfCannotParse(t *testing.T) {
	assert := assert.New(t)

	cacert := ""
	serviceURL := "§¶£¡:random.com"
	client := NewClient(cacert, serviceURL)
	_, receivedError := client.GetClient()
	assert.Error(receivedError)
}
