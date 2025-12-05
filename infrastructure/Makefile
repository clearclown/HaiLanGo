# HaiLanGo Infrastructure Makefile

.PHONY: help setup-gcloud terraform-init terraform-plan terraform-apply terraform-destroy check-costs create-budget clean

# Default target
help:
	@echo "HaiLanGo Infrastructure Management"
	@echo ""
	@echo "Google Cloud Setup:"
	@echo "  make setup-gcloud     - Initial GCP setup (login, APIs, service account)"
	@echo "  make create-budget    - Create monthly budget with alerts"
	@echo "  make check-costs      - Display current costs and usage"
	@echo ""
	@echo "Terraform:"
	@echo "  make terraform-init   - Initialize Terraform"
	@echo "  make terraform-plan   - Show Terraform plan"
	@echo "  make terraform-apply  - Apply Terraform changes"
	@echo "  make terraform-destroy - Destroy all resources (DANGER!)"
	@echo "  make terraform-output - Show Terraform outputs"
	@echo ""
	@echo "Utilities:"
	@echo "  make env-file         - Generate .env file from Terraform outputs"
	@echo "  make clean            - Clean temporary files"
	@echo ""

# Google Cloud Setup
setup-gcloud:
	@echo "Setting up Google Cloud..."
	./scripts/setup-gcloud.sh

create-budget:
	@echo "Creating budget..."
	./scripts/create-budget.sh

check-costs:
	@echo "Checking costs..."
	./scripts/check-costs.sh

# Terraform commands
terraform-init:
	@echo "Initializing Terraform..."
	cd terraform && terraform init

terraform-plan:
	@echo "Planning Terraform changes..."
	cd terraform && terraform plan

terraform-apply:
	@echo "Applying Terraform changes..."
	cd terraform && terraform apply

terraform-destroy:
	@echo "WARNING: This will destroy all resources!"
	@read -p "Are you sure? (yes/no): " confirm && [ "$$confirm" = "yes" ]
	cd terraform && terraform destroy

terraform-output:
	@echo "Terraform outputs:"
	cd terraform && terraform output

# Generate .env file
env-file:
	@echo "Generating .env file..."
	cd terraform && terraform output -raw env_file_content > ../../backend/.env.gcloud
	@echo "Generated: backend/.env.gcloud"
	@echo "Merge this with your existing .env file"

# Get service account credentials
get-credentials:
	@echo "Extracting service account credentials..."
	cd terraform && terraform output -raw service_account_key_json | base64 -d > ../credentials.json
	@echo "Saved to: infrastructure/credentials.json"
	@echo "WARNING: Keep this file secure!"

# Clean temporary files
clean:
	rm -f terraform/.terraform.lock.hcl
	rm -f terraform/terraform.tfstate.backup
	rm -rf terraform/.terraform/
	rm -f credentials.json
	@echo "Cleaned temporary files"

# Quick start
quickstart: setup-gcloud terraform-init terraform-apply env-file
	@echo ""
	@echo "=================================="
	@echo "Quick Start Complete!"
	@echo "=================================="
	@echo ""
	@echo "Next steps:"
	@echo "1. Copy env variables: cat ../backend/.env.gcloud >> ../backend/.env"
	@echo "2. Set credentials: export GOOGLE_APPLICATION_CREDENTIALS=./credentials.json"
	@echo "3. Start the backend: cd ../backend && go run cmd/server/main.go"
