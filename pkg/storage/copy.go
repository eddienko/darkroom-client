package storage

import (
	"context"
	"crypto/md5"
	"darkroom/pkg/config"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"io"

	"github.com/minio/minio-go/v7"

	"github.com/schollz/progressbar/v3"
)

// Copy supports both upload (local → remote) and download (remote → local).
// If src is a directory or remote prefix, it recursively copies all contents.
func Copy(cfg *config.Config, src, dst string, recursive bool) error {
	if cfg.UserName == "" || cfg.S3AccessToken == "" {
		return fmt.Errorf("S3 credentials not found in user info. Please login again")
	}

	client, err := MinioClient(cfg.UserName, cfg.S3AccessToken, true, cfg.UserId)
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

// uploadFile uploads a file to MinIO with a schollz/progressbar that works with multi-part uploads.
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

	// Open file
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	// Get file info
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("cannot stat file: %w", err)
	}
	size := stat.Size()
	md5sum, _ := computeFileMD5(localPath)

	// Create progress bar with nicer config for large files
	bar := progressbar.NewOptions64(
		size,
		progressbar.OptionSetDescription("Uploading"),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(65), // smooth updates
		progressbar.OptionShowCount(),  // show counters
		// progressbar.OptionShowIts(),            // show speed
		progressbar.OptionSetPredictTime(true), // ETA
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "▶",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	// Wrap file with TeeReader so every read updates the bar
	reader := io.TeeReader(file, bar)

	// Upload file (MinIO Go handles multipart under the hood for large files)
	objInfo, err := client.PutObject(
		context.Background(),
		bucket,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
			UserMetadata: map[string]string{
				"X-Amz-Meta-Md5sum": md5sum,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	// Verify MD5 for small files only
	if size <= 5*1024*1024 { // 5 MB
		remoteETag := objInfo.ETag
		remoteETag = strings.Trim(remoteETag, `"`)
		if md5sum != remoteETag {
			return fmt.Errorf("MD5 mismatch for %s: local=%s remote=%s", objectName, md5sum, remoteETag)
		}
	}

	fmt.Printf("\n✅ Uploaded %s → %s/%s\n", localPath, bucket, objectName)
	return nil
}

// uploadDir uploads a whole directory recursively with a global progress bar.
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

	// 1. Pre-scan directory to get total size
	var totalSize int64
	err := filepath.WalkDir(localDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, statErr := os.Stat(path)
			if statErr != nil {
				return statErr
			}
			totalSize += info.Size()
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	// 2. Create a single global progress bar
	bar := progressbar.NewOptions64(
		totalSize,
		progressbar.OptionSetDescription("Uploading"),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(65),
		progressbar.OptionShowCount(),
		// progressbar.OptionShowIts(),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "▶",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	// 3. Walk directory again and upload with progress tracking
	err = filepath.WalkDir(localDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(localDir, path)
		objectName := filepath.ToSlash(filepath.Join(prefix, relPath))

		// open file
		file, openErr := os.Open(path)
		if openErr != nil {
			return openErr
		}
		defer file.Close()

		// get file size
		info, statErr := file.Stat()
		if statErr != nil {
			return statErr
		}
		size := info.Size()
		md5sum, _ := computeFileMD5(path)

		// wrap with progress bar
		reader := io.TeeReader(file, bar)

		objInfo, putErr := client.PutObject(
			context.Background(),
			bucket,
			objectName,
			reader,
			size,
			minio.PutObjectOptions{
				ContentType: "application/octet-stream",
				UserMetadata: map[string]string{
					"X-Amz-Meta-Md5sum": md5sum,
				},
			},
		)
		if putErr != nil {
			return fmt.Errorf("upload failed: %w", putErr)
		}

		// Verify MD5 for small files only
		if size <= 5*1024*1024 { // 5 MB
			remoteETag := objInfo.ETag
			remoteETag = strings.Trim(remoteETag, `"`)
			if md5sum != remoteETag {
				return fmt.Errorf("MD5 mismatch for %s: local=%s remote=%s", objectName, md5sum, remoteETag)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	fmt.Printf("\n✅ Finished uploading directory %s → %s/%s\n", localDir, bucket, prefix)
	return nil
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

func computeFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("cannot read file: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
