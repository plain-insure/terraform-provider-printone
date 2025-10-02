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
  filters = [
    {
      event  = "order_status_update"
      key    = "isBillable"
      type   = "equals"
      value  = "false"
      values = null
    },
    {
      event = "order_status_update"
      key   = "status"
      type  = "in"
      value = null
      values = [
        "order_created",
        "order_ready",
      ]
    },
  ]

  headers = {
    X-Test = "demo21"
  }
  secret_headers = {
    X-Secret-Test = "secret"
  }
  events = ["order_status_update"]
}

output "first_order" {
  value = data.printone_webhook.first_order
}
