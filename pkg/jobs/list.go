package jobs

import (
	"context"
	"darkroom/pkg/auth"
	"darkroom/pkg/config"
	"fmt"
	"strings"

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

	fmt.Printf("%-30s %-10s %-20s\n", "NAME", "STATUS", "SUBMITTED BY")
	fmt.Println(strings.Repeat("-", 70))
	for _, job := range jobList.Items {
		name := job.GetName()
		submittedBy := "<unknown>"
		if ann, ok := job.GetAnnotations()["submitted-by"]; ok {
			submittedBy = ann
		}

		status := "<unknown>"
		if spec, ok := job.Object["status"].(map[string]interface{}); ok {
			if phase, ok := spec["phase"].(string); ok {
				status = phase
			}
		}

		fmt.Printf("%-30s %-10s %-20s\n", name, status, submittedBy)
	}

	return nil
}
