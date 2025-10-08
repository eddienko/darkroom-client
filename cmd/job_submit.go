package cmd

import (
	"darkroom/pkg/colorfmt"
	"darkroom/pkg/config"
	"darkroom/pkg/jobs"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	jobName string
	image   string
	script  string
	cpu     string
	memory  string
	// submittedBy string
	jobType string
	workers string
)

// validateJobInputs checks all job parameters for correctness.
func validateJobInputs(cfg *config.Config, jobName, cpu, memory, jobType, workers string) error {
	var errs []string

	// --- Validate job name ---
	validName := regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	if !validName.MatchString(jobName) {
		errs = append(errs, fmt.Sprintf("invalid job name %q: must match [a-z0-9]([-a-z0-9]*[a-z0-9])?", jobName))
	}

	// --- Validate CPU ---
	cpuVal, err := strconv.ParseFloat(cpu, 64)
	if err != nil || cpuVal <= 0 {
		errs = append(errs, fmt.Sprintf("invalid CPU value %q: must be a positive number", cpu))
	} else if cpuVal > cfg.ResourceLimit.MaxCPU {
		errs = append(errs, fmt.Sprintf("requested CPU %.2f exceeds the limit of %.2f", cpuVal, cfg.ResourceLimit.MaxCPU))
	}

	// --- Validate memory ---
	memVal, err := parseMemory(memory)
	if err != nil {
		errs = append(errs, fmt.Sprintf("invalid memory value %q: %v", memory, err))
	} else if memVal > cfg.ResourceLimit.MaxMemory {
		errs = append(errs, fmt.Sprintf("requested memory %dMi exceeds the limit of %dMi", memVal, cfg.ResourceLimit.MaxMemory))
	}

	// --- Validate job type ---
	validTypes := map[string]bool{"default": true, "batch": true, "dask": true}
	if !validTypes[strings.ToLower(jobType)] {
		errs = append(errs, fmt.Sprintf("invalid job type %q: must be one of default, batch, or dask", jobType))
	}

	// --- Validate workers ---
	if workers != "" {
		w, err := strconv.Atoi(workers)
		if err != nil || w < 1 {
			errs = append(errs, fmt.Sprintf("invalid number of workers %q: must be a positive integer", workers))
		}
	}

	// --- Collect errors ---
	if len(errs) > 0 {
		return colorfmt.Error("invalid job: %v", errors.New(strings.Join(errs, "\n")))
	}

	return nil
}

// parseMemory converts values like "512Mi" or "2Gi" into MiB.
func parseMemory(mem string) (int, error) {
	mem = strings.TrimSpace(strings.ToLower(mem))
	if strings.HasSuffix(mem, "gi") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(mem, "gi"), 64)
		if err != nil {
			return 0, err
		}
		return int(v * 1024), nil
	}
	if strings.HasSuffix(mem, "mi") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(mem, "mi"), 64)
		if err != nil {
			return 0, err
		}
		return int(v), nil
	}
	return 0, fmt.Errorf("unknown memory unit (must end with Mi or Gi)")
}

var jobSubmitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a UserJob to the cluster",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateJobInputs(cfg, jobName, cpu, memory, jobType, workers)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name, err := jobs.SubmitJob(cfg, jobName, image, script, cpu, memory, jobType, workers)
		if err != nil {
			return err
		}
		fmt.Printf("UserJob %s submitted successfully\n", name)
		return nil
	},
}

func init() {
	jobCmd.AddCommand(jobSubmitCmd)

	jobSubmitCmd.Flags().StringVar(&jobName, "name", "pi-job", "Job name")
	jobSubmitCmd.Flags().StringVar(&image, "image", "docker.io/6darkroom/jh-darkroom:latest", "Container image")
	jobSubmitCmd.Flags().StringVar(&script, "script", "sleep 3600", "Script to run inside the job")
	jobSubmitCmd.Flags().StringVar(&cpu, "cpu", "1", "CPU request")
	jobSubmitCmd.Flags().StringVar(&memory, "memory", "1Gi", "Memory request")
	jobSubmitCmd.Flags().StringVar(&jobType, "type", "default", "job type (default/batch/dask)")
	jobSubmitCmd.Flags().StringVar(&workers, "nworkers", "", "number of worker nodes (for batch/dask)")
	// jobSubmitCmd.Flags().StringVar(&submittedBy, "submitted-by", "unknown", "Submitter username")
}
