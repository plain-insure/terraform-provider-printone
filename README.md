# Terraform Provider for PrintOne

This Terraform provider allows you to manage PrintOne API resources using Infrastructure as Code. Currently supports webhook management with an extensible design for additional resources.

## Features

- **Webhook Management**: Create, read, update, and delete webhooks
- **API Key Authentication**: Secure authentication using PrintOne API keys
- **Environment Variable Support**: Configure API key via `PRINTONE_API_KEY` environment variable
- **Extensible Design**: Easy to extend for additional PrintOne API resources

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23
- PrintOne API key

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Using the Provider

### Provider Configuration

```hcl
terraform {
  required_providers {
    printone = {
      source = "plain-insure/printone"
    }
  }
}

provider "printone" {
  # API key for authentication (can also use PRINTONE_API_KEY env var)
  api_key = var.printone_api_key
  
  # Optional: Custom API endpoint (defaults to https://api.print.one)
  # endpoint = "https://api.print.one"
}
```

### Environment Variable Configuration

You can set the API key using an environment variable:

```shell
export PRINTONE_API_KEY="your-api-key"
```

### Webhook Resource

```hcl
resource "printone_webhook" "example" {
  name   = "Order Status Updates"
  url    = "https://api.example.com/webhooks/orders"
  active = true
  events = [
    "order_status_update",
    "batch_status_update"
  ]
}
```

### Webhook Data Source

```hcl
data "printone_webhook" "existing" {
  id = "webhook-id-123"
}

output "webhook_url" {
  value = data.printone_webhook.existing.url
}
```

### Complete Example

See `examples/complete-example.tf` for a comprehensive example showing multiple webhooks and data sources.

## API Authentication

The provider uses the PrintOne API key for authentication via the `x-api-key` header. You can provide the API key in two ways:

1. **Provider Configuration** (not recommended for production):
   ```hcl
   provider "printone" {
     api_key = "your-api-key"
   }
   ```

2. **Environment Variable** (recommended):
   ```shell
   export PRINTONE_API_KEY="your-api-key"
   ```

## Resources and Data Sources

### Resources

- `printone_webhook`: Manage PrintOne webhooks

### Data Sources

- `printone_webhook`: Read existing PrintOne webhooks

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).

To add a new dependency:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
