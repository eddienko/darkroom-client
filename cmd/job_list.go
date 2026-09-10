package cmd

import (
	"darkroom/pkg/colorfmt"
	"darkroom/pkg/config"
	"darkroom/pkg/jobs"
	"fmt"

	"github.com/spf13/cobra"
)

var jobStatus string

var validJobStatuses = map[string]bool{
	"pending":   true,
	"running":   true,
	"completed": true,
	"failed":    true,
	"cancelled": true,
}

var jobListCmd = &cobra.Command{
	Use:   "list",
	Short: "List submitted jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if jobStatus != "" && !validJobStatuses[jobStatus] {
			return fmt.Errorf(
				"invalid status %q: must be one of pending, running, completed, failed, cancelled",
				jobStatus,
			)
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		resp, err := jobs.ListJobs(cfg, jobStatus)
		if err != nil {
			return colorfmt.Error("%v", err)
		}

		jobs.PrintJobs(resp)

		return nil
	},
}

func init() {
	jobCmd.AddCommand(jobListCmd)

	jobListCmd.Flags().StringVar(
		&jobStatus,
		"status",
		"",
		"filter jobs by status (pending, running, completed, failed, cancelled)",
	)
}
