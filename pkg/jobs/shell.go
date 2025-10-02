package jobs

import (
	"context"
	"darkroom/pkg/auth"
	"darkroom/pkg/config"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// OpenShell attaches to a job's pod and starts an interactive shell
func OpenShell(cfg *config.Config, jobName string) error {
	if cfg.AuthToken == "" {
		return fmt.Errorf("not authenticated, please login first")
	}

	userInfo := auth.GetUserInfo(cfg.AuthToken)
	namespace := "jupyter-" + userInfo.Username

	// Load kubeconfig
	restCfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(cfg.KubeConfig))
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Find pods for the job
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s-runner", jobName),
	})
	if err != nil {
		return fmt.Errorf("failed to list pods for job %s: %w", jobName, err)
	}
	if len(pods.Items) == 0 {
		return fmt.Errorf("no pods found for job %s", jobName)
	}

	// Pick the first pod (usually only one per Job)
	pod := pods.Items[0]
	container := pod.Spec.Containers[0].Name

	fmt.Printf("Opening shell in pod %s, container %s\n", pod.Name, container)

	// Exec request
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   []string{"/bin/sh"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("failed to initialize executor: %w", err)
	}

	// Attach stdio
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		return fmt.Errorf("error executing shell: %w", err)
	}

	return nil
}
