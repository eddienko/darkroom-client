package storage

import (
	"context"
	"darkroom/pkg/config"
	"fmt"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
)

// Stat fetches metadata for a remote object
func Stat(cfg *config.Config, target string) error {
	accessKey := cfg.S3AccessUser
	secretKey := cfg.S3AccessToken
	if accessKey == "" || secretKey == "" {
		fmt.Println("S3 credentials not found in user info. Please login again.")
		os.Exit(1)
	}

	// Initialize minio client
	minioClient, err := MinioClient(cfg.UserName, cfg.S3AccessToken, true, cfg.UserId)
	if err != nil {
		return fmt.Errorf("failed to create MinIO client: %w", err)
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
	fmt.Printf("Bucket       : %s\n", bucket)
	fmt.Printf("Object       : %s\n", object)
	fmt.Printf("Size         : %s bytes\n", humanizeSize(info.Size))
	fmt.Printf("LastModified : %s\n", info.LastModified.Format("2006-01-02 15:04:05"))
	fmt.Printf("ETag         : %s\n", info.ETag)
	fmt.Printf("ContentType  : %s\n", info.ContentType)

	// Access user-defined metadata
	for k, v := range info.UserMetadata {
		fmt.Printf("%-12s : %s\n", k, v)
	}

	fmt.Println("------------------------------------------------------------------------------------")

	return nil
}
