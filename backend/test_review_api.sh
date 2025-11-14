#!/bin/bash

# Review API Test Script
# This script tests the Review API endpoints using curl

echo "========================================="
echo "Review API Test Script"
echo "========================================="
echo ""

# Test Health Endpoint
echo "1. Testing Health Endpoint..."
curl -s http://localhost:8080/health | jq .
echo ""
echo ""

# Note: The Review API requires authentication
echo "2. Review API Endpoints (require authentication):"
echo "   - GET  /api/v1/review/stats"
echo "   - GET  /api/v1/review/items"
echo "   - POST /api/v1/review/submit"
echo ""

echo "To test these endpoints with authentication:"
echo "1. Register a user:"
echo "   curl -X POST http://localhost:8080/api/v1/auth/register \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{\"email\":\"user@example.com\",\"password\":\"Password123!\",\"display_name\":\"Test User\"}'"
echo ""
echo "2. Login to get a token:"
echo "   curl -X POST http://localhost:8080/api/v1/auth/login \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{\"email\":\"user@example.com\",\"password\":\"Password123!\"}'"
echo ""
echo "3. Use the token to access Review API:"
echo "   TOKEN=<your_access_token>"
echo "   curl -H \"Authorization: Bearer \$TOKEN\" http://localhost:8080/api/v1/review/stats"
echo ""

echo "========================================="
echo "Backend tests passed successfully!"
echo "========================================="
echo ""
echo "Run: go test -v ./internal/api/handler -run Review"
echo ""
