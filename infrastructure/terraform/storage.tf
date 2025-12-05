# HaiLanGo - Cloud Storage Configuration

# Storage bucket for user uploads (books, images)
resource "google_storage_bucket" "user_uploads" {
  name          = "${var.project_id}-user-uploads-${var.environment}"
  location      = var.storage_location
  storage_class = var.storage_class
  project       = var.project_id

  # Enable uniform bucket-level access
  uniform_bucket_level_access = true

  # Lifecycle rules for cost optimization
  lifecycle_rule {
    condition {
      age = 30  # Move to nearline after 30 days
    }
    action {
      type          = "SetStorageClass"
      storage_class = "NEARLINE"
    }
  }

  lifecycle_rule {
    condition {
      age = 90  # Move to coldline after 90 days
    }
    action {
      type          = "SetStorageClass"
      storage_class = "COLDLINE"
    }
  }

  # Versioning for data protection
  versioning {
    enabled = true
  }

  # CORS configuration for web uploads
  cors {
    origin          = ["*"]  # Restrict in production
    method          = ["GET", "HEAD", "PUT", "POST", "DELETE"]
    response_header = ["*"]
    max_age_seconds = 3600
  }

  labels = local.common_labels

  depends_on = [time_sleep.wait_for_apis]
}

# Storage bucket for generated audio (TTS cache)
resource "google_storage_bucket" "audio_cache" {
  name          = "${var.project_id}-audio-cache-${var.environment}"
  location      = var.storage_location
  storage_class = var.storage_class
  project       = var.project_id

  uniform_bucket_level_access = true

  # Auto-delete old audio files after 7 days
  lifecycle_rule {
    condition {
      age = 7
    }
    action {
      type = "Delete"
    }
  }

  # CORS for audio streaming
  cors {
    origin          = ["*"]
    method          = ["GET", "HEAD"]
    response_header = ["*"]
    max_age_seconds = 86400
  }

  labels = local.common_labels

  depends_on = [time_sleep.wait_for_apis]
}

# Storage bucket for OCR results cache
resource "google_storage_bucket" "ocr_cache" {
  name          = "${var.project_id}-ocr-cache-${var.environment}"
  location      = var.storage_location
  storage_class = var.storage_class
  project       = var.project_id

  uniform_bucket_level_access = true

  # Auto-delete OCR results after 30 days
  lifecycle_rule {
    condition {
      age = 30
    }
    action {
      type = "Delete"
    }
  }

  labels = local.common_labels

  depends_on = [time_sleep.wait_for_apis]
}

# IAM bindings for storage buckets
resource "google_storage_bucket_iam_member" "user_uploads_access" {
  count = var.create_api_keys ? 1 : 0

  bucket = google_storage_bucket.user_uploads.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.hailango_api[0].email}"

  depends_on = [
    google_storage_bucket.user_uploads,
    google_service_account.hailango_api,
  ]
}

resource "google_storage_bucket_iam_member" "audio_cache_access" {
  count = var.create_api_keys ? 1 : 0

  bucket = google_storage_bucket.audio_cache.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.hailango_api[0].email}"

  depends_on = [
    google_storage_bucket.audio_cache,
    google_service_account.hailango_api,
  ]
}

resource "google_storage_bucket_iam_member" "ocr_cache_access" {
  count = var.create_api_keys ? 1 : 0

  bucket = google_storage_bucket.ocr_cache.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.hailango_api[0].email}"

  depends_on = [
    google_storage_bucket.ocr_cache,
    google_service_account.hailango_api,
  ]
}

# Outputs
output "user_uploads_bucket" {
  description = "User uploads bucket name"
  value       = google_storage_bucket.user_uploads.name
}

output "audio_cache_bucket" {
  description = "Audio cache bucket name"
  value       = google_storage_bucket.audio_cache.name
}

output "ocr_cache_bucket" {
  description = "OCR cache bucket name"
  value       = google_storage_bucket.ocr_cache.name
}
