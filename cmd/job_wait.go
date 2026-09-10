package cmd

import (
	"darkroom/pkg/colorfmt"
	"darkroom/pkg/config"
	"darkroom/pkg/jobs"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var jobWaitCmd = &cobra.Command{
	Use:   "wait JOB_ID",
	Short: "Wait for a job to finish",
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

		fmt.Printf("Waiting for job %d...\n", jobID)

		job, err := jobs.WaitForJob(
			cmd.Context(),
			cfg,
			jobID,
		)

		if err != nil {
			return colorfmt.Error("%v", err)
		}

		fmt.Printf(
			"Job %d (%s) completed successfully.\n",
			job.ID,
			job.Name,
		)

		return nil
	},
}

func init() {
	jobCmd.AddCommand(jobWaitCmd)
}
