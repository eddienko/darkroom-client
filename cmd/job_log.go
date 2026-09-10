package cmd

import (
	"darkroom/pkg/colorfmt"
	"darkroom/pkg/config"
	"darkroom/pkg/jobs"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var jobLogsCmd = &cobra.Command{
	Use:   "logs JOB_ID",
	Short: "Show logs for a job",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		jobID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid job ID %q", args[0])
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		logs, err := jobs.GetJobLogs(cfg, jobID)
		if err != nil {
			return colorfmt.Error("%v", err)
		}

		jobs.PrintJobLogs(logs)

		return nil
	},
}

func init() {
	jobCmd.AddCommand(jobLogsCmd)
}
