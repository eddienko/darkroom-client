package cmd

import (
	"darkroom/pkg/colorfmt"
	"darkroom/pkg/config"
	"darkroom/pkg/jobs"
	"darkroom/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v2"
)

var (
	jobName     string
	image       string
	script      string
	cpu         string
	memory      string
	gpu         int
	command     []string
	commandArgs []string
	// submittedBy string
	jobType      string
	workers      string
	outputFormat string
)

type JobSpec struct {
	Name      string            `yaml:"name"`
	Image     string            `yaml:"image"`
	Command   []string          `yaml:"command"`
	Args      []string          `yaml:"args"`
	Script    string            `yaml:"script"`
	Workers   string            `yaml:"workers"`
	JobType   string            `yaml:"job_type"`
	Env       map[string]string `yaml:"env"`
	Resources ResourcesSpec     `yaml:"resources"`
	Outputs   []OutputSpec      `yaml:"outputs"`
}

type ResourcesSpec struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
	GPU    int    `yaml:"gpu"`
}

type OutputSpec struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Type string `yaml:"type"`
}

type JobSubmitResult struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func loadJobSpec(path string) (JobSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return JobSpec{}, err
	}

	var spec JobSpec

	if err := yaml.Unmarshal(data, &spec); err != nil {
		return JobSpec{}, fmt.Errorf(
			"invalid job YAML: %w", err,
		)
	}

	if spec.Name == "" {
		return JobSpec{}, fmt.Errorf(
			"job name is required",
		)
	}

	if spec.Image == "" {
		spec.Image = "docker.io/6darkroom/jh-darkroom:latest"
	}

	return spec, nil
}

// validateJobInputs checks all job parameters for correctness.
func validateJobInputs(cfg *config.Config, jobName, image, script, cpu, memory string, gpu int, jobType, workers, outputFormat string) error {
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

	// --- Validate GPU ---
	if gpu < 0 {
		errs = append(errs, fmt.Sprintf("invalid GPU value %d: must be a non-negative integer", gpu))
	}

	// Optional limit check
	// else if gpuVal > cfg.ResourceLimit.MaxGPU {
	//     errs = append(errs, fmt.Sprintf(
	//         "requested GPUs %d exceeds the limit of %d",
	//         gpuVal,
	//         cfg.ResourceLimit.MaxGPU,
	//     ))
	// }

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

	// --- Validate output format ---
	validFormats := map[string]bool{"text": true, "json": true}
	if !validFormats[strings.ToLower(outputFormat)] {
		errs = append(errs, fmt.Sprintf("invalid output format %q: must be one of text or json", outputFormat))
	}

	// --- Collect errors ---
	if len(errs) > 0 {
		return colorfmt.Error("invalid job: %v", errors.New(strings.Join(errs, "\n")))
	}

	return nil
}

var jobSubmitCmd = &cobra.Command{
	Use:   "submit [job.yaml]",
	Short: "Submit a UserJob to the cluster",
	Args:  cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// YAML mode
		if len(args) == 1 {
			spec, err := loadJobSpec(args[0])
			if err != nil {
				return err
			}

			// Populate the existing CLI variables from YAML.
			jobName = spec.Name
			image = spec.Image

			command = nil
			commandArgs = nil
			script = spec.Script

			if len(spec.Command) > 0 {
				command = spec.Command
				commandArgs = spec.Args
				script = ""
			}

			cpu = spec.Resources.CPU
			memory = spec.Resources.Memory
			gpu = spec.Resources.GPU
			jobType = spec.JobType
			workers = spec.Workers

			if cmd.Flags().Changed("cpu") {
				cpu, _ = cmd.Flags().GetString("cpu")
			}
			if cmd.Flags().Changed("memory") {
				memory, _ = cmd.Flags().GetString("memory")
			}
			if cmd.Flags().Changed("gpu") {
				gpu, _ = cmd.Flags().GetInt("gpu")
			}
			if cmd.Flags().Changed("name") {
				jobName, _ = cmd.Flags().GetString("name")
			}

		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		return validateJobInputs(cfg, jobName, image, script, cpu, memory, gpu, jobType, workers, outputFormat)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		resp, err := jobs.SubmitJob(cfg, jobName, image, script, command, commandArgs, cpu, memory, gpu, jobType, workers)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			jsonOutput, err := json.Marshal(resp)
			if err != nil {
				return fmt.Errorf("failed to convert response to JSON: %w", err)
			}
			fmt.Println(string(jsonOutput))
		} else {
			fmt.Printf("UserJob %s submitted successfully. Status %s.\n", resp.Name, resp.Status)
		}

		return nil
	},
}

func init() {
	jobCmd.AddCommand(jobSubmitCmd)

	jobSubmitCmd.Flags().StringVar(&jobName, "name", "", "Job name")
	jobSubmitCmd.Flags().StringVar(&image, "image", "docker.io/6darkroom/jh-darkroom:latest", "Container image")
	jobSubmitCmd.Flags().StringVar(&script, "script", "sleep 600", "Script to execute")
	jobSubmitCmd.Flags().StringVar(&cpu, "cpu", "1", "CPU request")
	jobSubmitCmd.Flags().StringVar(&memory, "memory", "1Gi", "Memory request")
	jobSubmitCmd.Flags().IntVar(&gpu, "gpu", 0, "Number of GPUs to request")
	jobSubmitCmd.Flags().StringVar(&jobType, "type", "default", "job type (default/batch/dask)")
	jobSubmitCmd.Flags().StringVar(&workers, "nworkers", "", "number of worker nodes (for batch/dask)")
	jobSubmitCmd.Flags().StringVar(&outputFormat, "output", "text", "Output format (text/json)")
}
