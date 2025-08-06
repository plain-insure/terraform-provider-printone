# Copyright (c) HashiCorp, Inc.

provider "printone" {
  # API key for PrintOne API authentication
  # Can also be set via PRINTONE_API_KEY environment variable
  api_key = var.printone_api_key

  # Optional: Custom API endpoint (defaults to https://api.print.one)
  # endpoint = "https://api.print.one"
}
