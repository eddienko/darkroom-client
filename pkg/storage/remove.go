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

// Remove deletes a single object or, if recursive is true, all objects under the prefix
func Remove(cfg *config.Config, target string, recursive bool) error {
	accessKey := cfg.UserName
	secretKey := cfg.S3AccessToken
	if accessKey == "" || secretKey == "" {
		return fmt.Errorf("S3 credentials not found")
	}

	client, err := minio.New(strings.TrimPrefix(config.BaseURL, "https://")+":9443", &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return err
	}

	// Parse bucket and object/prefix
	parts := strings.SplitN(target, "/", 2)
	if len(parts) < 2 {
		return fmt.Errorf("please specify target as bucket/object or bucket/prefix")
	}
	bucket := parts[0]
	prefix := parts[1]

	ctx := context.Background()

	if !recursive {
		// Delete single object
		err := client.RemoveObject(ctx, bucket, prefix, minio.RemoveObjectOptions{})
		if err != nil {
			return err
		}
		fmt.Println("Deleted", target)
		return nil
	}

	// Recursive delete
	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}

	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for obj := range client.ListObjects(ctx, bucket, opts) {
			if obj.Err != nil {
				fmt.Fprintln(os.Stderr, "Error listing object:", obj.Err)
				continue
			}
			objectsCh <- obj
		}
	}()

	for rErr := range client.RemoveObjects(ctx, bucket, objectsCh, minio.RemoveObjectsOptions{}) {
		fmt.Fprintln(os.Stderr, "Error deleting:", rErr)
	}

	fmt.Println("Deleted all objects under", target)
	return nil
}
