package netutil

import (
	"darkroom/pkg/config"
	"fmt"
	"io"
	"net/http"
        "strings"

	"sigs.k8s.io/yaml"
)

// FetchKubeconfig retrieves the kubeconfig YAML using the stored auth token
func FetchKubeconfig(cfg *config.Config) (string, error) {
	req, err := http.NewRequest("GET", config.KubeConfigURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.AuthToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch kubeconfig: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	// Unmarshal + Marshal to normalize indentation
	var obj interface{}
	if err := yaml.Unmarshal(raw, &obj); err != nil {
		return "", fmt.Errorf("failed to parse kubeconfig: %w", err)
	}

	out, err := yaml.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("failed to format kubeconfig: %w", err)
	}

        // Convert to string and remove the first line if it's "|"
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "|" {
		lines = lines[1:]
	}

	return strings.Join(lines, "\n"), nil
}
