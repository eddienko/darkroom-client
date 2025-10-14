package utils

import (
	"darkroom/pkg/config"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// validateJobInputs checks all job parameters for correctness.
func ValidateJobInputs(cfg *config.Config, jobName, image, script, cpu, memory, jobType, workers string) error {
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
	memVal, err := ParseMemory(memory)
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
		return fmt.Errorf("invalid job: %v", errors.New(strings.Join(errs, "\n")))
	}

	return nil
}
