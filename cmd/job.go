package cmd

import (
	"github.com/spf13/cobra"
)

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Manage jobs on the cluster",
}

func init() {
	rootCmd.AddCommand(jobCmd)
}
