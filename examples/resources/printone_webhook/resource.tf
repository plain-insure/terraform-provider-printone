# Copyright (c) Plain Technologies Aps

resource "printone_webhook" "example" {
  name   = "Example Webhook"
  url    = "https://example.com/webhook"
  active = true
  events = [
    "order_status_update",
    "batch_status_update"
  ]

  filters = [
    {
      key   = "status"
      event = "order_status_update"
      type  = "equals"
      value = "order_created"
    },
    {
      key    = "status"
      event  = "batch_status_update"
      type   = "in"
      values = ["batch_created", "batch_processed"]
    }
  ]
}
