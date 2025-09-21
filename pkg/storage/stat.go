package storage

import (
	"context"
	"darkroom/pkg/config"
	"fmt"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Stat fetches metadata for a remote object
func Stat(cfg *config.Config, target string) error {
	accessKey := cfg.UserName
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
		return fmt.Errorf("failed to connect to storage: %w", err)
	}

	parts := strings.SplitN(target, "/", 2)
	if len(parts) < 2 {
		return fmt.Errorf("please specify object path as bucket/object")
	}
	bucket := parts[0]
	object := parts[1]

	info, err := minioClient.StatObject(context.Background(), bucket, object, minio.StatObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to stat %s/%s: %w", bucket, object, err)
	}

	fmt.Println("------------------------------------------------------------------------------------")
	fmt.Printf("Bucket:       %s\n", bucket)
	fmt.Printf("Object:       %s\n", object)
	fmt.Printf("Size:         %d bytes\n", info.Size)
	fmt.Printf("LastModified: %s\n", info.LastModified.Format("2006-01-02 15:04:05"))
	fmt.Printf("ETag:         %s\n", info.ETag)
	fmt.Printf("ContentType:  %s\n", info.ContentType)
	fmt.Println("------------------------------------------------------------------------------------")

	return nil
}
