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

// Whitelisted Docker images
var allowedImages = []string{
	"docker.io/6darkroom/jh-darkroom:latest",
	// add more as needed
}

func isImageAllowed(image string) bool {
	for _, allowed := range allowedImages {
		if image == allowed {
			return true
		}
	}
	return false
}

// SubmitJob creates a UserJob custom resource in Kubernetes
func SubmitJob(cfg *config.Config, jobName, image, script, cpu, memory, jobType, workers string) (string, error) {
	// Ensure user is authenticated
	if cfg.AuthToken == "" {
		return "", fmt.Errorf("not authenticated, please login first")
	}

	// Validate image against whitelist
	if !isImageAllowed(image) {
		return "", fmt.Errorf("image %q is not in the whitelist of allowed images", image)
	}

	// Get user info
	userInfo, err := auth.GetUserInfo(cfg.AuthToken)
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %w", err)
	}
	fmt.Printf("Submitting job %s as user: %s\n", jobName, userInfo.Username)

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
	var userJob *unstructured.Unstructured
	switch jobType {
	case "default":
		userJob = &unstructured.Unstructured{
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
	case "batch":
	case "dask":
		userJob = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": fmt.Sprintf("%s/%s", userJobGVR.Group, userJobGVR.Version),
				"kind":       "UserDaskJob",
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
					"scheduler": map[string]interface{}{
						"cpu":    "4",
						"memory": "8GB",
					},
					"worker": map[string]interface{}{
						"cpu":      cpu,
						"memory":   memory,
						"replicas": workers,
					},
				},
			},
		}
	default:
		return "", fmt.Errorf("unknown job type: %s", jobType)
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
