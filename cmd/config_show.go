package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"darkroom/pkg/colorfmt"
	"darkroom/pkg/config"
)

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current darkroom config (decrypted)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return colorfmt.Error("failed to load config: %v", err)
		}

		// Make a copy with sensitive fields redacted
		redacted := *cfg
		redacted.KubeConfig = "<redacted>"
		redacted.S3AccessToken = "<redacted>"
		redacted.AuthToken = "<redacted>"

		out, err := yaml.Marshal(&redacted)
		if err != nil {
			return colorfmt.Error("failed to marshal config: %v", err)
		}

		fmt.Println(string(out))
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
}
