package cmd

import (
	"darkroom/pkg/config"
	"darkroom/pkg/jobs"
	"fmt"

	"github.com/spf13/cobra"
)

var jobStatusCmd = &cobra.Command{
	Use:   "status <jobName>",
	Short: "Show detailed status of a UserJob",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		return jobs.JobStatus(cfg, args[0])
	},
}

func init() {
	jobCmd.AddCommand(jobStatusCmd)
}
