package cmd

import (
	"darkroom/pkg/netutil"
	"fmt"

	"github.com/spf13/cobra"
)

var kubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Fetch and print the Kubernetes config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.AuthToken == "" {
			return fmt.Errorf("no auth token found, please run `darkroom login` first")
		}

		data, err := netutil.FetchKubeconfig(cfg)
		if err != nil {
			return err
		}

		fmt.Println(data)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(kubeconfigCmd)
}
