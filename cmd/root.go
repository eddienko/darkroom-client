package cmd

import (
	"darkroom/pkg/config"
	"fmt"

	"github.com/spf13/cobra"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "darkroom",
	Short: "Darkroom manages jobs and connections to external services",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		fmt.Printf("Warning: could not load config: %v\n", err)
		cfg = config.New()
	}

	// Persistent flags → override config
	rootCmd.PersistentFlags().StringVar(&cfg.APIEndpoint, "api-endpoint", cfg.APIEndpoint, "API endpoint URL")
	//rootCmd.PersistentFlags().StringVar(&cfg.Namespace, "namespace", cfg.Namespace, "Namespace to use")
	// rootCmd.PersistentFlags().StringVar(&cfg.KubeConfig, "kubeconfig", cfg.KubeConfig, "Path to the kubeconfig file")
}
