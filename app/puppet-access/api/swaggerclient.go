package api

import (
	"fmt"
	"net/http"
	"net/url"

	openapihttptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/puppetlabs/pe-sdk-go/app/puppet-access/api/client"
	"github.com/puppetlabs/pe-sdk-go/log"
	"github.com/puppetlabs/pe-sdk-go/log/loglevel"
)

//SwaggerClient represents a pe-sdk-go swagger client
type SwaggerClient struct {
	login, password, lifetime, url, label, cacert string
}

//NewClient creates a new NewClient
func NewClient(login, password, lifetime, url, label, cacert string) Client {
	sc := SwaggerClient{
		login:    login,
		password: password,
		lifetime: lifetime,
		url:      url,
		label:    label,
		cacert:   cacert,
	}
	return &sc
}

//GetClient configures and creates a swagger generated client
func (sc *SwaggerClient) GetClient() (*client.PuppetAccess, error) {

	url, err := url.Parse(sc.url)
	if err != nil {
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
	cfg, err := openapihttptransport.TLSClientAuth(tlsClientOptions)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: cfg,
	}

	return &http.Client{Transport: transport}, nil
}

func newOpenAPITransport(url url.URL, httpclient *http.Client) *openapihttptransport.Runtime {
	schemes := []string{url.Scheme}

	return openapihttptransport.NewWithClient(
		fmt.Sprintf("%s:%s", url.Hostname(), url.Port()),
		fmt.Sprintf("%s", url.Path),
		schemes, httpclient)
}
