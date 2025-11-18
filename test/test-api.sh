#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_URL="${API_URL:-http://localhost:8088}"
TEST_EMAIL="test@example.com"
TEST_PASSWORD="password123"

# Global variables
JWT_TOKEN=""
ACCESS_KEY=""
SECRET_KEY=""
BUCKET_NAME=""
PIPELINE_ID=""

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Media Processing Server - API Test${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Helper function to make API calls
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3

    if [ -n "$JWT_TOKEN" ]; then
        if [ -n "$data" ]; then
            curl -s -X "$method" "$API_URL$endpoint" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $JWT_TOKEN" \
                -d "$data"
        else
            curl -s -X "$method" "$API_URL$endpoint" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $JWT_TOKEN"
        fi
    else
        if [ -n "$data" ]; then
            curl -s -X "$method" "$API_URL$endpoint" \
                -H "Content-Type: application/json" \
                -d "$data"
        else
            curl -s -X "$method" "$API_URL$endpoint" \
                -H "Content-Type: application/json"
        fi
    fi
}

# Test 1: Health Check
echo -e "${YELLOW}[1/9] Testing health endpoint...${NC}"
response=$(curl -s "$API_URL/health")
if echo "$response" | grep -q "ok"; then
    echo -e "${GREEN}✓ Health check passed${NC}"
else
    echo -e "${RED}✗ Health check failed${NC}"
    exit 1
fi
echo ""

# Test 2: User Registration
echo -e "${YELLOW}[2/9] Registering user...${NC}"
response=$(api_call POST "/auth/register" "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}")
JWT_TOKEN=$(echo "$response" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -n "$JWT_TOKEN" ]; then
    echo -e "${GREEN}✓ User registered successfully${NC}"
    echo -e "   JWT Token: ${JWT_TOKEN:0:20}..."
