# Copyright (c) Plain Technologies Aps

terraform {
  required_providers {
    printone = {
      source = "hashicorp.com/plain/printone"
    }
  }
}

provider "printone" {}

data "printone_webhook" "first_order" {
  id = "a989f2f7-8572-4b9a-89ca-12524003f596"
}

resource "printone_webhook" "leads" {
  name   = "demo-leads"
  url    = "https://www.example.com/webhook"
  active = false
  /*
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
  */
  headers = {
    X-Test = "demo"
  }
  secret_headers = {
    X-Secret-Test = "secret"
  }
  events = ["order_status_update"]
}

output "first_order" {
  sensitive = true
  value = data.printone_webhook.first_order
}
