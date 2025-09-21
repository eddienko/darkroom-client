package jobs

import (
	"context"
	"darkroom/pkg/auth"
	"darkroom/pkg/config"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

// Hardcoded Group/Version/Resource for UserJob CRD
var userJobGVR = schema.GroupVersionResource{
	Group:    "edu.dev",
	Version:  "v1",
	Resource: "userjobs",
}

// SubmitJob creates a UserJob custom resource in Kubernetes
func SubmitJob(cfg *config.Config, jobName, image, script, cpu, memory string) (string, error) {
	// Ensure user is authenticated
	if cfg.AuthToken == "" {
		return "", fmt.Errorf("not authenticated, please login first")
	}

	// Get user info
	userInfo := auth.GetUserInfo(cfg.AuthToken)
	fmt.Printf("Submitting job as user: %s\n", userInfo.Username)

	// Load kubeconfig
	namespace := "jupyter-" + userInfo.Username

	restCfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(cfg.KubeConfig))
	if err != nil {
		return "", fmt.Errorf("failed to load kubeconfig content: %w", err)
	}

	// Create dynamic client
	dynClient, err := dynamic.NewForConfig(restCfg)
	if err != nil {
		return "", fmt.Errorf("failed to create dynamic client: %w", err)
	}

	// Build unstructured object representing the CR
	userJob := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": fmt.Sprintf("%s/%s", userJobGVR.Group, userJobGVR.Version),
			"kind":       "UserJob",
			"metadata": map[string]interface{}{
				"name":      jobName,
				"namespace": namespace,
				"annotations": map[string]interface{}{
					"submitted-by": userInfo.Username,
				},
			},
			"spec": map[string]interface{}{
				"image":  image,
				"script": script,
				"resources": map[string]interface{}{
					"cpu":    cpu,
					"memory": memory,
				},
			},
		},
	}

	// Submit to Kubernetes
	created, err := dynClient.
		Resource(userJobGVR).
		Namespace(namespace).
		Create(context.TODO(), userJob, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create UserJob: %w", err)
	}

	return created.GetName(), nil
}
