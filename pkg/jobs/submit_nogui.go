//go:build !gui
// +build !gui

package jobs

import (
	"darkroom/pkg/config"
	"fmt"
)

type FormData struct {
	Name    string `json:"name"`
	Script  string `json:"script"`
	CPU     int    `json:"cpu"`
	Memory  int    `json:"memory"`
	Image   string `json:"image"`
	JobType string `json:"jobtype"`
	Workers string `json:"workers"`
}

func SubmitGUI(cfg *config.Config, jobName, image, script, cpu, memory, jobType, workers string) (*FormData, error) {
	return nil, fmt.Errorf("this version does not support GUI")
}
