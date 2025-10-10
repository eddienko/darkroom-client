package jobs

import (
	"context"
	"darkroom/pkg/auth"
	"darkroom/pkg/config"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

// JobStatus prints detailed info for a single UserJob
func JobStatus(cfg *config.Config, jobName string) error {
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

	job, err := dynClient.Resource(userJobGVR).Namespace(namespace).Get(context.TODO(), jobName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get UserJob %s: %w", jobName, err)
	}

	fmt.Println("Job Name:", job.GetName())
	fmt.Println("Namespace:", namespace)

	// Submitted by annotation
	submittedBy := "<unknown>"
	if ann, ok := job.GetAnnotations()["submitted-by"]; ok {
		submittedBy = ann
	}
	fmt.Println("Submitted by:", submittedBy)

	// Status
	status := "<unknown>"
	startTime := "<unknown>"
	completionTime := "<unknown>"
	message := "<unknown>"
	logs := ""
	if s, ok := job.Object["status"].(map[string]interface{}); ok {
		if phase, ok := s["phase"].(string); ok {
			status = phase
		}
		if st, ok := s["startTime"].(string); ok {
			startTime = st
		}
		if ct, ok := s["completionTime"].(string); ok {
			completionTime = ct
		}
		if ct, ok := s["message"].(string); ok {
			message = ct
		}
		if ct, ok := s["logs"].(string); ok {
			logs = ct
		}
	}
	fmt.Println("Status:", status)
	fmt.Println("Message:", message)
	fmt.Println("Start Time:", startTime)
	fmt.Println("Completion Time:", completionTime)

	// Spec info
	if spec, ok := job.Object["spec"].(map[string]interface{}); ok {
		if script, ok := spec["script"].(string); ok {
			fmt.Println("Script:", script)
		}
		if image, ok := spec["image"].(string); ok {
			fmt.Println("Image:", image)
		}
		if resources, ok := spec["resources"].(map[string]interface{}); ok {
			fmt.Println("Resources:")
			if cpu, ok := resources["cpu"].(string); ok {
				fmt.Println("  CPU:", cpu)
			}
			if mem, ok := resources["memory"].(string); ok {
				fmt.Println("  Memory:", mem)
			}
		}
	}

	fmt.Println("Logs:\n", logs)

	return nil
}
