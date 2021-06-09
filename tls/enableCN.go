package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

// This code is based https://github.com/golang/go/issues/40748#issuecomment-673612108
// and should be used only as a workaround until all puppetdb servers certificates
// are properly generated (and have the SAN fields added)

// EnableCNVerification enables CN verification
func EnableCNVerification(cfg *tls.Config) {
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
