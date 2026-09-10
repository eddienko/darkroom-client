package jobs

import (
	"context"
	"darkroom/pkg/config"
	"fmt"
	"math/rand/v2"
	"time"
)

const (
	initialPollInterval = 1 * time.Second
	maxPollInterval     = 60 * time.Second
	backoffFactor       = 2.0
)

func nextPollInterval(current time.Duration) time.Duration {
	next := time.Duration(float64(current) * 2)

	if next > 60*time.Second {
		next = 60 * time.Second
	}

	// Add ±20% jitter.
	jitter := time.Duration(rand.Float64()*0.4*float64(next)) -
		time.Duration(0.2*float64(next))

	return next + jitter
}

func WaitForJob(
	ctx context.Context,
	cfg *config.Config,
	jobID int,
) (*Job, error) {
	interval := initialPollInterval

	var previousStatus string

	for {
		job, err := GetJob(cfg, jobID)
		if err != nil {
			return nil, err
		}

		if job.Status != previousStatus {
			fmt.Printf(
				"Job %d (%s): %s\n",
				job.ID,
				job.Name,
				job.Status,
			)

			previousStatus = job.Status
		}

		switch job.Status {
		case "completed":
			return job, nil

		case "failed":
			return job, fmt.Errorf(
				"job failed: %s",
				job.Message,
			)

		case "cancelled":
			return job, fmt.Errorf(
				"job cancelled: %s",
				job.Message,
			)

		case "pending", "running":
			// Continue polling.

		default:
			return nil, fmt.Errorf(
				"unknown job status %q",
				job.Status,
			)
		}

		timer := time.NewTimer(interval)

		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()

		case <-timer.C:
		}

		interval = nextPollInterval(interval)
	}
}
