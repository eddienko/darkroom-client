package storage

import (
	"context"
	"darkroom/pkg/config"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Copy supports both upload (local → remote) and download (remote → local).
// If src is a directory or remote prefix, it recursively copies all contents.
func Copy(cfg *config.Config, src, dst string, recursive bool) error {
	if cfg.UserName == "" || cfg.S3AccessToken == "" {
		return fmt.Errorf("S3 credentials not found in user info. Please login again")
	}

	endpoint := strings.TrimPrefix(config.BaseURL, "https://") + ":9443"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.UserName, cfg.S3AccessToken, ""),
		Secure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create MinIO client: %w", err)
	}

	if fileExists(src) {
		// --- Upload single file ---
		return uploadFile(client, src, dst)
	} else if dirExists(src) {
		// --- Upload folder recursively ---
		if !recursive {
			return fmt.Errorf("source is a directory, use --recursive")
		}
		return uploadDir(client, src, dst)
	} else if strings.Contains(src, "/") {
		// --- Download (single or recursive) ---
		if recursive {
			return downloadPrefix(client, src, dst)
		}
		return downloadFile(client, src, dst)
	}

	return fmt.Errorf("invalid source path: %s", src)
}

// --- Upload helpers ---

func uploadFile(client *minio.Client, localPath, remotePath string) error {
	parts := strings.SplitN(remotePath, "/", 2)
	if len(parts) < 1 {
		return fmt.Errorf("remote path must be bucket/object")
	}
	bucket := parts[0]
	objectName := filepath.Base(localPath)
	if len(parts) > 1 {
		if strings.HasSuffix(parts[1], "/") {
			objectName = parts[1] + objectName
		} else {
			objectName = parts[1]
		}
	}

	_, err := client.FPutObject(
		context.Background(),
		bucket,
		objectName,
		localPath,
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	fmt.Printf("✅ Uploaded %s → %s/%s\n", localPath, bucket, objectName)
	return nil
}

func uploadDir(client *minio.Client, localDir, remotePath string) error {
	parts := strings.SplitN(remotePath, "/", 2)
	if len(parts) < 1 {
		return fmt.Errorf("remote path must be bucket/object")
	}
	bucket := parts[0]
	prefix := ""
	if len(parts) > 1 {
		prefix = strings.TrimSuffix(parts[1], "/") + "/"
	}

	err := filepath.WalkDir(localDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		relPath, _ := filepath.Rel(localDir, path)
		objectName := filepath.ToSlash(filepath.Join(prefix, relPath))

		_, err = client.FPutObject(
			context.Background(),
			bucket,
			objectName,
			path,
			minio.PutObjectOptions{ContentType: "application/octet-stream"},
		)
		if err != nil {
			return fmt.Errorf("upload failed: %w", err)
		}
		fmt.Printf("✅ Uploaded %s → %s/%s\n", path, bucket, objectName)
		return nil
	})
	return err
}

// --- Download helpers ---

func downloadFile(client *minio.Client, remotePath, localPath string) error {
	parts := strings.SplitN(remotePath, "/", 2)
	if len(parts) < 2 {
		return fmt.Errorf("remote path must be bucket/object")
	}
	bucket := parts[0]
	objectName := parts[1]

	dst := localPath
	if info, err := os.Stat(localPath); err == nil && info.IsDir() {
		dst = filepath.Join(localPath, filepath.Base(objectName))
	}

	err := client.FGetObject(
		context.Background(),
		bucket,
		objectName,
		dst,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	fmt.Printf("✅ Downloaded %s/%s → %s\n", bucket, objectName, dst)
	return nil
}

func downloadPrefix(client *minio.Client, remotePath, localDir string) error {
	parts := strings.SplitN(remotePath, "/", 2)
	if len(parts) < 2 {
		return fmt.Errorf("remote path must be bucket/prefix")
	}
	bucket := parts[0]
	prefix := strings.TrimSuffix(parts[1], "/") + "/"

	ctx := context.Background()
	for object := range client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}) {
		if object.Err != nil {
			return object.Err
		}

		dst := filepath.Join(localDir, object.Key[len(prefix):])
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return fmt.Errorf("mkdir failed: %w", err)
		}

		err := client.FGetObject(ctx, bucket, object.Key, dst, minio.GetObjectOptions{})
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
		fmt.Printf("✅ Downloaded %s/%s → %s\n", bucket, object.Key, dst)
	}
	return nil
}

// --- Utils ---
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
