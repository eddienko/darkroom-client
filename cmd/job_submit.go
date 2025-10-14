package cmd

import (
	"darkroom/pkg/colorfmt"
	"darkroom/pkg/config"
	"darkroom/pkg/jobs"
	"darkroom/pkg/utils"
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
	gui     bool
)

// validateJobInputs checks all job parameters for correctness.
func validateJobInputs(cfg *config.Config, jobName, image, script, cpu, memory, jobType, workers string) error {
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
	memVal, err := utils.ParseMemory(memory)
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

var jobSubmitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a UserJob to the cluster",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if gui {
			return nil
		}
		return validateJobInputs(cfg, jobName, image, script, cpu, memory, jobType, workers)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if gui {
			data, err := jobs.SubmitGUI(cfg, jobName, image, script, cpu, memory, jobType, workers)
			if err != nil {
				return err
			}
			if data == nil {
				return nil
			}
			return nil
		}

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

	jobSubmitCmd.Flags().StringVar(&jobName, "name", "", "Job name")
	jobSubmitCmd.Flags().StringVar(&image, "image", "docker.io/6darkroom/jh-darkroom:latest", "Container image")
	jobSubmitCmd.Flags().StringVar(&script, "script", "sleep 600", "Script to run inside the job")
	jobSubmitCmd.Flags().StringVar(&cpu, "cpu", "1", "CPU request")
	jobSubmitCmd.Flags().StringVar(&memory, "memory", "1Gi", "Memory request")
	jobSubmitCmd.Flags().StringVar(&jobType, "type", "default", "job type (default/batch/dask)")
	jobSubmitCmd.Flags().StringVar(&workers, "nworkers", "", "number of worker nodes (for batch/dask)")
	jobSubmitCmd.Flags().BoolVarP(&gui, "gui", "", false, "launch a GUI")
}
