package cmd

import (
	"darkroom/pkg/storage"

	"github.com/spf13/cobra"
)

var (
	mkBucketRegion string
	mkBucketACL    string
)

var storageMakeBucketCmd = &cobra.Command{
	Use:   "mb <bucket>",
	Short: "Create a new bucket in remote storage",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]

		// cfg, err := config.Load()
		// if err != nil {
		// 	return fmt.Errorf("failed to load config: %w", err)
		// }

		opts := &storage.MakeBucketOptions{
			Region: mkBucketRegion,
			ACL:    mkBucketACL,
		}

		return storage.MakeBucket(cfg, bucket, opts)
	},
}

func init() {
	storageCmd.AddCommand(storageMakeBucketCmd)

	storageMakeBucketCmd.Flags().StringVar(&mkBucketRegion, "region", "", "Region for the bucket (optional)")
	storageMakeBucketCmd.Flags().StringVar(&mkBucketACL, "acl", "private", "ACL for the bucket: private or public-read")
}
