package tlscfg

import (
	"crypto/x509"
	_ "embed"
)

//go:embed cacerts
var caCerts []byte

// ClientCertPool returns *x509.CertPool with certificates
func ClientCertPool() *x509.CertPool {
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCerts)

	return pool
}
