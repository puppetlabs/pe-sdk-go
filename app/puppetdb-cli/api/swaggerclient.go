package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	openapihttptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/client"
	"github.com/puppetlabs/pe-sdk-go/log"
	"github.com/puppetlabs/pe-sdk-go/log/loglevel"
)

//SwaggerClientCfg represents a puppetdb-cli swagger client cfg
type SwaggerClientCfg struct {
	Cacert, Cert, Key, URL, Token string
	UseCNVerification             bool
}

//SwaggerClient represents a puppetdb-cli swagger client
type SwaggerClient struct {
	SwaggerClientCfg
}

//NewClientWithConfig creates a new SwaggerClient
func NewClientWithConfig(cfg SwaggerClientCfg) Client {
	return &SwaggerClient{
		cfg,
	}
}

//ArgError represents an argument error
type ArgError struct {
	msg string
}

func (e *ArgError) Error() string {
	return e.msg
}

func supportedScheme(urlScheme string) bool {
	switch urlScheme {
	case "http", "https":
		return true
	default:
		return false
	}
}

func (sc *SwaggerClient) validateSchemeParameters(urlScheme string) error {
	if urlScheme == "https" && (sc.Token == "" && (sc.Cert == "" || sc.Key == "")) {
		return &ArgError{"ssl requires a token, please use `puppet access login` to retrieve a token (alternatively use 'cert' and 'key' for whitelist validation)"}
	}
	return nil
}

//GetClient configures and creates a swagger generated client
func (sc *SwaggerClient) GetClient() (*client.PuppetdbCli, error) {
	url, err := url.Parse(sc.URL)
	if err != nil {
		return nil, err
	}
	if !supportedScheme(url.Scheme) {
		err = fmt.Errorf("Invalid scheme for %v", strings.Title(url.Scheme))
		return nil, err
	}

	if err := sc.validateSchemeParameters(url.Scheme); err != nil {
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
		CA:          sc.Cacert,
		Certificate: sc.Cert,
		Key:         sc.Key,
	}
	cfg, err := openapihttptransport.TLSClientAuth(tlsClientOptions)
	if err != nil {
		return nil, err
	}

	if sc.UseCNVerification { // check server name against CN only
		enableCNVerification(cfg)
	}

	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: cfg,
	}

	return &http.Client{Transport: transport}, nil
}

// This code is based https://github.com/golang/go/issues/40748#issuecomment-673612108
// and should be used only as a workaround until all puppetdb servers certificates
// are properly generated (and have the SAN fields added)
func enableCNVerification(cfg *tls.Config) {
	cfg.InsecureSkipVerify = true
	cfg.VerifyConnection = func(cs tls.ConnectionState) error {
		commonName := cs.PeerCertificates[0].Subject.CommonName
		if commonName != cs.ServerName {
			return fmt.Errorf("invalid certificate name %q, expected %q", commonName, cs.ServerName)
		}
		opts := x509.VerifyOptions{
			Roots:         cfg.RootCAs,
			Intermediates: x509.NewCertPool(),
		}
		for _, cert := range cs.PeerCertificates[1:] {
			opts.Intermediates.AddCert(cert)
		}
		_, err := cs.PeerCertificates[0].Verify(opts)
		return err
	}
}

func newOpenAPITransport(url url.URL, httpclient *http.Client) *openapihttptransport.Runtime {
	schemes := []string{url.Scheme}

	return openapihttptransport.NewWithClient(
		fmt.Sprintf("%s:%s", url.Hostname(), url.Port()),
		fmt.Sprintf("%s", url.Path),
		schemes, httpclient)
}
