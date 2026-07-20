package cmd

import (
	"darkroom/pkg/colorfmt"
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
			return colorfmt.Error("%v", err)
		}
		err = jobs.CancelJob(cfg, args[0])
		if err != nil {
			return colorfmt.Error("%v", err)
		}
		return nil
	},
}

func init() {
	jobCmd.AddCommand(jobCancelCmd)
}
