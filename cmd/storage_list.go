package cmd

import (
	"darkroom/pkg/storage"

	"github.com/spf13/cobra"
)

var storageListCmd = &cobra.Command{
	Use:     "ls [path]",
	Example: "  darkroom storage ls /projects/myproject",
	Short:   "List contents of a storage path",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 { // path provided, list contents of the path
			path = args[0]
			err := storage.List(cfg, path)
			if err != nil {
				return err
			}
		} else { // no path provided, list buckets
			err := storage.List(cfg, path)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	storageCmd.AddCommand(storageListCmd)

	// jobSubmitCmd.Flags().StringVar(&jobName, "name", "pi-job", "Job name")
	// jobSubmitCmd.Flags().StringVar(&image, "image", "docker.io/6darkroom/jh-darkroom:latest", "Container image")
	// jobSubmitCmd.Flags().StringVar(&script, "script", "sleep 3600", "Script to run inside the job")
	// jobSubmitCmd.Flags().StringVar(&cpu, "cpu", "1", "CPU request")
	// jobSubmitCmd.Flags().StringVar(&memory, "memory", "1Gi", "Memory request")
	// jobSubmitCmd.Flags().StringVar(&submittedBy, "submitted-by", "unknown", "Submitter username")
}
