package jobs

import (
	"bufio"
	"context"
	"darkroom/pkg/auth"
	"darkroom/pkg/config"
	"fmt"
	"io"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// JobLog fetches logs from the Pod(s) of a UserJob
func JobLog(cfg *config.Config, jobName string, follow bool, tail int64) error {
	if cfg.AuthToken == "" {
		return fmt.Errorf("not authenticated, please login first")
	}

	userInfo := auth.GetUserInfo(cfg.AuthToken)
	namespace := "jupyter-" + userInfo.Username

	// Build kubeconfig
	restCfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(cfg.KubeConfig))
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// List Pods with label selector: userJob=<jobName>
	labelSelector := fmt.Sprintf("job-name=%s-runner", jobName)
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return fmt.Errorf("failed to list pods: %w", err)
	}
	if len(pods.Items) == 0 {
		return fmt.Errorf("no pods found for job %s", jobName)
	}

	// Stream or fetch logs from each pod
	for _, pod := range pods.Items {
		fmt.Printf("Logs for pod: %s\n", pod.Name)

		req := clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &v1.PodLogOptions{
			Follow: follow,
			TailLines: func() *int64 {
				if tail > 0 {
					return &tail
				}
				return nil
			}(),
		})

		stream, err := req.Stream(context.TODO())
		if err != nil {
			return fmt.Errorf("failed to open log stream: %w", err)
		}
		defer stream.Close()

		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			return fmt.Errorf("error reading logs: %w", err)
		}
	}

	return nil
}
