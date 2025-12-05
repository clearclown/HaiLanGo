# HaiLanGo - API Quotas and Rate Limiting Configuration

# Note: Google Cloud API quotas are managed through the console or gcloud CLI
# This file documents the recommended quota settings and monitoring

# Monitoring alert policy for API quota usage
resource "google_monitoring_alert_policy" "api_quota_alert" {
  display_name = "HaiLanGo API Quota Alert"
  project      = var.project_id

  combiner = "OR"

  conditions {
    display_name = "Vision API Quota Usage"

    condition_threshold {
      filter          = "resource.type=\"consumer_quota\" AND resource.labels.service=\"vision.googleapis.com\""
      duration        = "60s"
      comparison      = "COMPARISON_GT"
      threshold_value = 80

      aggregations {
        alignment_period     = "300s"
        per_series_aligner   = "ALIGN_PERCENT_CHANGE"
        cross_series_reducer = "REDUCE_MAX"
      }

      trigger {
        count = 1
      }
    }
  }

  conditions {
    display_name = "TTS API Quota Usage"

    condition_threshold {
      filter          = "resource.type=\"consumer_quota\" AND resource.labels.service=\"texttospeech.googleapis.com\""
      duration        = "60s"
      comparison      = "COMPARISON_GT"
      threshold_value = 80

      aggregations {
        alignment_period     = "300s"
        per_series_aligner   = "ALIGN_PERCENT_CHANGE"
        cross_series_reducer = "REDUCE_MAX"
      }

      trigger {
        count = 1
      }
    }
  }

  conditions {
    display_name = "STT API Quota Usage"

    condition_threshold {
      filter          = "resource.type=\"consumer_quota\" AND resource.labels.service=\"speech.googleapis.com\""
      duration        = "60s"
      comparison      = "COMPARISON_GT"
      threshold_value = 80

      aggregations {
        alignment_period     = "300s"
        per_series_aligner   = "ALIGN_PERCENT_CHANGE"
        cross_series_reducer = "REDUCE_MAX"
      }

      trigger {
        count = 1
      }
    }
  }

  # Alert documentation
  documentation {
    content   = <<-EOT
      API quota usage is high. Consider:
      1. Implementing more aggressive caching
      2. Rate limiting user requests
      3. Requesting quota increase from Google Cloud

      Dashboard: https://console.cloud.google.com/apis/dashboard?project=${var.project_id}
    EOT
    mime_type = "text/markdown"
  }

  # Notification channels would be added here
  # notification_channels = [google_monitoring_notification_channel.email.id]

  depends_on = [time_sleep.wait_for_apis]
}

# API usage dashboard
resource "google_monitoring_dashboard" "api_usage" {
  dashboard_json = jsonencode({
    displayName = "HaiLanGo API Usage Dashboard"
    gridLayout = {
      columns = 2
      widgets = [
        {
          title = "Vision API Requests"
          xyChart = {
            dataSets = [{
              timeSeriesQuery = {
                timeSeriesFilter = {
                  filter = "resource.type=\"consumed_api\" AND resource.labels.service=\"vision.googleapis.com\""
                  aggregation = {
                    perSeriesAligner   = "ALIGN_RATE"
                    alignmentPeriod    = "60s"
                    crossSeriesReducer = "REDUCE_SUM"
                  }
                }
              }
            }]
          }
        },
        {
          title = "TTS API Requests"
          xyChart = {
            dataSets = [{
              timeSeriesQuery = {
                timeSeriesFilter = {
                  filter = "resource.type=\"consumed_api\" AND resource.labels.service=\"texttospeech.googleapis.com\""
                  aggregation = {
                    perSeriesAligner   = "ALIGN_RATE"
                    alignmentPeriod    = "60s"
                    crossSeriesReducer = "REDUCE_SUM"
                  }
                }
              }
            }]
          }
        },
        {
          title = "STT API Requests"
          xyChart = {
            dataSets = [{
              timeSeriesQuery = {
                timeSeriesFilter = {
                  filter = "resource.type=\"consumed_api\" AND resource.labels.service=\"speech.googleapis.com\""
                  aggregation = {
                    perSeriesAligner   = "ALIGN_RATE"
                    alignmentPeriod    = "60s"
                    crossSeriesReducer = "REDUCE_SUM"
                  }
                }
              }
            }]
          }
        },
        {
          title = "Translation API Requests"
          xyChart = {
            dataSets = [{
              timeSeriesQuery = {
                timeSeriesFilter = {
                  filter = "resource.type=\"consumed_api\" AND resource.labels.service=\"translate.googleapis.com\""
                  aggregation = {
                    perSeriesAligner   = "ALIGN_RATE"
                    alignmentPeriod    = "60s"
                    crossSeriesReducer = "REDUCE_SUM"
                  }
                }
              }
            }]
          }
        },
        {
          title = "Estimated Daily Cost"
          scorecard = {
            timeSeriesQuery = {
              timeSeriesFilter = {
                filter = "resource.type=\"global\" AND metric.type=\"billing.googleapis.com/cost\""
              }
            }
          }
        },
        {
          title = "API Error Rate"
          xyChart = {
            dataSets = [{
              timeSeriesQuery = {
                timeSeriesFilter = {
                  filter = "resource.type=\"consumed_api\" AND metric.labels.response_code!=\"200\""
                  aggregation = {
                    perSeriesAligner   = "ALIGN_RATE"
                    alignmentPeriod    = "300s"
                    crossSeriesReducer = "REDUCE_SUM"
                  }
                }
              }
            }]
          }
        }
      ]
    }
  })
  project = var.project_id

  depends_on = [time_sleep.wait_for_apis]
}

# Output dashboard URL
output "monitoring_dashboard_url" {
  description = "URL to the API monitoring dashboard"
  value       = "https://console.cloud.google.com/monitoring/dashboards/builder/${google_monitoring_dashboard.api_usage.id}?project=${var.project_id}"
}
