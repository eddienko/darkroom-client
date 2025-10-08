package cmd

import (
	"darkroom/pkg/storage"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var expireStr string
var method string

// storagePresignCmd represents `darkroom storage presign`
var storagePresignCmd = &cobra.Command{
	Use:   "presign <remote-path>",
	Short: "Generate a presigned URL for a remote object",
	Long: `Generate a temporary signed URL for downloading or uploading a remote object.

Examples:
  darkroom storage presign mybucket/file.txt --expire 15m
  darkroom storage presign mybucket/file.txt --method PUT --expire 24h
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// cfg, err := config.Load()
		// if err != nil {
		// 	return fmt.Errorf("failed to load config: %w", err)
		// }

		expiry, err := time.ParseDuration(expireStr)
		if err != nil {
			return fmt.Errorf("invalid expiry duration: %w", err)
		}

		url, err := storage.Presign(cfg, args[0], expiry, method)
		if err != nil {
			return err
		}

		fmt.Println(url.String())
		return nil
	},
}

func init() {
	storagePresignCmd.Flags().StringVar(&expireStr, "expire", "1h", "Expiry duration (e.g. 15m, 1h, 24h)")
	storagePresignCmd.Flags().StringVar(&method, "method", "GET", "HTTP method (GET for download, PUT for upload)")
	storageCmd.AddCommand(storagePresignCmd)
}
