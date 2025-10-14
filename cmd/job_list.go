package cmd

import (
	"darkroom/pkg/colorfmt"
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

		// err = jobs.ListJobs(cfg)
		err = jobs.ListJobsViaQueryJob(cfg)
		if err != nil {
			return colorfmt.Error("%v", err)
		}
		return nil
	},
}

func init() {
	jobCmd.AddCommand(jobListCmd)
}
