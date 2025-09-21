package cmd

import (
	"darkroom/pkg/config"
	"darkroom/pkg/jobs"
	"fmt"

	"github.com/spf13/cobra"
)

var jobListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all submitted UserJobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		return jobs.ListJobs(cfg)
	},
}

func init() {
	jobCmd.AddCommand(jobListCmd)
}
