package cmd

import (
	"darkroom/pkg/colorfmt"
	"darkroom/pkg/storage"

	"github.com/spf13/cobra"
)

var rmRecursive bool

// storageRmCmd represents `darkroom storage rm`
var storageRmCmd = &cobra.Command{
	Use:   "rm <remote-path>",
	Short: "Remove a file, prefix, or bucket from remote storage",
	Long: `Remove an object, prefix, or bucket from the remote storage.
Examples:
  # Delete a single object
  darkroom storage rm mybucket/file.txt

  # Delete all objects under a prefix
  darkroom storage rm -r mybucket/folder/

  # Delete an entire bucket
  darkroom storage rm -r mybucket
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// cfg, err := config.Load()
		// if err != nil {
		// 	return colorfmt.Error("%v", err)
		// }
		target := args[0]

		if err := storage.Remove(cfg, target, rmRecursive); err != nil {
			return colorfmt.Error("%v", err)
		}
		return nil
	},
}

func init() {
	storageRmCmd.Flags().BoolVarP(&rmRecursive, "recursive", "r", false, "Remove directories or buckets recursively")
	storageCmd.AddCommand(storageRmCmd)
}
