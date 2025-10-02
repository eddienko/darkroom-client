package storage

import (
	"darkroom/pkg/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// NewUserIDClient creates a MinIO client that always sets x-user-id = userID+1000
func MinioClient(accessKey, secretKey string, secure bool, userID int) (*minio.Client, error) {

	endpoint := strings.TrimPrefix(config.BaseURL, "https://") + ":9443"

	idWithOffset := fmt.Sprintf("%d", userID+1000)

	// Clone default transport
	tr := http.DefaultTransport.(*http.Transport).Clone()

	// Wrap transport to inject headers
	rt := roundTripperWithHeader(tr, map[string]string{
		"X-Amz-Meta-User-Id": idWithOffset,
	})

	// Create MinIO client with our custom transport
	return minio.New(endpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure:    secure,
		Transport: rt,
	})
}

// roundTripperWithHeader injects headers into every request
func roundTripperWithHeader(rt http.RoundTripper, headers map[string]string) http.RoundTripper {
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		return rt.RoundTrip(req)
	})
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
