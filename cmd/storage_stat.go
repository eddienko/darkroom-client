package cmd

import (
	"darkroom/pkg/storage"

	"github.com/spf13/cobra"
)

// storageStatCmd represents `darkroom storage stat`
var storageStatCmd = &cobra.Command{
	Use:   "stat <remote-path>",
	Short: "Show metadata for a remote object",
	Long: `Display detailed metadata for a remote object in storage.

Examples:
  darkroom storage stat mybucket/file.txt
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// cfg, err := config.Load()
		// if err != nil {
		// 	return fmt.Errorf("failed to load config: %w", err)
		// }
		target := args[0]
		return storage.Stat(cfg, target)
	},
}

func init() {
	storageCmd.AddCommand(storageStatCmd)
}
