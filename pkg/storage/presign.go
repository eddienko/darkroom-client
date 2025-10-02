package storage

import (
	"context"
	"darkroom/pkg/config"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Presign generates a presigned URL for an object, supporting GET or PUT
func Presign(cfg *config.Config, target string, expiry time.Duration, method string) (*url.URL, error) {
	accessKey := cfg.S3AccessUser
	secretKey := cfg.S3AccessToken
	if accessKey == "" || secretKey == "" {
		fmt.Println("S3 credentials not found in user info. Please login again.")
		os.Exit(1)
	}

	// Initialize minio client
	minioClient, err := minio.New(strings.TrimPrefix(config.BaseURL, "https://")+":9443", &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to storage: %w", err)
	}

	parts := strings.SplitN(target, "/", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("please specify object path as bucket/object")
	}
	bucket := parts[0]
	object := parts[1]

	reqParams := make(url.Values)

	var presignedURL *url.URL
	switch strings.ToUpper(method) {
	case "GET":
		presignedURL, err = minioClient.PresignedGetObject(context.Background(), bucket, object, expiry, reqParams)
	case "PUT":
		presignedURL, err = minioClient.PresignedPutObject(context.Background(), bucket, object, expiry)
	default:
		return nil, fmt.Errorf("unsupported method: %s (use GET or PUT)", method)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to presign object %s/%s: %w", bucket, object, err)
	}

	return presignedURL, nil
}
