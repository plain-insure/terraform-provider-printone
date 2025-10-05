# Copyright (c) Plain Technologies Aps

terraform {
  required_providers {
    printone = {
      source = "hashicorp.com/plain/printone"
    }
  }
}

provider "printone" {}

//data "printone_webhook" "first_order" {
//  id = ""
//}


resource "printone_webhook" "leads2" {
  name   = "demo-leads2"
  url    = "https://www.example.com/webhook2"
  active = false
  events = ["order_status_update"]
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
    X-Test = "demo2"
  }
  secret_headers = {
    X-Secret-Test = "secret"
  }
  events = ["order_status_update"]
}

//output "first_order" {
//  value = data.printone_webhook.first_order
//}
