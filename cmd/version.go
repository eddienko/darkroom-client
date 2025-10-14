package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// Fill defaults dynamically if not set by ldflags
func InitVersion() {
	if GitCommit == "unknown" {
		out, err := exec.Command("git", "rev-parse", "HEAD").Output()
		if err == nil {
			GitCommit = strings.TrimSpace(string(out))
		} else {
			GitCommit = "unknown"
		}
	}

	if Version == "dev" {
		out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
		if err == nil {
			Version = strings.TrimSpace(string(out)) + "-" + GitCommit[:7]
		}
		branch, err := exec.Command("git", "branch", "--show-current").Output()
		if err == nil {
			Version = Version + "-" + strings.TrimSpace(string(branch))
		}
	}

	if BuildDate == "unknown" {
		BuildDate = time.Now().UTC().Format(time.RFC3339)
	}
}

func PrintVersion() {
	fmt.Printf("darkroom version: %s\n", Version)
	fmt.Printf("Git commit: %s\n", GitCommit)
	fmt.Printf("Build date: %s\n", BuildDate)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the darkroom version",
	Run: func(cmd *cobra.Command, args []string) {
		InitVersion()
		PrintVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
