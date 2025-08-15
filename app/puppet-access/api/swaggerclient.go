package api

import (
	"fmt"
	"net/http"
	"net/url"

	tlshelper "github.com/puppetlabs/pe-sdk-go/tls"

	openapihttptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-access/api/client"
	"github.com/puppetlabs/pe-sdk-go/log"
	"github.com/puppetlabs/pe-sdk-go/log/loglevel"
)

// SwaggerClientCfg represent a pe-sdk-go swagger client config
type SwaggerClientCfg struct {
	Login, Password, Lifetime, URL, Label, Cacert, Username string
	UseCNVerification                                       bool
}

// SwaggerClient represents a pe-sdk-go swagger client
type SwaggerClient struct {
	SwaggerClientCfg
}

// NewClient creates a new NewClient
func NewClient(cfg SwaggerClientCfg) Client {
	return &SwaggerClient{
		cfg,
	}
}

// GetClient configures and creates a swagger generated client
func (sc *SwaggerClient) GetClient() (*client.PuppetAccess, error) {
	if sc.URL == "" {
		err := fmt.Errorf("Please provide the service URL. For example, `puppet-access login [username] --service-url https://<HOSTNAME OF PUPPET ENTERPRISE CONSOLE>:4433/rbac-ap`")
		return nil, err
	}

	url, err := url.Parse(sc.URL)
	if err != nil {
		return nil, err
	}

	if url.Scheme != "http" && url.Scheme != "https" {
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
