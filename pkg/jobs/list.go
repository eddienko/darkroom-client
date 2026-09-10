package jobs

import (
	"darkroom/pkg/api"
	"darkroom/pkg/config"
	"fmt"
	"net/url"
	"time"
)

type Job struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Namespace      string     `json:"namespace"`
	Image          string     `json:"image"`
	Command        []string   `json:"command"`
	Args           []string   `json:"args"`
	CPU            string     `json:"cpu"`
	Memory         string     `json:"memory"`
	GPU            int        `json:"gpu"`
	JobType        string     `json:"job_type"`
	Workers        int        `json:"workers"`
	Status         string     `json:"status"`
	Message        string     `json:"message"`
	StartTime      *time.Time `json:"start_time"`
	CompletionTime *time.Time `json:"completion_time"`
}

// ListJobs lists all UserJobs in the user's namespace
func ListJobs(cfg *config.Config, jobStatus string) ([]Job, error) {
	client := api.NewClient(cfg)
	u, err := url.Parse(config.JobListURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	if jobStatus != "" {
		q := u.Query()
		q.Set("status", jobStatus)
		u.RawQuery = q.Encode()
	}

	endpoint := u.String()

	var result []Job

	err = client.Get(endpoint, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	// fmt.Printf("Job List Response: %+v\n", result)

	return result, nil

}

func PrintJobs(jobs []Job) {
	if len(jobs) == 0 {
		fmt.Println("No jobs found.")
		return
	}

	fmt.Printf(
		"%-6s %-20s %-10s %-8s %-8s %-6s %-20s\n",
		"ID",
		"NAME",
		"STATUS",
		"CPU",
		"MEMORY",
		"GPU",
		"STARTED",
	)

	for _, job := range jobs {
		started := "-"

		if job.StartTime != nil {
			started = job.StartTime.Local().Format("2006-01-02 15:04")
		}

		fmt.Printf(
			"%-6d %-20s %-10s %-8s %-8s %-6d %-20s\n",
			job.ID,
			job.Name,
			job.Status,
			job.CPU,
			job.Memory,
			job.GPU,
			started,
		)
	}
}
