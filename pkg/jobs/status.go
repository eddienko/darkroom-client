package jobs

import (
	"darkroom/pkg/api"
	"darkroom/pkg/config"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"sigs.k8s.io/yaml"
)

type JobDetail struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Status         string     `json:"status"`
	Message        string     `json:"message"`
	Image          string     `json:"image"`
	Command        []string   `json:"command"`
	Args           []string   `json:"args"`
	CPU            string     `json:"cpu"`
	Memory         string     `json:"memory"`
	GPU            int        `json:"gpu"`
	JobType        string     `json:"job_type"`
	Workers        int        `json:"workers"`
	StartTime      *time.Time `json:"start_time"`
	CompletionTime *time.Time `json:"completion_time"`
}

func GetJob(cfg *config.Config, jobID int) (*Job, error) {
	client := api.NewClient(cfg)

	endpoint := fmt.Sprintf(config.JobDetailURL, jobID)

	var job Job

	if err := client.Get(endpoint, &job); err != nil {
		return nil, fmt.Errorf("failed to get job %d: %w", jobID, err)
	}

	return &job, nil
}

func PrintJobStatus(job *Job) {
	fmt.Printf("ID:           %d\n", job.ID)
	fmt.Printf("Name:         %s\n", job.Name)
	fmt.Printf("Status:       %s\n", job.Status)
	fmt.Printf("Message:      %s\n", job.Message)
	fmt.Printf("Image:        %s\n", job.Image)

	if len(job.Command) > 0 {
		fmt.Printf("Command:      %s\n", strings.Join(job.Command, " "))
	}

	if len(job.Args) > 0 {
		fmt.Printf("Args:         %s\n", strings.Join(job.Args, " "))
	}

	fmt.Printf("CPU:          %s\n", job.CPU)
	fmt.Printf("Memory:       %s\n", job.Memory)
	fmt.Printf("GPU:          %d\n", job.GPU)
	fmt.Printf("Type:          %s\n", job.JobType)
	fmt.Printf("Workers:      %d\n", job.Workers)

	if job.StartTime != nil {
		fmt.Printf(
			"Start time:   %s\n",
			job.StartTime.Local().Format("2006-01-02 15:04:05"),
		)
	}

	if job.CompletionTime != nil {
		fmt.Printf(
			"Completion:   %s\n",
			job.CompletionTime.Local().Format("2006-01-02 15:04:05"),
		)
	}
}

func jsonToYAML(jsonStr string) (string, error) {
	var data interface{}

	// Decode JSON into a generic structure
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Encode as YAML
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return string(yamlBytes), nil
}
