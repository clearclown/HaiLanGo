# HaiLanGo - Billing Cap (Hard Limit) Implementation
# Automatically disables billing when budget threshold is exceeded

# Pub/Sub topic for budget alerts
resource "google_pubsub_topic" "budget_alerts" {
  name   = "hailango-budget-alerts"
  labels = local.common_labels

  depends_on = [google_project_service.apis]
}

# Service account for Cloud Function
resource "google_service_account" "billing_cap_function" {
  account_id   = "billing-cap-function"
  display_name = "Billing Cap Cloud Function"
  description  = "Service account for the billing cap Cloud Function"

  depends_on = [google_project_service.apis]
}

# Grant billing admin role to the service account
resource "google_billing_account_iam_member" "billing_admin" {
  billing_account_id = var.billing_account_id
  role               = "roles/billing.admin"
  member             = "serviceAccount:${google_service_account.billing_cap_function.email}"
}

# Grant project billing manager role
resource "google_project_iam_member" "billing_manager" {
  project = var.project_id
  role    = "roles/billing.projectManager"
  member  = "serviceAccount:${google_service_account.billing_cap_function.email}"
}

# Grant Artifact Registry reader role to Cloud Functions service account
resource "google_project_iam_member" "cloudfunctions_artifact_reader" {
  project = var.project_id
  role    = "roles/artifactregistry.reader"
  member  = "serviceAccount:${var.project_id}@appspot.gserviceaccount.com"

  depends_on = [google_project_service.apis]
}

# Cloud Storage bucket for function source code
resource "google_storage_bucket" "function_source" {
  name                        = "${var.project_id}-function-source"
  location                    = var.region
  uniform_bucket_level_access = true
  force_destroy               = true

  labels = local.common_labels

  depends_on = [google_project_service.apis]
}

# Create function source code archive
data "archive_file" "billing_cap_function" {
  type        = "zip"
  output_path = "${path.module}/functions/billing_cap.zip"

  source {
    content  = <<-EOF
      const {CloudBillingClient} = require('@google-cloud/billing');
      const {ProjectsClient} = require('@google-cloud/resource-manager');

      const PROJECT_ID = process.env.GCP_PROJECT || process.env.GOOGLE_CLOUD_PROJECT;
      const PROJECT_NAME = 'projects/' + PROJECT_ID;

      exports.stopBilling = async (pubsubEvent, context) => {
        const pubsubData = JSON.parse(
          Buffer.from(pubsubEvent.data, 'base64').toString()
        );

        console.log('Budget notification received:', JSON.stringify(pubsubData));

        // Check if we've exceeded the budget
        if (pubsubData.costAmount <= pubsubData.budgetAmount) {
          console.log('Current cost (' + pubsubData.costAmount + ') is within budget (' + pubsubData.budgetAmount + '). No action needed.');
          return;
        }

        console.log('Cost (' + pubsubData.costAmount + ') exceeds budget (' + pubsubData.budgetAmount + '). Disabling billing...');

        // Disable billing
        const billingClient = new CloudBillingClient();
        const [billingInfo] = await billingClient.getProjectBillingInfo({
          name: PROJECT_NAME,
        });

        if (billingInfo.billingEnabled) {
          console.log('Billing is currently enabled. Disabling...');
          await billingClient.updateProjectBillingInfo({
            name: PROJECT_NAME,
            projectBillingInfo: {
              billingAccountName: '', // Empty string disables billing
            },
          });
          console.log('Billing has been disabled for project:', PROJECT_ID);
        } else {
          console.log('Billing is already disabled.');
        }
      };
    EOF
    filename = "index.js"
  }

  source {
    content  = <<-EOF
      {
        "name": "billing-cap-function",
        "version": "1.0.0",
        "dependencies": {
          "@google-cloud/billing": "^4.0.0",
          "@google-cloud/resource-manager": "^5.0.0"
        }
      }
    EOF
    filename = "package.json"
  }
}

# Upload function source to bucket
resource "google_storage_bucket_object" "billing_cap_function_source" {
  name   = "billing_cap_${data.archive_file.billing_cap_function.output_md5}.zip"
  bucket = google_storage_bucket.function_source.name
  source = data.archive_file.billing_cap_function.output_path
}

# Cloud Function (Gen1 for simplicity)
resource "google_cloudfunctions_function" "billing_cap" {
  name        = "billing-cap-function"
  description = "Disables billing when budget is exceeded"
  runtime     = "nodejs20"
  region      = var.region

  available_memory_mb   = 256
  source_archive_bucket = google_storage_bucket.function_source.name
  source_archive_object = google_storage_bucket_object.billing_cap_function_source.name
  entry_point           = "stopBilling"

  event_trigger {
    event_type = "google.pubsub.topic.publish"
    resource   = google_pubsub_topic.budget_alerts.id
  }

  service_account_email = google_service_account.billing_cap_function.email

  environment_variables = {
    GCP_PROJECT = var.project_id
  }

  labels = local.common_labels

  depends_on = [
    google_project_service.apis,
    time_sleep.wait_for_apis,
    google_project_iam_member.billing_manager,
    google_project_iam_member.cloudfunctions_artifact_reader,
  ]
}

# Budget with Pub/Sub notification (¥15,000 = ~$100)
resource "google_billing_budget" "hailango_budget" {
  billing_account = var.billing_account_id
  display_name    = "HaiLanGo Budget with Auto-Stop"

  budget_filter {
    projects = ["projects/${var.project_id}"]
  }

  amount {
    specified_amount {
      currency_code = "JPY"
      units         = 15000
    }
  }

  # Alert at 50%, 80%, 90%, 100%
  threshold_rules {
    threshold_percent = 0.5
    spend_basis       = "CURRENT_SPEND"
  }

  threshold_rules {
    threshold_percent = 0.8
    spend_basis       = "CURRENT_SPEND"
  }

  threshold_rules {
    threshold_percent = 0.9
    spend_basis       = "CURRENT_SPEND"
  }

  threshold_rules {
    threshold_percent = 1.0
    spend_basis       = "CURRENT_SPEND"
  }

  # Send notifications to Pub/Sub for automatic action
  all_updates_rule {
    pubsub_topic = google_pubsub_topic.budget_alerts.id
  }

  depends_on = [
    google_project_service.apis,
    google_cloudfunctions_function.billing_cap,
  ]
}

# Outputs
output "billing_cap_function_url" {
  description = "URL of the billing cap function"
  value       = google_cloudfunctions_function.billing_cap.https_trigger_url
}

output "budget_pubsub_topic" {
  description = "Pub/Sub topic for budget alerts"
  value       = google_pubsub_topic.budget_alerts.id
}

output "budget_id" {
  description = "Budget ID"
  value       = google_billing_budget.hailango_budget.id
}
