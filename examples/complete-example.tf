terraform {
  required_providers {
    printone = {
      source = "plain-insure/printone"
    }
  }
}

provider "printone" {
  # API key can be set via PRINTONE_API_KEY environment variable
  # or passed directly (not recommended for production)
  api_key = var.printone_api_key
}

# Create a webhook for order status updates
resource "printone_webhook" "order_updates" {
  name   = "Order Status Updates"
  url    = "https://api.example.com/webhooks/order-updates"
  active = true
  events = [
    "order_status_update"
  ]
}

# Create a webhook for batch status updates
resource "printone_webhook" "batch_updates" {
  name   = "Batch Status Updates"
  url    = "https://api.example.com/webhooks/batch-updates"
  active = true
  events = [
    "batch_status_update"
  ]
}

# Create a comprehensive webhook that listens to multiple events
resource "printone_webhook" "comprehensive" {
  name   = "Comprehensive Webhook"
  url    = "https://api.example.com/webhooks/all-events"
  active = true
  events = [
    "order_status_update",
    "batch_status_update",
    "template_preview_rendered",
    "coupon_code_used"
  ]
}

# Data source to read an existing webhook
data "printone_webhook" "existing" {
  id = "existing-webhook-id"
}

# Output webhook information
output "order_webhook_id" {
  description = "ID of the order updates webhook"
  value       = printone_webhook.order_updates.id
}

output "existing_webhook_name" {
  description = "Name of the existing webhook"
  value       = data.printone_webhook.existing.name
}

output "existing_webhook_url" {
  description = "URL of the existing webhook"
  value       = data.printone_webhook.existing.url
}