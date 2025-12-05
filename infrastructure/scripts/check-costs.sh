#!/bin/bash
# HaiLanGo - Cost Monitoring Script

set -e

PROJECT_ID="${GOOGLE_CLOUD_PROJECT:-hailango}"

echo "=================================="
echo "HaiLanGo Cost Report"
echo "=================================="
echo "Project: $PROJECT_ID"
echo "Date: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo "Error: gcloud CLI is not installed."
    exit 1
fi

# Get billing account
BILLING_ACCOUNT=$(gcloud billing projects describe $PROJECT_ID --format="value(billingAccountName)" 2>/dev/null | sed 's/billingAccounts\///')

if [ -z "$BILLING_ACCOUNT" ]; then
    echo "Warning: No billing account linked to project"
else
    echo "Billing Account: $BILLING_ACCOUNT"
fi

echo ""
echo "API Usage (Last 24 hours)"
echo "--------------------------"

# Check API usage
APIS=("vision.googleapis.com" "texttospeech.googleapis.com" "speech.googleapis.com" "translate.googleapis.com")

for api in "${APIS[@]}"; do
    API_NAME=$(echo $api | sed 's/.googleapis.com//')
    echo ""
    echo "📊 $API_NAME:"

    # Get API metrics (simplified - actual metrics require Cloud Monitoring API)
    STATUS=$(gcloud services list --enabled --filter="name:$api" --format="value(state)" 2>/dev/null || echo "UNKNOWN")
    echo "   Status: $STATUS"

    # Note: Detailed usage requires Cloud Monitoring API queries
    # gcloud monitoring metrics list --filter="metric.type has '$api'"
done

echo ""
echo "Storage Usage"
echo "--------------"

# List storage buckets and their sizes
echo ""
gcloud storage ls --project=$PROJECT_ID 2>/dev/null | while read bucket; do
    if [[ $bucket == gs://${PROJECT_ID}* ]]; then
        SIZE=$(gcloud storage du --summarize "$bucket" 2>/dev/null | awk '{print $1}' || echo "N/A")
        echo "  $bucket: $SIZE"
    fi
done

echo ""
echo "Quick Links"
echo "-----------"
echo "📈 Cost Breakdown: https://console.cloud.google.com/billing/reports?project=$PROJECT_ID"
echo "📊 API Dashboard: https://console.cloud.google.com/apis/dashboard?project=$PROJECT_ID"
echo "💰 Budgets: https://console.cloud.google.com/billing/budgets?project=$PROJECT_ID"
echo "📉 Quotas: https://console.cloud.google.com/apis/api/vision.googleapis.com/quotas?project=$PROJECT_ID"
echo ""

# Cost optimization tips
echo "Cost Optimization Tips"
echo "----------------------"
echo "1. Enable caching for TTS/OCR results"
echo "2. Use lifecycle rules for storage (auto-delete old files)"
echo "3. Monitor API quotas to avoid unexpected charges"
echo "4. Consider batch processing for OCR operations"
echo "5. Use lower-quality TTS for non-critical audio"
echo ""

# API pricing reference
echo "API Pricing Reference (as of 2024)"
echo "-----------------------------------"
echo "Vision API:"
echo "  - Text Detection: \$1.50 per 1,000 images"
echo "  - Document Text: \$1.50 per 1,000 pages"
echo ""
echo "Text-to-Speech:"
echo "  - Standard: \$4.00 per 1M characters"
echo "  - WaveNet: \$16.00 per 1M characters"
echo "  - Neural2: \$16.00 per 1M characters"
echo ""
echo "Speech-to-Text:"
echo "  - Standard: \$0.024 per minute"
echo "  - Enhanced: \$0.048 per minute"
echo ""
echo "Cloud Translation:"
echo "  - Basic: \$20 per 1M characters"
echo "  - Advanced: \$80 per 1M characters"
echo ""
echo "Cloud Storage (Standard):"
echo "  - Storage: \$0.020 per GB/month"
echo "  - Operations: \$0.05 per 10K (Class A)"
echo ""
