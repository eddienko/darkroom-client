package cmd

import (
	"darkroom/pkg/colorfmt"
	"darkroom/pkg/config"
	"darkroom/pkg/jobs"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	completed bool
	failed    bool
	pending   bool
	running   bool
)

var jobListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all submitted UserJobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		selected := 0
		for _, f := range []bool{completed, failed, pending, running} {
			if f {
				selected++
			}
		}
		if selected > 1 {
			return fmt.Errorf("only one of --completed, --failed, --running or --pending can be used at a time")
		}

		var status string
		switch {
		case completed:
			status = "completed"
		case failed:
			status = "failed"
		case pending:
			status = "pending"
		case running:
			status = "running"
		default:
			status = "" // or "all" if you want a default
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// err = jobs.ListJobs(cfg)

		err = jobs.ListJobsViaQueryJob(cfg, status)
		if err != nil {
			return colorfmt.Error("%v", err)
		}
		return nil
	},
}

func init() {
	jobCmd.AddCommand(jobListCmd)

	jobListCmd.Flags().BoolVarP(&completed, "completed", "c", false, "show completed jobs")
	jobListCmd.Flags().BoolVarP(&failed, "failed", "f", false, "show failed jobs")
	jobListCmd.Flags().BoolVarP(&pending, "pending", "p", false, "show pending jobs")
	jobListCmd.Flags().BoolVarP(&running, "running", "r", false, "show running jobs")

}
