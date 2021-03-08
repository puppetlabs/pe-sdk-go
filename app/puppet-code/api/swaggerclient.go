package api

import (
	"fmt"
	"net/http"
	"net/url"

	openapihttptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/client"
	"github.com/puppetlabs/pe-sdk-go/log"
	"github.com/puppetlabs/pe-sdk-go/log/loglevel"
)

//SwaggerClient represents a puppec-code swagger client
type SwaggerClient struct {
	cacert, serviceURL string
}

//NewClient creates a new SwaggerClient
func NewClient(cacert, serviceURL string) Client {
	sc := SwaggerClient{
		cacert:     cacert,
		serviceURL: serviceURL,
	}
	return &sc
}

//GetClient configures and creates a swagger generated client
func (sc *SwaggerClient) GetClient() (*client.PuppetCode, error) {
	url, err := url.Parse(sc.serviceURL)
	if err != nil {
		return nil, err
	}
	if url.Scheme != "https" {
		err = fmt.Errorf("Unsupported protocol scheme: %v", url.Scheme)
		return nil, err
	}

	httpclient, err := getHTTPClient(sc.cacert)
	if err != nil {
		return nil, err
	}

	openapitransport := newOpenAPITransport(*url, httpclient)
	openapitransport.SetDebug(log.GetLogLevel() == loglevel.Debug)

	return client.New(openapitransport, strfmt.Default), nil
}

func getHTTPClient(cacert string) (*http.Client, error) {
	tlsClientOptions := openapihttptransport.TLSClientOptions{
		CA: cacert,
	}
	return openapihttptransport.TLSClient(tlsClientOptions)
}

func newOpenAPITransport(url url.URL, httpclient *http.Client) *openapihttptransport.Runtime {
	schemes := []string{url.Scheme}

	return openapihttptransport.NewWithClient(
		fmt.Sprintf("%s:%s", url.Hostname(), url.Port()),
		fmt.Sprintf("%s/v1", url.Path),
		schemes, httpclient)
}
