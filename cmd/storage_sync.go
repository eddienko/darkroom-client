package cmd

import (
	"darkroom/pkg/storage"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	syncDelete    bool
	syncDirection string
	syncChecksums bool
)

var storageSyncCmd = &cobra.Command{
	Use:   "sync <localdir> <remote>",
	Short: "Synchronize local files with remote object storage",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		localPath := args[0]
		remotePath := args[1]

		// cfg, err := config.Load()
		// if err != nil {
		// 	return fmt.Errorf("failed to load config: %w", err)
		// }

		opts := storage.SyncOptions{
			Delete:    syncDelete,
			Direction: syncDirection,
			Checksums: syncChecksums,
		}

		if err := storage.Sync(cfg, localPath, remotePath, opts); err != nil {
			return fmt.Errorf("sync failed: %w", err)
		}
		return nil
	},
}

func init() {
	storageCmd.AddCommand(storageSyncCmd)

	storageSyncCmd.Flags().BoolVar(&syncDelete, "delete", false, "remove files in destination that are not in source")
	storageSyncCmd.Flags().StringVar(&syncDirection, "direction", "up", "sync direction: 'up' (local → remote) or 'down' (remote → local)")
	storageSyncCmd.Flags().BoolVar(&syncChecksums, "checksums", false, "compare checksums to skip unchanged files")
}
