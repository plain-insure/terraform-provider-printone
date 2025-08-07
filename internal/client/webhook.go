// Copyright (c) Plain Technologies Aps

package client

import (
	"context"
	"fmt"
	"net/http"
)

// WebhookRequest represents a webhook creation/update request.
type WebhookRequest struct {
	Name          string                 `json:"name"`
	URL           string                 `json:"url"`
	Active        bool                   `json:"active"`
	Events        []string               `json:"events"`
	Headers       map[string]interface{} `json:"headers,omitempty"`
	SecretHeaders map[string]interface{} `json:"secretHeaders,omitempty"`
	Filters       []interface{}          `json:"filters,omitempty"`
}

// WebhookResponse represents a webhook response from the API.
type WebhookResponse struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	URL           string                 `json:"url"`
	Active        bool                   `json:"active"`
	Events        []string               `json:"events"`
	Headers       map[string]interface{} `json:"headers"`
	SecretHeaders map[string]interface{} `json:"secretHeaders"`
	SuccessRate   interface{}            `json:"successRate"`
	Filters       []interface{}          `json:"filters,omitempty"`
}

// CreateWebhook creates a new webhook.
func (c *Client) CreateWebhook(ctx context.Context, webhook *WebhookRequest) (*WebhookResponse, error) {
	resp, err := c.makeRequest(ctx, http.MethodPost, "/v2/webhooks", webhook)
	if err != nil {
		return nil, err
	}

	var result WebhookResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetWebhook retrieves a webhook by ID.
func (c *Client) GetWebhook(ctx context.Context, id string) (*WebhookResponse, error) {
	path := fmt.Sprintf("/v2/webhooks/%s", id)
	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result WebhookResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateWebhook updates an existing webhook.
func (c *Client) UpdateWebhook(ctx context.Context, id string, webhook *WebhookRequest) (*WebhookResponse, error) {
	path := fmt.Sprintf("/v2/webhooks/%s", id)
	resp, err := c.makeRequest(ctx, http.MethodPatch, path, webhook)
	if err != nil {
		return nil, err
	}

	var result WebhookResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteWebhook deletes a webhook by ID.
func (c *Client) DeleteWebhook(ctx context.Context, id string) error {
	path := fmt.Sprintf("/v2/webhooks/%s", id)
	resp, err := c.makeRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return c.handleResponse(resp, nil)
}

// WebhookSecretResponse represents a webhook secret response from the API.
type WebhookSecretResponse struct {
	Secret string `json:"secret"`
}

// GetWebhookSecret retrieves the webhook secret.
func (c *Client) GetWebhookSecret(ctx context.Context) (*WebhookSecretResponse, error) {
	resp, err := c.makeRequest(ctx, http.MethodPost, "/v2/webhooks/secret", nil)
	if err != nil {
		return nil, err
	}

	var result WebhookSecretResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
