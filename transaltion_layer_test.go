package gcsslcache

import (
	"reflect"
	"testing"

	compute "google.golang.org/api/compute/v1"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgF+zOTDSNRvJXba/q+ud93L2BGQK
vAae7YuHSztfPR1jQNjcqkIAKSg5opeEKYgNEoHv36z5rY1j889wUyUFcvJqmEIE
KagomUVSdltJxmZO0lQW5TuqPbedgB1zV0d/I58yCWPMLJWdRsU0Il3oUwKM5B9A
784Q03nARtzl4Wc3AgMBAAE=
-----END PUBLIC KEY-----
`

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICWgIBAAKBgF+zOTDSNRvJXba/q+ud93L2BGQKvAae7YuHSztfPR1jQNjcqkIA
KSg5opeEKYgNEoHv36z5rY1j889wUyUFcvJqmEIEKagomUVSdltJxmZO0lQW5Tuq
PbedgB1zV0d/I58yCWPMLJWdRsU0Il3oUwKM5B9A784Q03nARtzl4Wc3AgMBAAEC
gYA4Dw5zUM+XZU+mG4Uj0jb/eql7yOX3ouVXlHs3XkS4kEmOP0TkwJ9bVtetldeW
QMIUp2UJOIC3kFNjslKiHx0Dt5UXDcyRGAq6cdj2yt447+m+wOy+hhSp1LCfes80
0Xds6shQ3C+i3dql9v6kKr3ZZxUjOqghlIbeCLDpt0VnYQJBAK3gw+P7DGT+S61i
puxTg2as8/PfHEvCvs89mcV5I5+cCw/HBU1SHd+6hzvrLSO02/RFRE0K8bjSO3Ex
PH4D3DkCQQCM5iCTTtMARpOdSs7W94yClOgkHvItr7g03sRrePg3oYQkQGXM/YzD
O4dniKuvJTmYS/3bQ+JGNegdt/kKeT7vAkA/vnjSKYUfuUJRLCt51BwGFj3RF+gt
thVxsGmhRYnTx8ceX54H/KTLEnzlcJA52OISKRqjC/IWCayVELHWmN+xAkA9p739
d/KxHjEeFUwpmS2tPofOtpP3FfufdxOwi8DiZxUx39QsPY9JJ1V7Ir0t6TYoxKgT
OMNdQd2Ok6CwypmVAkA0PGph3PG2BCDS9SPzO1THG+tinZIuJLwE1nEDsX47FBJY
60rRY/HMtYXs2Kja9AA5mOU1g65Nb0m9Z/noOabZ
-----END RSA PRIVATE KEY-----
`

func TestStableReverse(t *testing.T) {
	cert := &compute.SslCertificate{
		PrivateKey:  privateKey,
		Certificate: publicKey,
	}

	//Let's run through a couple of times.
	d, _ := GCSSLCertificateToAutoCertBytes(cert)
	cert2, _ := AutoCertBytesToGCSSLCertificate(d)
	d, _ = GCSSLCertificateToAutoCertBytes(cert)
	cert2, _ = AutoCertBytesToGCSSLCertificate(d)
	d, _ = GCSSLCertificateToAutoCertBytes(cert)
	cert2, _ = AutoCertBytesToGCSSLCertificate(d)
	d, _ = GCSSLCertificateToAutoCertBytes(cert)
	cert2, _ = AutoCertBytesToGCSSLCertificate(d)

	if !reflect.DeepEqual(cert, cert2) {
		t.Error("Cert1 & Cert2 should be equal but are not, %v, %v", cert, cert2)
	}
}
