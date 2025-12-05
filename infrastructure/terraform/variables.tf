# HaiLanGo - Terraform Variables

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
  description = "Monthly budget amount in USD"
  type        = number
  default     = 50  # $50 USD monthly budget
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

# Storage configuration
variable "storage_location" {
  description = "Location for Cloud Storage buckets"
  type        = string
  default     = "ASIA"
}

variable "storage_class" {
  description = "Storage class for buckets"
  type        = string
  default     = "STANDARD"
}

# Alert notification emails
variable "alert_notification_emails" {
  description = "Email addresses for budget alerts"
  type        = list(string)
  default     = []  # Add email addresses for notifications
}
