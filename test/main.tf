# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    printone = {
      source = "hashicorp.com/plain/printone"
    }
  }
}

provider "printone" {}

data "printone_webhook" "first_order" {
  id = "1"
}

output "first_order" {
  value = data.printone_webhook.first_order
}