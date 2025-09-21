package cmd

import (
	"darkroom/pkg/config"
	"darkroom/pkg/jobs"

	"github.com/spf13/cobra"
)

var jobCancelCmd = &cobra.Command{
	Use:   "cancel <jobName>",
	Short: "Cancel a submitted UserJob",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		return jobs.CancelJob(cfg, args[0])
	},
}

func init() {
	jobCmd.AddCommand(jobCancelCmd)
}
