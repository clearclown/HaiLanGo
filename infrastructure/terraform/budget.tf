# HaiLanGo - Budget and Cost Management

# Budget for the project
resource "google_billing_budget" "hailango_budget" {
  count = var.billing_account_id != "" ? 1 : 0

  billing_account = var.billing_account_id
  display_name    = "HaiLanGo Monthly Budget - ${var.environment}"

  budget_filter {
    projects = ["projects/${var.project_number}"]

    # Optional: Filter by specific services
    # services = [
    #   "services/24E6-581D-38E5",  # Cloud Vision API
    #   "services/7F2E-BE19-E8E0",  # Cloud Text-to-Speech API
    # ]
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units         = tostring(var.monthly_budget_amount)
    }
  }

  # Alert thresholds
  dynamic "threshold_rules" {
    for_each = var.budget_alert_thresholds
    content {
      threshold_percent = threshold_rules.value
      spend_basis       = "CURRENT_SPEND"
    }
  }

  # Also alert on forecasted spend
  threshold_rules {
    threshold_percent = 1.0
    spend_basis       = "FORECASTED_SPEND"
  }

  # Notification channels (email)
  all_updates_rule {
    # Pub/Sub topic for programmatic alerts (optional)
    # pubsub_topic = google_pubsub_topic.budget_alerts.id

    # Schema version for notifications
    schema_version = "1.0"

    # Disable default IAM recipients if needed
    disable_default_iam_recipients = false

    # Add monitoring notification channels if configured
    # monitoring_notification_channels = []
  }

  depends_on = [google_project_service.apis]
}

# Optional: Pub/Sub topic for budget alerts (for programmatic handling)
# resource "google_pubsub_topic" "budget_alerts" {
#   name = "hailango-budget-alerts"
#
#   labels = local.common_labels
#
#   depends_on = [google_project_service.apis]
# }

# Output budget information
output "budget_id" {
  description = "Budget ID"
  value       = var.billing_account_id != "" ? google_billing_budget.hailango_budget[0].id : null
}

output "budget_display_name" {
  description = "Budget display name"
  value       = var.billing_account_id != "" ? google_billing_budget.hailango_budget[0].display_name : null
}
