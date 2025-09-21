package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"darkroom/pkg/config"
)

var configShowCmd = &cobra.Command{
	Use:   "config show",
	Short: "Show the current darkroom config (decrypted)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		out, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}

		fmt.Println(string(out))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configShowCmd)
}
