# autocert-gc-cache - NON WORKING
A cache implementation for acme/autocert using google cloud compute/sslCertificates as a certificate storage.

DOES NOT WORK (ATM). The google API implementation for ssl certs does not allow for updating existing certs meaning the processes needs to look like 
1) Create new cert 
2) Update all sources using old cert to use new 
3) Delete old cert
#2 means a significantly larger scope then expected for this project. Leaving here for other parties who might be considering the same approach.


# NOTE
Highly exprimental & limited test coverage

# Installation 
`go get github.com/zatte/autocert-gc-cache`


# Usage - Example
```golang
package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

func main() {

	ctx := context.Background()
	projectId := "google-cloud-project-id"
	googleClient, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}
	GCSSLCache, err := NewGoogleCloudSSLCache(ctx, googleClient, projectId, nil)
	if err != nil {
		log.Fatal(err)
	}

	//From Autocert : https://godoc.org/golang.org/x/crypto/acme/autocert#pkg-files
	m := &autocert.Manager{
		Cache:      GCSSLCache,
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("example.org"),
	}
	go http.ListenAndServe(":http", m.HTTPHandler(nil))
	s := &http.Server{
		Addr:      ":https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
	}
	s.ListenAndServeTLS("", "")
}

```
