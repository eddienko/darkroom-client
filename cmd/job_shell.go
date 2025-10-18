package cmd

import (
	"darkroom/pkg/colorfmt"
	"darkroom/pkg/config"
	"darkroom/pkg/jobs"

	"github.com/spf13/cobra"
)

var jobShellCmd = &cobra.Command{
	Use:   "shell <jobname>",
	Short: "Open an interactive shell inside a job's pod",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jobName := args[0]

		// Load config
		cfg, err := config.Load()
		if err != nil {
			return colorfmt.Error("failed to load config: %v", err)
		}

		// Run OpenShell
		if err := jobs.OpenShell(cfg, jobName); err != nil {
			return colorfmt.Error("failed to open shell: %v", err)
		}
		return nil
	},
}

func init() {
	jobCmd.AddCommand(jobShellCmd)
}
