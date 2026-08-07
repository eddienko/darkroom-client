package api

import (
	"bytes"
	"darkroom/pkg/auth"
	"darkroom/pkg/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	cfg    *config.Config
	client *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) Do(method, path string, request any, response any) error {

	// Ensure we have a valid access token
	if err := auth.EnsureAccessToken(c.cfg); err != nil {
		return err
	}

	var body io.Reader

	if request != nil {
		data, err := json.Marshal(request)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		body = bytes.NewReader(data)
	}

	url := path

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.AuthToken)
	req.Header.Set("Accept", "application/json")

	if request != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s %s failed (%d): %s",
			method,
			path,
			resp.StatusCode,
			string(data),
		)
	}

	if response != nil && len(data) > 0 {
		if err := json.Unmarshal(data, response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) Get(path string, response any) error {
	return c.Do(http.MethodGet, path, nil, response)
}

func (c *Client) Post(path string, request, response any) error {
	return c.Do(http.MethodPost, path, request, response)
}

func (c *Client) Delete(path string, response any) error {
	return c.Do(http.MethodDelete, path, nil, response)
}

func (c *Client) Put(path string, request, response any) error {
	return c.Do(http.MethodPut, path, request, response)
}
