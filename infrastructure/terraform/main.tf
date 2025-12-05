# HaiLanGo - Google Cloud Infrastructure
# Terraform configuration for managing GCP resources

terraform {
  required_version = ">= 1.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 5.0"
    }
  }

  # Backend configuration for state management
  # Uncomment and configure when ready for production
  # backend "gcs" {
  #   bucket = "hailango-terraform-state"
  #   prefix = "terraform/state"
  # }
}

# Provider configuration
provider "google" {
  project = var.project_id
  region  = var.region
}

provider "google-beta" {
  project = var.project_id
  region  = var.region
}

# Local values
locals {
  # APIs required for HaiLanGo
  required_apis = [
    "vision.googleapis.com",           # Cloud Vision API (OCR)
    "texttospeech.googleapis.com",     # Cloud Text-to-Speech API
    "speech.googleapis.com",           # Cloud Speech-to-Text API
    "translate.googleapis.com",        # Cloud Translation API
    "storage.googleapis.com",          # Cloud Storage
    "cloudbilling.googleapis.com",     # Cloud Billing API
    "billingbudgets.googleapis.com",   # Cloud Billing Budget API
    "cloudresourcemanager.googleapis.com", # Resource Manager API
    "serviceusage.googleapis.com",     # Service Usage API
    "iam.googleapis.com",              # IAM API
    "secretmanager.googleapis.com",    # Secret Manager API
  ]

  # Common labels for all resources
  common_labels = {
    project     = "hailango"
    environment = var.environment
    managed_by  = "terraform"
  }
}

# Enable required APIs
resource "google_project_service" "apis" {
  for_each = toset(local.required_apis)

  project = var.project_id
  service = each.value

  disable_dependent_services = false
  disable_on_destroy         = false
}

# Wait for APIs to be enabled
resource "time_sleep" "wait_for_apis" {
  depends_on = [google_project_service.apis]

  create_duration = "30s"
}
