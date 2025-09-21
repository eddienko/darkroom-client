package storage

import (
	"context"
	"darkroom/pkg/config"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Remove deletes an object or prefix from remote storage
func Remove(cfg *config.Config, target string, recursive bool) error {
	accessKey := cfg.UserName
	secretKey := cfg.S3AccessToken
	if accessKey == "" || secretKey == "" {
		fmt.Println("S3 credentials not found in user info. Please login again.")
		os.Exit(1)
	}

	// Initialize minio client
	endpoint := strings.TrimPrefix(config.BaseURL, "https://") + ":9443"
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to storage: %w", err)
	}

	parts := strings.SplitN(target, "/", 2)
	bucket := parts[0]
	var prefix string
	if len(parts) > 1 {
		prefix = parts[1]
	}

	if prefix == "" {
		// Trying to delete whole bucket
		if recursive {
			if err := minioClient.RemoveBucket(context.Background(), bucket); err != nil {
				return fmt.Errorf("failed to delete bucket %s: %w", bucket, err)
			}
			fmt.Printf("Bucket %s deleted.\n", bucket)
			return nil
		}
		return fmt.Errorf("refusing to delete entire bucket without --recursive")
	}

	// If recursive, delete all objects under prefix
	if recursive {
		opts := minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		}

		ch := minioClient.ListObjects(context.Background(), bucket, opts)

		for obj := range ch {
			if obj.Err != nil {
				log.Println("list error:", obj.Err)
				continue
			}
			err := minioClient.RemoveObject(context.Background(), bucket, obj.Key, minio.RemoveObjectOptions{})
			if err != nil {
				log.Printf("failed to delete %s: %v\n", obj.Key, err)
			} else {
				fmt.Printf("Deleted %s/%s\n", bucket, obj.Key)
			}
		}
		return nil
	}

	// Non-recursive → delete single object
	if err := minioClient.RemoveObject(context.Background(), bucket, prefix, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("failed to delete %s/%s: %w", bucket, prefix, err)
	}
	fmt.Printf("Deleted %s/%s\n", bucket, prefix)
	return nil
}
