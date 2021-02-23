package app

import (
	"fmt"
)

type puppetConfig struct {
	ServiceURL string `json:"service-url"`
	Cacert     string `json:"cacert"`
	TokenFile  string `json:"token-file"`
	Token      string `json:"token"`
}

//GetConfig will return the configuration used
func (puppetCode *PuppetCode) GetConfig() string {

	token, err := puppetCode.Token.Read()
	if err != nil {
		token = "NotFound"
	}
	return fmt.Sprintf("service-url: \"%s\"\ncacert: \"%s\"\ntoken-file: \"%s\"\ntoken: \"%s\"\n",
		puppetCode.ServiceURL, puppetCode.Cacert, puppetCode.TokenFile, token)
}
