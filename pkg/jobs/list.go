package jobs

import (
	"context"
	"darkroom/pkg/auth"
	"darkroom/pkg/config"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

// ListJobs lists all UserJobs in the user's namespace
func ListJobs(cfg *config.Config) error {
	if cfg.AuthToken == "" {
		return fmt.Errorf("not authenticated, please login first")
	}

	userInfo, err := auth.GetUserInfo(cfg.AuthToken)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}
	namespace := "jupyter-" + userInfo.Username

	restCfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(cfg.KubeConfig))
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig content: %w", err)
	}

	dynClient, err := dynamic.NewForConfig(restCfg)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	jobList, err := dynClient.Resource(userJobGVR).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list UserJobs: %w", err)
	}

	if len(jobList.Items) == 0 {
		fmt.Println("No jobs found")
		return nil
	}

	fmt.Printf("%-30s %-10s %-20s %-25s %-25s\n", "NAME", "STATUS", "SUBMITTED BY", "STARTED", "COMPLETED")
	fmt.Println(strings.Repeat("-", 109))
	for _, job := range jobList.Items {
		name := job.GetName()
		submittedBy := "<unknown>"
		if ann, ok := job.GetAnnotations()["submitted-by"]; ok {
			submittedBy = ann
		}

		status := "<unknown>"
		startTime := ""
		endTime := ""
		if spec, ok := job.Object["status"].(map[string]interface{}); ok {
			if phase, ok := spec["phase"].(string); ok {
				status = phase
			}
			if ann, ok := spec["startTime"].(string); ok {
				startTime = formatTime(ann)
			}
			if ann, ok := spec["completionTime"].(string); ok {
				endTime = formatTime(ann)
			}
		}

		fmt.Printf("%-30s %-10s %-20s %-25s %-25s\n", name, status, submittedBy, startTime, endTime)
	}

	return nil
}

func formatTime(s string) string {
	if s == "" {
		return ""
	}

	// parse ISO-like timestamp
	layout := "2006-01-02T15:04:05.000000Z"
	t, err := time.Parse(layout, s)
	if err != nil {
		return ""
	}

	// format in a readable way
	return t.Format("02 Jan 2006 15:04:05")
}
