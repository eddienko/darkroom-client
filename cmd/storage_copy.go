package cmd

import (
	"darkroom/pkg/config"
	"darkroom/pkg/storage"
	"fmt"

	"github.com/spf13/cobra"
)

var recursive bool

// storageCopyCmd represents the `darkroom storage cp` command
var storageCopyCmd = &cobra.Command{
	Use:   "cp <src> <dst>",
	Short: "Copy files between local and remote storage",
	Long: `Copy files and directories between local and remote storage.

Examples:
  # Upload a single file
  darkroom storage cp ./file.txt mybucket/path/

  # Download a single file
  darkroom storage cp mybucket/path/file.txt ./localdir/

  # Upload a directory recursively
  darkroom storage cp -r ./localdir mybucket/path/

  # Download a bucket/prefix recursively
  darkroom storage cp -r mybucket/path ./localdir/
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		src := args[0]
		dst := args[1]

		if err := storage.Copy(cfg, src, dst, recursive); err != nil {
			return fmt.Errorf("copy failed: %w", err)
		}
		return nil
	},
}

func init() {
	storageCopyCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Copy directories recursively")
	storageCmd.AddCommand(storageCopyCmd)
}
