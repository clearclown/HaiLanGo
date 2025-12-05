# HaiLanGo - API Keys Configuration

# Service Account for API access
resource "google_service_account" "hailango_api" {
  count = var.create_api_keys ? 1 : 0

  account_id   = "hailango-api-sa"
  display_name = "HaiLanGo API Service Account"
  description  = "Service account for HaiLanGo API access"
  project      = var.project_id

  depends_on = [time_sleep.wait_for_apis]
}

# IAM roles for the service account
resource "google_project_iam_member" "hailango_api_roles" {
  for_each = var.create_api_keys ? toset([
    "roles/cloudvision.user",           # Cloud Vision API
    "roles/cloudtranslate.user",        # Cloud Translation API
    "roles/storage.objectViewer",        # Read storage objects
  ]) : toset([])

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.hailango_api[0].email}"

  depends_on = [google_service_account.hailango_api]
}

# Custom role for TTS and STT (more granular permissions)
resource "google_project_iam_custom_role" "tts_stt_user" {
  count = var.create_api_keys ? 1 : 0

  role_id     = "hailangoTtsSttUser"
  title       = "HaiLanGo TTS/STT User"
  description = "Custom role for Text-to-Speech and Speech-to-Text access"
  project     = var.project_id

  permissions = [
    "texttospeech.synthesize",
    "speech.recognize",
    "speech.longrunningrecognize",
  ]

  depends_on = [time_sleep.wait_for_apis]
}

# Assign custom role to service account
resource "google_project_iam_member" "hailango_tts_stt" {
  count = var.create_api_keys ? 1 : 0

  project = var.project_id
  role    = google_project_iam_custom_role.tts_stt_user[0].id
  member  = "serviceAccount:${google_service_account.hailango_api[0].email}"

  depends_on = [
    google_service_account.hailango_api,
    google_project_iam_custom_role.tts_stt_user,
  ]
}

# Service account key (for local development)
resource "google_service_account_key" "hailango_api_key" {
  count = var.create_api_keys && var.environment == "development" ? 1 : 0

  service_account_id = google_service_account.hailango_api[0].name
  key_algorithm      = "KEY_ALG_RSA_2048"

  depends_on = [google_service_account.hailango_api]
}

# Store API key in Secret Manager (for production use)
resource "google_secret_manager_secret" "api_credentials" {
  count = var.create_api_keys ? 1 : 0

  secret_id = "hailango-api-credentials"
  project   = var.project_id

  labels = local.common_labels

  replication {
    auto {}
  }

  depends_on = [time_sleep.wait_for_apis]
}

resource "google_secret_manager_secret_version" "api_credentials" {
  count = var.create_api_keys && var.environment == "development" ? 1 : 0

  secret      = google_secret_manager_secret.api_credentials[0].id
  secret_data = google_service_account_key.hailango_api_key[0].private_key

  depends_on = [google_secret_manager_secret.api_credentials]
}

# API Key for client-side access (restricted)
resource "google_apikeys_key" "hailango_web_key" {
  count = var.create_api_keys ? 1 : 0

  name         = "hailango-web-key"
  display_name = "HaiLanGo Web API Key"
  project      = var.project_id

  restrictions {
    # API restrictions - only allow specific APIs
    api_targets {
      service = "vision.googleapis.com"
    }
    api_targets {
      service = "texttospeech.googleapis.com"
    }
    api_targets {
      service = "speech.googleapis.com"
    }
    api_targets {
      service = "translate.googleapis.com"
    }

    # Browser key restrictions (add your domains)
    # browser_key_restrictions {
    #   allowed_referrers = [
    #     "localhost:3000",
    #     "*.hailango.com",
    #   ]
    # }
  }

  depends_on = [time_sleep.wait_for_apis]
}

# Outputs
output "service_account_email" {
  description = "Service account email"
  value       = var.create_api_keys ? google_service_account.hailango_api[0].email : null
  sensitive   = false
}

output "web_api_key" {
  description = "Web API Key (restricted)"
  value       = var.create_api_keys ? google_apikeys_key.hailango_web_key[0].key_string : null
  sensitive   = true
}

output "service_account_key_json" {
  description = "Service account key JSON (base64 encoded)"
  value       = var.create_api_keys && var.environment == "development" ? google_service_account_key.hailango_api_key[0].private_key : null
  sensitive   = true
}
