package cmd

import (
	"darkroom/pkg/config"
	"darkroom/pkg/jobs"
	"fmt"

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
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Run OpenShell
		if err := jobs.OpenShell(cfg, jobName); err != nil {
			return fmt.Errorf("failed to open shell: %w", err)
		}
		return nil
	},
}

func init() {
	jobCmd.AddCommand(jobShellCmd)
}
