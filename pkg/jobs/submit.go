package jobs

import (
	"darkroom/pkg/api"
	"darkroom/pkg/config"
	"fmt"

	"github.com/mattn/go-shellwords"
)

// Whitelisted Docker images
var allowedImages = []string{
	"docker.io/6darkroom/jh-darkroom:latest",
	// add more as needed
}

type SubmitJobRequest struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Command []string `json:"command"`
	Args    []string `json:"args"`
	CPU     string   `json:"cpu"`
	Memory  string   `json:"memory"`
	GPU     int      `json:"gpu"`
	JobType string   `json:"job_type"`
	Workers string   `json:"workers,omitempty"`
}

type SubmitJobResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func isImageAllowed(image string) bool {
	for _, allowed := range allowedImages {
		if image == allowed {
			return true
		}
	}
	return false
}

// SubmitJob creates a UserJob custom resource in Kubernetes
func SubmitJob(cfg *config.Config, jobName, image, script string, command, commandArgs []string, cpu, memory string, gpu int, jobType, workers string) (SubmitJobResponse, error) {
	client := api.NewClient(cfg)

	var parsedCommand []string
	var parsedArgs []string

	if script == "" && len(command) == 0 {
		return SubmitJobResponse{}, fmt.Errorf("either script or command must be provided")
	}

	if script == "" {
		parsedCommand = command
		parsedArgs = commandArgs
	} else {
		parsedCommand, _ = shellwords.Parse("/usr/local/bin/entrypoint.sh")
		parsedArgs, _ = shellwords.Parse(script)

		if !isImageAllowed(image) {
			return SubmitJobResponse{}, fmt.Errorf("image %q is not in the whitelist of allowed images", image)
		}

	}

	req := SubmitJobRequest{
		Name:    jobName,
		Image:   image,
		Command: parsedCommand,
		Args:    parsedArgs,
		CPU:     cpu,
		Memory:  memory,
		GPU:     gpu,
		JobType: jobType,
		Workers: workers,
	}

	var resp SubmitJobResponse

	err := client.Post(config.JobSubmitURL, req, &resp)

	if err != nil {
		return SubmitJobResponse{}, err
	}

	return resp, nil

}
