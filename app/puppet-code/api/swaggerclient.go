package api

import (
	"fmt"
	"net/http"
	"net/url"

	tlshelper "github.com/puppetlabs/pe-sdk-go/tls"

	openapihttptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-code/api/client"
	"github.com/puppetlabs/pe-sdk-go/log"
	"github.com/puppetlabs/pe-sdk-go/log/loglevel"
)

// SwaggerClientCfg represent a pe-sdk-go swagger client config
type SwaggerClientCfg struct {
	Cacert, ServiceURL, Token string
	UseCNVerification         bool
}

// SwaggerClient represents a puppec-code swagger client
type SwaggerClient struct {
	SwaggerClientCfg
}

// NewClient creates a new SwaggerClient
func NewClient(cfg SwaggerClientCfg) Client {
	return &SwaggerClient{
		cfg,
	}
}

// GetClient configures and creates a swagger generated client
func (sc *SwaggerClient) GetClient() (*client.PuppetCode, error) {
	url, err := url.Parse(sc.ServiceURL)
	if err != nil {
		return nil, err
	}
	if url.Scheme != "https" {
		err = fmt.Errorf("Unsupported protocol scheme: %v", url.Scheme)
		return nil, err
	}

	httpclient, err := sc.getHTTPClient()
	if err != nil {
		return nil, err
	}

	openapitransport := newOpenAPITransport(*url, httpclient)
	openapitransport.SetDebug(log.GetLogLevel() == loglevel.Debug)

	return client.New(openapitransport, strfmt.Default), nil
}

func (sc *SwaggerClient) getHTTPClient() (*http.Client, error) {
	tlsClientOptions := openapihttptransport.TLSClientOptions{
		CA: sc.Cacert,
	}

	transport, err := tlshelper.GetHTTPTransport(tlsClientOptions, sc.UseCNVerification)
	if err != nil {
		return nil, err
	}

	return &http.Client{Transport: transport}, nil
}

func newOpenAPITransport(url url.URL, httpclient *http.Client) *openapihttptransport.Runtime {
	schemes := []string{url.Scheme}

	return openapihttptransport.NewWithClient(
		fmt.Sprintf("%s:%s", url.Hostname(), url.Port()),
		fmt.Sprintf("%s/v1", url.Path),
		schemes, httpclient)
}
