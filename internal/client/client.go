package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	serverURL string
	timeout   time.Duration
	logger    *zap.Logger
	httpClient *http.Client
}

func New(serverURL string, timeoutSeconds int, logger *zap.Logger) *Client {
	return &Client{
		serverURL: serverURL,
		timeout:   time.Duration(timeoutSeconds) * time.Second,
		logger:    logger,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

func (c *Client) Get(endpoint string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s", c.serverURL, endpoint)
	c.logger.Debug("GET request", zap.String("url", url))

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return result, nil
}

func (c *Client) Post(endpoint string, data interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s", c.serverURL, endpoint)
	c.logger.Debug("POST request", zap.String("url", url))

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return result, nil
}

func (c *Client) CheckHealth() (map[string]interface{}, error) {
	return c.Get("health")
}

func (c *Client) ExecuteCommand(command string, useCache bool) (map[string]interface{}, error) {
	return c.Post("api/command", map[string]interface{}{
		"command":   command,
		"use_cache": useCache,
	})
}
