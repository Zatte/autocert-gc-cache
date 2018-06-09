package gcsslcache

import (
	"encoding/pem"
	"strings"

	compute "google.golang.org/api/compute/v1"
)

// GCSSLCertificateToAutoCertBytes takes the Google Certificate response and turns it into
// an pem-encoded byte slice that is interpretable by autocert. This functionality depends
// heavily on the inner working of autocert and can change / break at any time
func GCSSLCertificateToAutoCertBytes(cert *compute.SslCertificate) ([]byte, error) {
	//since the input is multiple PEM encoded fields we can just concatenate, however order is important!
	return []byte(strings.Join([]string{
		cert.PrivateKey,
		cert.Certificate,
	}, "\n")), nil
}

// AutoCertBytesToGCSSLCertificate takes an byte stream from autocert and turn the
// data into a Google Certificate. This functionality depends
// heavily on the inner working of autocert and can change / break at any time
func AutoCertBytesToGCSSLCertificate(data []byte) (*compute.SslCertificate, error) {
	priv, rest := pem.Decode(data)
	pub, rest := pem.Decode(rest)
	rb := &compute.SslCertificate{
		PrivateKey:  string(pem.EncodeToMemory(priv)),
		Certificate: string(pem.EncodeToMemory(pub)),
	}
	return rb, nil
}
