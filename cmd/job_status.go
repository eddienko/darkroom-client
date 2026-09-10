package cmd

import (
	"darkroom/pkg/colorfmt"
	"darkroom/pkg/config"
	"darkroom/pkg/jobs"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var output string

var jobStatusCmd = &cobra.Command{
	Use:   "status JOB_ID",
	Short: "Show the status of a job",
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

		job, err := jobs.GetJob(cfg, jobID)
		if err != nil {
			return colorfmt.Error("%v", err)
		}

		switch output {
		case "json":
			jobJSON, err := json.MarshalIndent(job, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal job to JSON: %w", err)
			}
			fmt.Println(string(jobJSON))
		default:
			jobs.PrintJobStatus(job)
		}

		return nil
	},
}

func init() {
	jobCmd.AddCommand(jobStatusCmd)

	jobStatusCmd.Flags().StringVar(
		&output,
		"output",
		"text",
		"output format: text or json",
	)
}
