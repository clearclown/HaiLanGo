#!/bin/bash
# HaiLanGo - Google Cloud Setup Script

set -e

PROJECT_ID="${GOOGLE_CLOUD_PROJECT:-hailango}"
REGION="${GOOGLE_CLOUD_REGION:-asia-northeast1}"

echo "=================================="
echo "HaiLanGo Google Cloud Setup"
echo "=================================="
echo "Project: $PROJECT_ID"
echo "Region: $REGION"
echo ""

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo "Error: gcloud CLI is not installed."
    echo "Install from: https://cloud.google.com/sdk/install"
    exit 1
fi

# Check if authenticated
if ! gcloud auth list --filter="status:ACTIVE" --format="value(account)" 2>/dev/null | grep -q "@"; then
    echo "Not authenticated. Starting login..."
    gcloud auth login
fi

# Set project
echo "Setting project to $PROJECT_ID..."
gcloud config set project $PROJECT_ID

# Set region
echo "Setting default region to $REGION..."
gcloud config set compute/region $REGION

# Enable required APIs
echo ""
echo "Enabling required APIs..."

APIS=(
    "vision.googleapis.com"
    "texttospeech.googleapis.com"
    "speech.googleapis.com"
    "translate.googleapis.com"
    "storage.googleapis.com"
    "cloudbilling.googleapis.com"
    "billingbudgets.googleapis.com"
    "cloudresourcemanager.googleapis.com"
    "serviceusage.googleapis.com"
    "iam.googleapis.com"
    "secretmanager.googleapis.com"
    "monitoring.googleapis.com"
    "apikeys.googleapis.com"
)

for api in "${APIS[@]}"; do
    echo "  Enabling $api..."
    gcloud services enable $api --quiet || echo "    Warning: Failed to enable $api"
done

echo ""
echo "Waiting for APIs to propagate..."
sleep 10

# Create service account
SA_NAME="hailango-api-sa"
SA_EMAIL="${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

echo ""
echo "Creating service account: $SA_NAME..."
if ! gcloud iam service-accounts describe $SA_EMAIL &>/dev/null; then
    gcloud iam service-accounts create $SA_NAME \
        --display-name="HaiLanGo API Service Account" \
        --description="Service account for HaiLanGo API access"
    echo "  Service account created: $SA_EMAIL"
else
    echo "  Service account already exists: $SA_EMAIL"
fi

# Grant roles to service account
echo ""
echo "Granting roles to service account..."

ROLES=(
    "roles/cloudvision.user"
    "roles/cloudtranslate.user"
    "roles/storage.objectAdmin"
)

for role in "${ROLES[@]}"; do
    echo "  Granting $role..."
    gcloud projects add-iam-policy-binding $PROJECT_ID \
        --member="serviceAccount:$SA_EMAIL" \
        --role="$role" \
        --quiet || echo "    Warning: Failed to grant $role"
done

# Create API key (for client-side use)
echo ""
echo "Creating API key..."
API_KEY_NAME="hailango-web-key"

# Check if API key exists
EXISTING_KEY=$(gcloud services api-keys list --filter="displayName:$API_KEY_NAME" --format="value(name)" 2>/dev/null || true)

if [ -z "$EXISTING_KEY" ]; then
    gcloud services api-keys create \
        --display-name="$API_KEY_NAME" \
        --api-target="service=vision.googleapis.com" \
        --api-target="service=texttospeech.googleapis.com" \
        --api-target="service=speech.googleapis.com" \
        --api-target="service=translate.googleapis.com" \
        2>/dev/null || echo "  Warning: Failed to create API key"
    echo "  API key created"
else
    echo "  API key already exists"
fi

# Create service account key for local development
echo ""
echo "Creating service account key for local development..."
KEY_FILE="./credentials.json"

if [ ! -f "$KEY_FILE" ]; then
    gcloud iam service-accounts keys create $KEY_FILE \
        --iam-account=$SA_EMAIL \
        2>/dev/null || echo "  Warning: Failed to create service account key"
    echo "  Key saved to: $KEY_FILE"
    echo "  WARNING: Keep this file secure and never commit to git!"
else
    echo "  Key file already exists: $KEY_FILE"
fi

# Get billing account
echo ""
echo "Checking billing account..."
BILLING_ACCOUNT=$(gcloud billing accounts list --filter="open:true" --format="value(ACCOUNT_ID)" --limit=1 2>/dev/null || true)

if [ -n "$BILLING_ACCOUNT" ]; then
    echo "  Billing account: $BILLING_ACCOUNT"

    # Link billing account to project
    LINKED=$(gcloud billing projects describe $PROJECT_ID --format="value(billingAccountName)" 2>/dev/null || true)
    if [ -z "$LINKED" ]; then
        echo "  Linking billing account to project..."
        gcloud billing projects link $PROJECT_ID --billing-account=$BILLING_ACCOUNT 2>/dev/null || echo "  Warning: Failed to link billing account"
    else
        echo "  Billing account already linked"
    fi
else
    echo "  Warning: No billing account found. Some APIs may not work without billing."
fi

# Print API key
echo ""
echo "=================================="
echo "Setup Complete!"
echo "=================================="
echo ""
echo "Service Account: $SA_EMAIL"
echo "Credentials File: $KEY_FILE"
echo ""
echo "To use the credentials in your application:"
echo "  export GOOGLE_APPLICATION_CREDENTIALS=$KEY_FILE"
echo ""
echo "API Key (for client-side use):"
gcloud services api-keys list --filter="displayName:$API_KEY_NAME" --format="value(keyString)" 2>/dev/null || echo "  Run 'gcloud services api-keys list' to see your API keys"
echo ""
echo "Console Links:"
echo "  APIs: https://console.cloud.google.com/apis/dashboard?project=$PROJECT_ID"
echo "  Billing: https://console.cloud.google.com/billing?project=$PROJECT_ID"
echo "  Storage: https://console.cloud.google.com/storage?project=$PROJECT_ID"
echo ""
