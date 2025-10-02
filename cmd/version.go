package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "1.0.0"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the darkroom version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("darkroom version %s\n", Version)
		fmt.Printf("Git commit:       %s\n", GitCommit)
		fmt.Printf("Build date:       %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
