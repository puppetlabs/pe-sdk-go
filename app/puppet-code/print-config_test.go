package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericFailureCheck(t *testing.T) {
	assert := assert.New(t)
	expectedConfig := "service-url: \"\"\ncacert: \"\"\ntoken-file: \"\"\ntoken: \"NotFound\"\n"
	assert.Equal(expectedConfig, New().GetConfig(), "Empty config is not matching")
}
