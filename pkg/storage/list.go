package storage

import (
	"context"
	"darkroom/pkg/config"
	"fmt"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type entry struct {
	Name         string
	IsDir        bool
	Size         int64
	LastModified time.Time
}

// convert size in bytes into human-readable format
func humanizeSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	kb := float64(size) / 1024.0
	if kb < 1024 {
		return fmt.Sprintf("%.1f KB", kb)
	}
	mb := kb / 1024.0
	if mb < 1024 {
		return fmt.Sprintf("%.1f MB", mb)
	}
	gb := mb / 1024.0
	if gb < 1024 {
		return fmt.Sprintf("%.1f GB", gb)
	}
	tb := gb / 1024.0
	return fmt.Sprintf("%.1f TB", tb)
}

// List prints either all buckets or the objects in a given bucket/prefix.
func List(cfg *config.Config, arg string) error {
	// Validate credentials
	if cfg.UserName == "" || cfg.S3AccessToken == "" {
		return fmt.Errorf("S3 credentials not found in user info. Please login again")
	}

	// Build endpoint
	endpoint := ""
	if endpoint == "" {
		endpoint = strings.TrimPrefix(config.BaseURL, "https://") + ":9443"
	}

	// Init MinIO client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.UserName, cfg.S3AccessToken, ""),
		Secure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create MinIO client: %w", err)
	}

	ctx := context.Background()

	if arg == "" || arg == "/" {
		// --- List buckets ---
		buckets, err := client.ListBuckets(ctx)
		if err != nil {
			return fmt.Errorf("failed to list buckets: %w", err)
		}

		fmt.Println("------------------------------------------------------------------------------------")
		fmt.Printf("%-40s %-20s\n", "BUCKET NAME", "CREATED")
		fmt.Println("------------------------------------------------------------------------------------")
		for _, bucket := range buckets {
			fmt.Printf("%-40s %-20s\n",
				bucket.Name,
				bucket.CreationDate.Format(time.RFC3339),
			)
		}
		return nil
	}

	// --- List objects in bucket/prefix ---
	parts := strings.SplitN(arg, "/", 2)
	bucket := parts[0]
	prefix := ""
	if len(parts) > 1 {
		prefix = parts[1]
		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}
	}

	objectCh := client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: false,
	})

	var entries []entry

	for object := range objectCh {
		if object.Err != nil {
			return fmt.Errorf("failed to list objects: %w", object.Err)
		}

		if strings.HasSuffix(object.Key, "/") {
			// Directory
			entries = append(entries, entry{
				Name:  path.Base(strings.TrimRight(object.Key, "/")) + "/",
				IsDir: true,
			})
		} else {
			// File
			entries = append(entries, entry{
				Name:         path.Base(object.Key),
				IsDir:        false,
				Size:         object.Size,
				LastModified: object.LastModified,
			})
		}
	}

	// Sort: directories first, then files, both alphabetically
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir && !entries[j].IsDir {
			return true
		}
		if !entries[i].IsDir && entries[j].IsDir {
			return false
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	// Print results
	fmt.Println("------------------------------------------------------------------------------------")
	fmt.Printf("%-60s %-12s %-20s\n", "NAME", "SIZE", "LAST MODIFIED")
	fmt.Println("------------------------------------------------------------------------------------")

	for _, e := range entries {
		if e.IsDir {
			fmt.Printf("%-60s %-12s %-20s\n", e.Name, "-", "-")
		} else {
			fmt.Printf("%-60s %-12s %-20s\n", e.Name, humanizeSize(e.Size), e.LastModified.Format(time.RFC3339))
		}
	}

	return nil
}
