package cmd

import (
	"darkroom/pkg/auth"
	"fmt"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

func promptInput(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}

func promptPassword(prompt string) []byte {
	fmt.Print(prompt)
	pass, err := gopass.GetPasswdMasked()
	if err != nil {
		return nil
	}
	return pass
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with external service",
	RunE: func(cmd *cobra.Command, args []string) error {
		username := promptInput("Username: ")
		password := promptPassword("Password: ")

		if err := auth.Login(cfg, username, password); err != nil {
			return err
		}

		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Println("Login successful")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
