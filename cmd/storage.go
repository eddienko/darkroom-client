package cmd

import (
	"github.com/spf13/cobra"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Perform operations on storage",
}

func init() {
	rootCmd.AddCommand(storageCmd)
}
