package cmd

import (
	"fmt"
	"strings"

	"darkroom/pkg/config"

	"github.com/spf13/cobra"
)

var configSetCmd = &cobra.Command{
	Use:   "set KEY=VALUE [KEY=VALUE...]",
	Short: "Set one or more configuration values",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		for _, arg := range args {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid argument %q, must be KEY=VALUE", arg)
			}
			key := parts[0]
			value := parts[1]

			if err := cfg.UpdateField(key, value); err != nil {
				return fmt.Errorf("failed to update field %q: %w", key, err)
			}

			fmt.Printf("%s updated successfully\n", key)
		}

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
}
