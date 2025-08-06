variable "printone_api_key" {
  description = "PrintOne API key for authentication"
  type        = string
  sensitive   = true
  
  # You can set this via environment variable:
  # export TF_VAR_printone_api_key="your-api-key"
  # 
  # Or set the PRINTONE_API_KEY environment variable and omit this variable entirely
}