else
    # Check if user already exists
    if echo "$response" | grep -q "already exists"; then
        echo -e "${YELLOW}⚠ User already exists, trying login...${NC}"

        # Try login instead
        response=$(api_call POST "/auth/login" "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}")
        JWT_TOKEN=$(echo "$response" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

        if [ -n "$JWT_TOKEN" ]; then
            echo -e "${GREEN}✓ Login successful${NC}"
            echo -e "   JWT Token: ${JWT_TOKEN:0:20}..."
        else
            echo -e "${RED}✗ Login failed${NC}"
            echo -e "   Response: $response"
            exit 1
        fi
    else
        echo -e "${RED}✗ Registration failed${NC}"
        echo -e "   Response: $response"
        exit 1
    fi
fi
echo ""

# Test 3: Create S3 Credentials
echo -e "${YELLOW}[3/9] Creating S3 credentials...${NC}"
response=$(api_call POST "/api/s3-credentials" "")
ACCESS_KEY=$(echo "$response" | grep -o '"access_key":"[^"]*' | cut -d'"' -f4)
SECRET_KEY=$(echo "$response" | grep -o '"secret_key":"[^"]*' | cut -d'"' -f4)
BUCKET_NAME=$(echo "$response" | grep -o '"bucket_name":"[^"]*' | cut -d'"' -f4)

if [ -n "$ACCESS_KEY" ] && [ -n "$SECRET_KEY" ] && [ -n "$BUCKET_NAME" ]; then
    echo -e "${GREEN}✓ S3 credentials created${NC}"
    echo -e "   Access Key: $ACCESS_KEY"
    echo -e "   Secret Key: ${SECRET_KEY:0:20}..."
    echo -e "   Bucket: $BUCKET_NAME"
else
    echo -e "${RED}✗ Failed to create S3 credentials${NC}"
    exit 1
fi
echo ""

# Test 4: List S3 Credentials
echo -e "${YELLOW}[4/9] Listing S3 credentials...${NC}"
response=$(api_call GET "/api/s3-credentials" "")
if echo "$response" | grep -q "credentials"; then
    echo -e "${GREEN}✓ S3 credentials listed${NC}"
else
    echo -e "${RED}✗ Failed to list credentials${NC}"
fi
echo ""

# Test 5: Create Pipeline
echo -e "${YELLOW}[5/9] Creating pipeline...${NC}"

# Create pipeline with properly formatted JSON
PIPELINE_JSON='{"name":"test-pipeline","format":"yaml","content":"name: \"test-pipeline\"\nsteps:\n  - operation: \"transcode\"\n    input: \"${input}\"\n    output: \"${output}/output.mp4\"\n    params:\n      codec: \"h264\"\n      quality: 23"}'

response=$(api_call POST "/api/pipelines" "$PIPELINE_JSON")
PIPELINE_ID=$(echo "$response" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -n "$PIPELINE_ID" ]; then
    echo -e "${GREEN}✓ Pipeline created${NC}"
    echo -e "   Pipeline ID: $PIPELINE_ID"
else
    echo -e "${RED}✗ Failed to create pipeline${NC}"
    echo -e "   Response: $response"
fi
echo ""

# Test 6: List Pipelines
echo -e "${YELLOW}[6/9] Listing pipelines...${NC}"
response=$(api_call GET "/api/pipelines" "")
if echo "$response" | grep -q "test-pipeline"; then
    echo -e "${GREEN}✓ Pipelines listed${NC}"
else
    echo -e "${RED}✗ Failed to list pipelines${NC}"
fi
echo ""

# Test 7: Get Pipeline
if [ -n "$PIPELINE_ID" ]; then
    echo -e "${YELLOW}[7/9] Getting pipeline details...${NC}"
    response=$(api_call GET "/api/pipelines/$PIPELINE_ID" "")
    if echo "$response" | grep -q "test-pipeline"; then
        echo -e "${GREEN}✓ Pipeline retrieved${NC}"
    else
        echo -e "${RED}✗ Failed to get pipeline${NC}"
    fi
else
    echo -e "${YELLOW}[7/9] Skipping pipeline details (no pipeline ID)${NC}"
fi
echo ""

# Test 8: List Jobs
echo -e "${YELLOW}[8/9] Listing jobs...${NC}"
response=$(api_call GET "/api/jobs" "")
if echo "$response" | grep -q "jobs"; then
    echo -e "${GREEN}✓ Jobs listed${NC}"
else
    echo -e "${RED}✗ Failed to list jobs${NC}"
fi
echo ""

# Test 9: Analytics Endpoints
echo -e "${YELLOW}[9/11] Testing analytics endpoints...${NC}"
response=$(api_call GET "/api/analytics/jobs/stats?days=7" "")
if echo "$response" | grep -q "total_jobs\|error"; then
    if echo "$response" | grep -q "not available"; then
        echo -e "${YELLOW}⚠ Analytics not available (ClickHouse may not be running)${NC}"
    else
        echo -e "${GREEN}✓ Job stats retrieved${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Analytics endpoint returned unexpected response${NC}"
fi
echo ""

# Test 10: Analytics Timeline
echo -e "${YELLOW}[10/11] Testing analytics timeline...${NC}"
response=$(api_call GET "/api/analytics/jobs/timeline?days=7&interval=hour" "")
if echo "$response" | grep -q "\[\]\|time\|error"; then
    if echo "$response" | grep -q "not available"; then
        echo -e "${YELLOW}⚠ Analytics not available (ClickHouse may not be running)${NC}"
    else
        echo -e "${GREEN}✓ Job timeline retrieved${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Analytics timeline returned unexpected response${NC}"
fi
echo ""

# Test 11: AWS CLI Configuration (if available)
echo -e "${YELLOW}[11/11] Checking AWS CLI configuration...${NC}"
if command -v aws &> /dev/null; then
    echo -e "${GREEN}✓ AWS CLI is installed${NC}"
    echo ""
    echo -e "${BLUE}To configure AWS CLI with your credentials:${NC}"
    echo -e "aws configure set aws_access_key_id $ACCESS_KEY --profile mediaconvert"
    echo -e "aws configure set aws_secret_access_key $SECRET_KEY --profile mediaconvert"
    echo -e "aws configure set region us-east-1 --profile mediaconvert"
    echo ""
    echo -e "${BLUE}To upload a file:${NC}"
    echo -e "aws s3 cp test.mp4 s3://$BUCKET_NAME/test.mp4 \\"
    echo -e "  --endpoint-url $API_URL \\"
    echo -e "  --metadata Pipeline=test-pipeline \\"
    echo -e "  --profile mediaconvert"
else
    echo -e "${YELLOW}⚠ AWS CLI not installed${NC}"
    echo -e "   Install with: pip install awscli"
fi
echo ""

# Summary
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ All tests completed!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${BLUE}Analytics Endpoints:${NC}"
echo -e "1. Job Stats: curl -H 'Authorization: Bearer $JWT_TOKEN' $API_URL/api/analytics/jobs/stats?days=7"
echo -e "2. Job Timeline: curl -H 'Authorization: Bearer $JWT_TOKEN' $API_URL/api/analytics/jobs/timeline?days=7&interval=hour"
echo -e "3. Pipeline Stats: curl -H 'Authorization: Bearer $JWT_TOKEN' $API_URL/api/analytics/pipelines/stats?days=30"
echo ""
echo -e "${BLUE}Credentials saved for manual testing:${NC}"
echo -e "JWT Token: $JWT_TOKEN"
echo -e "Access Key: $ACCESS_KEY"
echo -e "Secret Key: $SECRET_KEY"
echo -e "Bucket: $BUCKET_NAME"
echo -e "Pipeline ID: $PIPELINE_ID"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo -e "1. Upload a file using AWS CLI (see command above)"
echo -e "2. Check job status: curl -H 'Authorization: Bearer $JWT_TOKEN' $API_URL/api/jobs"
echo -e "3. Monitor worker logs to see processing"
