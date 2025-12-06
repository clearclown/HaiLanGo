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
  project               = var.project_id
  region                = var.region
  user_project_override = true
  billing_project       = var.project_id
}

provider "google-beta" {
  project               = var.project_id
  region                = var.region
  user_project_override = true
  billing_project       = var.project_id
}

# Local values
locals {
  # APIs required for HaiLanGo (API access only, no GCP hosting)
  required_apis = [
    "vision.googleapis.com",           # Cloud Vision API (OCR)
    "texttospeech.googleapis.com",     # Cloud Text-to-Speech API
    "speech.googleapis.com",           # Cloud Speech-to-Text API
    "translate.googleapis.com",        # Cloud Translation API
    "cloudbilling.googleapis.com",     # Cloud Billing API (for budget)
    "billingbudgets.googleapis.com",   # Cloud Billing Budget API
    "cloudresourcemanager.googleapis.com", # Resource Manager API
    "serviceusage.googleapis.com",     # Service Usage API
    "iam.googleapis.com",              # IAM API (for API keys)
    "apikeys.googleapis.com",          # API Keys API
    "monitoring.googleapis.com",       # Cloud Monitoring API
    "pubsub.googleapis.com",           # Pub/Sub API (for budget alerts)
    "cloudfunctions.googleapis.com",   # Cloud Functions API (for auto-shutdown)
    "cloudbuild.googleapis.com",       # Cloud Build API (for deploying functions)
    "run.googleapis.com",              # Cloud Run API (for Gen2 functions)
    "eventarc.googleapis.com",         # Eventarc API (for Gen2 functions)
    "artifactregistry.googleapis.com", # Artifact Registry API (for Cloud Functions)
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
