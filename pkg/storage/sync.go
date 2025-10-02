package storage

import (
	"context"
	"crypto/md5"
	"darkroom/pkg/config"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type SyncOptions struct {
	Delete    bool
	Direction string // "up" or "down"
	Checksums bool   // whether to compare checksums
}

func Sync(cfg *config.Config, localPath, remotePath string, opts SyncOptions) error {
	accessKey := cfg.S3AccessUser
	secretKey := cfg.S3AccessToken
	if accessKey == "" || secretKey == "" {
		return fmt.Errorf("S3 credentials missing")
	}

	endpoint := strings.TrimPrefix(config.BaseURL, "https://") + ":9443"
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return err
	}

	if opts.Direction == "down" {
		return syncDown(client, localPath, remotePath, opts)
	}
	return syncUp(client, localPath, remotePath, opts)
}

// --- upload local → remote ---
func syncUp(client *minio.Client, localPath, remotePath string, opts SyncOptions) error {
	parts := strings.SplitN(strings.TrimPrefix(remotePath, "s3://"), "/", 2)
	bucket := parts[0]
	prefix := ""
	if len(parts) > 1 {
		prefix = parts[1]
	}

	return filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(localPath, path)
		remoteKey := filepath.ToSlash(filepath.Join(prefix, rel))

		// check if object exists
		stat, err := client.StatObject(context.Background(), bucket, remoteKey, minio.StatObjectOptions{})
		if err == nil {
			if opts.Checksums {
				// compare MD5 checksums
				localHash, err := md5File(path)
				if err != nil {
					return err
				}
				if strings.Trim(stat.ETag, `"`) == localHash {
					// skip unchanged
					return nil
				}
			}
		}

		// upload
		_, err = client.FPutObject(context.Background(), bucket, remoteKey, path, minio.PutObjectOptions{})
		if err != nil {
			return err
		}
		fmt.Println("Uploaded:", remoteKey)
		return nil
	})
}

// --- download remote → local ---
func syncDown(client *minio.Client, localPath, remotePath string, opts SyncOptions) error {
	parts := strings.SplitN(strings.TrimPrefix(remotePath, "s3://"), "/", 2)
	bucket := parts[0]
	prefix := ""
	if len(parts) > 1 {
		prefix = parts[1]
	}

	for object := range client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}) {
		if object.Err != nil {
			return object.Err
		}
		localFile := filepath.Join(localPath, strings.TrimPrefix(object.Key, prefix))

		// skip if exists and checksum matches
		if _, err := os.Stat(localFile); err == nil && opts.Checksums {
			localHash, err := md5File(localFile)
			if err == nil && strings.Trim(object.ETag, `"`) == localHash {
				continue // unchanged
			}
		}

		if err := os.MkdirAll(filepath.Dir(localFile), 0755); err != nil {
			return err
		}
		err := client.FGetObject(context.Background(), bucket, object.Key, localFile, minio.GetObjectOptions{})
		if err != nil {
			return err
		}
		fmt.Println("Downloaded:", localFile)
	}
	return nil
}

// --- helper: compute md5 of a file ---
func md5File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
