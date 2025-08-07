# Copyright (c) Plain Technologies Aps

data "printone_webhooksecret" "example" {
  # The webhook secret is retrieved from the PrintOne API
  # No configuration parameters are required
}

# Output the webhook secret (be careful with sensitive data)
output "webhook_secret" {
  value     = data.printone_webhooksecret.example.secret
  sensitive = true
}