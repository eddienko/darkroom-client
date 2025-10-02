package storage

import (
	"context"
	"darkroom/pkg/config"
	"fmt"

	"github.com/minio/minio-go/v7"
)

// MakeBucketOptions defines optional parameters for creating a bucket
type MakeBucketOptions struct {
	Region string
	ACL    string // "private" or "public-read" (MinIO ignores most AWS ACLs)
}

// MakeBucket creates a bucket in the remote S3 storage with optional region/ACL
func MakeBucket(cfg *config.Config, bucketName string, opts *MakeBucketOptions) error {
	client, err := MinioClient(cfg.S3AccessUser, cfg.S3AccessToken, true, cfg.UserId)
	if err != nil {
		return err
	}

	ctx := context.Background()

	bucketOpts := minio.MakeBucketOptions{}
	if opts != nil {
		bucketOpts.Region = opts.Region
		// MinIO mainly supports "private" and "public-read"
		if opts.ACL == "public-read" {
			bucketOpts.ObjectLocking = false // just to differentiate; MinIO ignores AWS ACLs mostly
		}
	}

	err = client.MakeBucket(ctx, bucketName, bucketOpts)
	if err != nil {
		// If bucket already exists, treat as non-fatal
		exists, errBucketExists := client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			fmt.Println("Bucket already exists:", bucketName)
			return nil
		}
		return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
	}

	fmt.Println("Bucket created:", bucketName)
	return nil
}
