package gcsslcache

import (
	"net/http"
	"time"

	cache "github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
	compute "google.golang.org/api/compute/v1"
)

//As per https://godoc.org/golang.org/x/crypto/acme/autocert
type GoogleCloudSSLCache struct {
	ctx            context.Context
	computeService *compute.Service
	projectID      string
	keyMapper      func(key string) string
	memCache       *cache.Cache
}

// Get returns a certificate data for the specified key.
// If there's no such key, Get returns ErrCacheMiss.
func (gcc *GoogleCloudSSLCache) Get(ctx context.Context, key string) ([]byte, error) {
	// Get the string associated with the key "foo" from the cache
	cachedCert, found := gcc.memCache.Get(key)
	if found {
		return cachedCert.([]byte), nil
	}

	resp, err := gcc.computeService.SslCertificates.Get(gcc.projectID, gcc.keyMapper(key)).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	res, err := GCSSLCertificateToAutoCertBytes(resp)
	if err != nil {
		return nil, err
	}

	gcc.memCache.Set(key, res, cache.DefaultExpiration)
	return res, nil
}

// Put stores the data in the cache under the specified key.
// Underlying implementations may use any data storage format,
// as long as the reverse operation, Get, results in the original data.
func (gcc *GoogleCloudSSLCache) Put(ctx context.Context, key string, data []byte) error {
	rb, err := AutoCertBytesToGCSSLCertificate(data)
	if err != nil {
		return err
	}

	rb.Name = gcc.keyMapper(key)
	_, err = gcc.computeService.SslCertificates.Insert(gcc.projectID, rb).Context(ctx).Do()

	gcc.memCache.Set(key, data, cache.DefaultExpiration)

	return err
}

// Delete removes a certificate data from the cache under the specified key.
// If there's no such key in the cache, Delete returns nil.
func (gcc *GoogleCloudSSLCache) Delete(ctx context.Context, key string) error {
	_, err := gcc.computeService.SslCertificates.Delete(gcc.projectID, gcc.keyMapper(key)).Context(ctx).Do()
	gcc.memCache.Delete(key)
	return err
}

func NewGoogleCloudSSLCache(ctx context.Context, client *http.Client, projectID string, keyMapper func(key string) string) (*GoogleCloudSSLCache, error) {
	c := cache.New(24*time.Hour, 7*24*time.Hour)

	if keyMapper == nil {
		keyMapper = func(key string) string {
			return key
		}
	}
	computeService, err := compute.New(client)
	if err != nil {
		return nil, err
	}

	return &GoogleCloudSSLCache{ctx, computeService, projectID, keyMapper, c}, nil
}
