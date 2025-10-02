package cmd

import (
	"darkroom/pkg/jobs"
	"fmt"

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

var jobSubmitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a UserJob to the cluster",
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
