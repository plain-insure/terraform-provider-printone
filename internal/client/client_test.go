package client

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	// Test with minimal config
	config := Config{}
	client := NewClient(config)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.baseURL != DefaultBaseURL {
		t.Errorf("Expected baseURL to be %s, got %s", DefaultBaseURL, client.baseURL)
	}

	// Test with custom config
	customConfig := Config{
		BaseURL: "https://custom.api.com",
		APIKey:  "test-key",
	}
	customClient := NewClient(customConfig)

	if customClient.baseURL != "https://custom.api.com" {
		t.Errorf("Expected baseURL to be https://custom.api.com, got %s", customClient.baseURL)
	}

	if customClient.apiKey != "test-key" {
		t.Errorf("Expected apiKey to be test-key, got %s", customClient.apiKey)
	}
}

func TestClientConstants(t *testing.T) {
	if DefaultBaseURL != "https://api.print.one" {
		t.Errorf("Expected DefaultBaseURL to be https://api.print.one, got %s", DefaultBaseURL)
	}

	if APIKeyHeader != "x-api-key" {
		t.Errorf("Expected APIKeyHeader to be x-api-key, got %s", APIKeyHeader)
	}

	if EnvAPIKey != "PRINTONE_API_KEY" {
		t.Errorf("Expected EnvAPIKey to be PRINTONE_API_KEY, got %s", EnvAPIKey)
	}
}
