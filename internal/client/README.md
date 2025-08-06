# PrintOne HTTP Client

This package provides an HTTP client for interacting with the PrintOne API.

## Features

- API key authentication via `x-api-key` header
- Support for environment variable fallback (`PRINTONE_API_KEY`)
- Webhook CRUD operations
- Extensible design for additional resources

## Usage

### Basic Configuration

```go
import "github.com/plain-insure/terraform-provider-printone/internal/client"

// Create client with API key
config := client.Config{
    APIKey: "your-api-key",
}
client := client.NewClient(config)

// Or use environment variable PRINTONE_API_KEY
config := client.Config{}
client := client.NewClient(config)
```

### Custom Base URL

```go
config := client.Config{
    BaseURL: "https://custom.api.print.one",
    APIKey:  "your-api-key",
}
client := client.NewClient(config)
```

## Webhook Operations

### Create Webhook

```go
webhook := &client.WebhookRequest{
    Name:   "My Webhook",
    URL:    "https://example.com/webhook",
    Active: true,
    Events: []string{"order_status_update"},
}

response, err := client.CreateWebhook(ctx, webhook)
```

### Get Webhook

```go
response, err := client.GetWebhook(ctx, "webhook-id")
```

### Update Webhook

```go
webhook := &client.WebhookRequest{
    Name:   "Updated Webhook",
    URL:    "https://example.com/webhook",
    Active: false,
    Events: []string{"order_status_update", "batch_status_update"},
}

response, err := client.UpdateWebhook(ctx, "webhook-id", webhook)
```

### Delete Webhook

```go
err := client.DeleteWebhook(ctx, "webhook-id")
```

## Error Handling

The client returns detailed error messages for API failures, including HTTP status codes and response bodies for debugging.

## Extension

To add support for additional PrintOne API resources, create new methods following the existing pattern:

1. Define request/response structs
2. Add methods to the Client struct
3. Use the common `makeRequest` and `handleResponse` methods