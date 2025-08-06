package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	DefaultBaseURL = "https://api.print.one"
	APIKeyHeader   = "x-api-key"
	EnvAPIKey      = "PRINTONE_API_KEY"
)

// Client represents a PrintOne API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// Config represents configuration for the PrintOne client
type Config struct {
	BaseURL string
	APIKey  string
}

// NewClient creates a new PrintOne API client
func NewClient(config Config) *Client {
	// Use default base URL if not provided
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	// Use API key from config or environment variable
	apiKey := config.APIKey
	if apiKey == "" {
		apiKey = os.Getenv(EnvAPIKey)
	}

	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// makeRequest makes an HTTP request to the PrintOne API
func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set(APIKeyHeader, c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return resp, nil
}

// handleResponse handles the HTTP response and unmarshals the body if successful
func (c *Client) handleResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if target != nil {
		if err := json.Unmarshal(body, target); err != nil {
			return fmt.Errorf("error unmarshaling response: %w", err)
		}
	}

	return nil
}
