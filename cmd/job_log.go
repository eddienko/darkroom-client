package cmd

import (
	"darkroom/pkg/config"
	"darkroom/pkg/jobs"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	logFollow bool
	logTail   int64
)

var jobLogCmd = &cobra.Command{
	Use:   "log <jobName>",
	Short: "Show logs for a submitted UserJob",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		return jobs.JobLog(cfg, args[0], logFollow, logTail)
	},
}

func init() {
	jobCmd.AddCommand(jobLogCmd)
	jobLogCmd.Flags().BoolVarP(&logFollow, "follow", "f", false, "Follow logs in real-time")
	jobLogCmd.Flags().Int64Var(&logTail, "tail", 0, "Number of lines from the end of logs to show (0 = all)")
}
