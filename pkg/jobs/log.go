package jobs

import (
	"darkroom/pkg/api"
	"darkroom/pkg/config"
	"fmt"
	"strings"
)

type JobLogs struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Logs   string `json:"logs"`
}

func GetJobLogs(cfg *config.Config, jobID int) (*JobLogs, error) {
	client := api.NewClient(cfg)

	endpoint := fmt.Sprintf(config.JobLogsURL, jobID)

	var result JobLogs

	if err := client.Get(endpoint, &result); err != nil {
		return nil, fmt.Errorf(
			"failed to get logs for job %d: %w",
			jobID,
			err,
		)
	}

	return &result, nil
}

func PrintJobLogs(logs *JobLogs) {
	if logs.Logs == "" {
		fmt.Println("No logs available.")
		return
	}

	fmt.Print(logs.Logs)

	if !strings.HasSuffix(logs.Logs, "\n") {
		fmt.Println()
	}
}
