# HaiLanGo - Terraform Variables (Simplified)

variable "project_id" {
  description = "GCP Project ID"
  type        = string
  default     = "hailango"
}

variable "project_number" {
  description = "GCP Project Number"
  type        = string
  default     = "1077128828066"
}

variable "region" {
  description = "Default GCP region"
  type        = string
  default     = "asia-northeast1"  # Tokyo region for lower latency
}

variable "environment" {
  description = "Environment (development, staging, production)"
  type        = string
  default     = "development"

  validation {
    condition     = contains(["development", "staging", "production"], var.environment)
    error_message = "Environment must be one of: development, staging, production."
  }
}

variable "billing_account_id" {
  description = "Billing Account ID for budget alerts"
  type        = string
  default     = ""  # Set via environment variable or tfvars
}

# Budget configuration
variable "monthly_budget_amount" {
  description = "Monthly budget amount in USD (hard limit)"
  type        = number
  default     = 100  # $100 USD monthly budget (hard limit)
}

variable "budget_alert_thresholds" {
  description = "Budget alert threshold percentages"
  type        = list(number)
  default     = [0.5, 0.8, 0.9, 1.0]  # Alert at 50%, 80%, 90%, 100%
}

# API Keys configuration
variable "create_api_keys" {
  description = "Whether to create API keys"
  type        = bool
  default     = true
}
