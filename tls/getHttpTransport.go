package tls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"strings"

	openapihttptransport "github.com/go-openapi/runtime/client"
	"github.com/puppetlabs/pe-sdk-go/log"
)

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

// GetHTTPTransport returns a http transport handling certificates without SAN
func GetHTTPTransport(tlsClientOptions openapihttptransport.TLSClientOptions, useCNVerification bool) (*http.Transport, error) {
	cfg, err := openapihttptransport.TLSClientAuth(tlsClientOptions)
	if err != nil {
		return nil, err
	}

	if useCNVerification { // check server name against CN only
		enableCNVerification(cfg)
	}

	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: cfg,
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {

			tlsDialer := tls.Dialer{
				Config: cfg,
			}

			conn, err := tlsDialer.DialContext(ctx, network, addr)
			if err != nil {
				if err != nil && strings.Contains(err.Error(), "temporarily enable Common Name matching with") {
					log.Warn("x509: certificate relies on legacy Common Name field, use SAN, or the -k/--use-cn-verification flag to avoid this warning")
					enableCNVerification(cfg)
					return tlsDialer.DialContext(ctx, network, addr)
				}
			}
			return conn, err
		},
	}

	return transport, nil
}
