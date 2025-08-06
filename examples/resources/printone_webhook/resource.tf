resource "printone_webhook" "example" {
  name   = "Example Webhook"
  url    = "https://example.com/webhook"
  active = true
  events = [
    "order_status_update",
    "batch_status_update"
  ]
}
