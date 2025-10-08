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

// CancelJob deletes a UserJob by name in the user's namespace
func CancelJob(cfg *config.Config, jobName string) error {
	if cfg.AuthToken == "" {
		return fmt.Errorf("not authenticated, please login first")
	}

	// Get user info
	userInfo, err := auth.GetUserInfo(cfg.AuthToken)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}
	namespace := "jupyter-" + userInfo.Username
	fmt.Printf("Cancelling job %s in namespace %s\n", jobName, namespace)

	// Load kubeconfig
	restCfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(cfg.KubeConfig))
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig content: %w", err)
	}

	// Create dynamic client
	dynClient, err := dynamic.NewForConfig(restCfg)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	// Delete the UserJob CR
	err = dynClient.
		Resource(userJobGVR).
		Namespace(namespace).
		Delete(context.TODO(), jobName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete UserJob: %w", err)
	}

	fmt.Printf("Job %s successfully cancelled\n", jobName)
	return nil
}
