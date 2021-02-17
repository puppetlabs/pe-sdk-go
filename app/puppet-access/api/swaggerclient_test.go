package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClientSuccess(t *testing.T) {
	assert := assert.New(t)
	login := "username"
	password := "pass"
	lifetime := "10m"
	url := "https://random1234.com"
	label := "test_token"
	cacert := ""

	client := NewClient(login, password, lifetime, url, label, cacert)
	_, receivedError := client.GetClient()
	assert.NoError(receivedError)
}

func TestGetClientFailsIfCannotParse(t *testing.T) {
	assert := assert.New(t)
	login := "username"
	password := "pass"
	lifetime := "10m"
	url := "§¶£¡:random.com"
	label := "test_token"
	cacert := ""

	client := NewClient(login, password, lifetime, url, label, cacert)
	_, receivedError := client.GetClient()
	assert.Error(receivedError)
}
