#!/bin/bash
# HaiLanGo - Budget Creation Script

set -e

PROJECT_ID="${GOOGLE_CLOUD_PROJECT:-hailango}"
BUDGET_AMOUNT="${BUDGET_AMOUNT:-50}"  # USD

echo "=================================="
echo "HaiLanGo Budget Setup"
echo "=================================="
echo "Project: $PROJECT_ID"
echo "Monthly Budget: \$$BUDGET_AMOUNT"
echo ""

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo "Error: gcloud CLI is not installed."
    exit 1
fi

# Get billing account
BILLING_ACCOUNT=$(gcloud billing projects describe $PROJECT_ID --format="value(billingAccountName)" 2>/dev/null | sed 's/billingAccounts\///')

if [ -z "$BILLING_ACCOUNT" ]; then
    echo "Error: No billing account linked to project $PROJECT_ID"
    echo ""
    echo "To link a billing account:"
    echo "  1. List billing accounts: gcloud billing accounts list"
    echo "  2. Link: gcloud billing projects link $PROJECT_ID --billing-account=ACCOUNT_ID"
    exit 1
fi

echo "Billing Account: $BILLING_ACCOUNT"

# Get project number
PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format="value(projectNumber)")
echo "Project Number: $PROJECT_NUMBER"

# Create budget using gcloud (requires alpha/beta components)
echo ""
echo "Creating budget..."

# Check if billing budget API is enabled
gcloud services enable billingbudgets.googleapis.com --quiet 2>/dev/null || true

# Create budget JSON
BUDGET_JSON=$(cat <<EOF
{
  "displayName": "HaiLanGo Monthly Budget",
  "budgetFilter": {
    "projects": ["projects/$PROJECT_NUMBER"]
  },
  "amount": {
    "specifiedAmount": {
      "currencyCode": "USD",
      "units": "$BUDGET_AMOUNT"
    }
  },
  "thresholdRules": [
    {"thresholdPercent": 0.5, "spendBasis": "CURRENT_SPEND"},
    {"thresholdPercent": 0.8, "spendBasis": "CURRENT_SPEND"},
    {"thresholdPercent": 0.9, "spendBasis": "CURRENT_SPEND"},
    {"thresholdPercent": 1.0, "spendBasis": "CURRENT_SPEND"},
    {"thresholdPercent": 1.0, "spendBasis": "FORECASTED_SPEND"}
  ]
}
EOF
)

echo "$BUDGET_JSON" > /tmp/budget.json

# Try to create budget via API
echo ""
echo "Note: Budget creation via CLI may require additional permissions."
echo "If this fails, create budget manually at:"
echo "  https://console.cloud.google.com/billing/budgets?project=$PROJECT_ID"
echo ""

# Using gcloud billing budgets (if available)
if gcloud billing budgets --help &>/dev/null 2>&1; then
    gcloud billing budgets create \
        --billing-account=$BILLING_ACCOUNT \
        --display-name="HaiLanGo Monthly Budget - \$$BUDGET_AMOUNT" \
        --budget-amount="${BUDGET_AMOUNT}USD" \
        --threshold-rules-from-file=/tmp/budget.json \
        2>/dev/null || {
            echo "Failed to create budget via CLI."
            echo "Please create manually at: https://console.cloud.google.com/billing/budgets"
        }
else
    echo "gcloud billing budgets command not available."
    echo ""
    echo "To create budget manually:"
    echo "  1. Go to: https://console.cloud.google.com/billing/budgets?project=$PROJECT_ID"
    echo "  2. Click 'CREATE BUDGET'"
    echo "  3. Set amount to \$$BUDGET_AMOUNT"
    echo "  4. Set thresholds at 50%, 80%, 90%, 100%"
fi

rm -f /tmp/budget.json

echo ""
echo "=================================="
echo "Budget Configuration Complete"
echo "=================================="
echo ""
echo "Budget Dashboard: https://console.cloud.google.com/billing/budgets?project=$PROJECT_ID"
echo "Cost Breakdown: https://console.cloud.google.com/billing/reports?project=$PROJECT_ID"
echo ""
