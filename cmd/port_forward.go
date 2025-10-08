package cmd

import (
	"context"
	"darkroom/pkg/auth"
	"darkroom/pkg/config"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"net/http"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// portForwardCmd represents the port-forward command
var portForwardCmd = &cobra.Command{
	Use:   "port-forward <podName> <localPort>:<podPort>",
	Short: "Forward one or more local ports to a pod",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		podName := args[0]
		portMapping := args[1]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

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
			return fmt.Errorf("failed to load kubeconfig: %w", err)
		}

		clientset, err := kubernetes.NewForConfig(restCfg)
		if err != nil {
			return fmt.Errorf("failed to create kubernetes client: %w", err)
		}

		// Verify pod exists
		pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get pod: %w", err)
		}
		if pod.Status.Phase != v1.PodRunning {
			return fmt.Errorf("pod %s is not running", podName)
		}

		// Build the request URL
		transport, upgrader, err := spdy.RoundTripperFor(restCfg)
		if err != nil {
			return err
		}
		req := clientset.CoreV1().RESTClient().Post().
			Resource("pods").
			Namespace(namespace).
			Name(podName).
			SubResource("portforward")
		url := req.URL()

		// PortForwarder
		stopChan := make(chan struct{}, 1)
		readyChan := make(chan struct{})
		forwarder, err := portforward.NewOnAddresses(
			spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, url),
			[]string{"127.0.0.1"},
			[]string{portMapping},
			stopChan,
			readyChan,
			os.Stdout,
			os.Stderr,
		)
		if err != nil {
			return fmt.Errorf("failed to create port forwarder: %w", err)
		}

		// Trap SIGINT / SIGTERM to stop gracefully
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-signalChan
			close(stopChan)
		}()

		fmt.Printf("Forwarding ports: %s in pod %s/%s\n", portMapping, namespace, podName)
		return forwarder.ForwardPorts()
	},
}

func init() {
	rootCmd.AddCommand(portForwardCmd)
}
