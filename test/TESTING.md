# Testing Guide

## Quick Start

### 1. Start the System

```bash
# Start infrastructure (Postgres, Redis, MinIO)
make docker-up

# Terminal 1: Start API server
make run

# Terminal 2: Start worker
make run-worker
```

### 2. Run Automated Tests

```bash
# Run the test script
./test/test-api.sh

# The script will:
# - Register a test user
# - Create S3 credentials
# - Create a test pipeline
# - List all resources
# - Provide AWS CLI commands for file upload
```

## Manual Testing

### Authentication

#### Register User

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com"
  }
}
```

#### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Save the JWT token for subsequent requests!**

### S3 Credentials

#### Create Credentials

```bash
curl -X POST http://localhost:8080/api/s3-credentials \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response:

```json
{
  "id": 1,
  "access_key": "AKIAIOSFODNN7EXAMPLE",
  "secret_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
  "bucket_name": "user-1-abc12345",
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Important**: Secret key is only shown once! Save it.

#### List Credentials

```bash
curl -X GET http://localhost:8080/api/s3-credentials \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Pipeline Management

#### Create Pipeline

```bash
# Using YAML
curl -X POST http://localhost:8080/api/pipelines \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "video-compress",
    "format": "yaml",
    "content": "name: \"video-compress\"\nsteps:\n  - operation: \"transcode\"\n    input: \"${input}\"\n    output: \"${output}/output.mp4\"\n    params:\n      codec: \"h264\"\n      quality: 23"
  }'
```

Or use the example pipelines:

```bash
# Read from file
PIPELINE_CONTENT=$(cat test/pipelines/video-compress.yaml)
curl -X POST http://localhost:8080/api/pipelines \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"video-compress\",
    \"format\": \"yaml\",
    \"content\": \"$PIPELINE_CONTENT\"
  }"
```

#### List Pipelines

```bash
curl -X GET http://localhost:8080/api/pipelines \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Pipeline

```bash
curl -X GET http://localhost:8080/api/pipelines/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Update Pipeline

```bash
curl -X PUT http://localhost:8080/api/pipelines/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "video-compress-updated",
    "format": "yaml",
    "content": "..."
  }'
```

#### Delete Pipeline

```bash
curl -X DELETE http://localhost:8080/api/pipelines/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Job Management

#### List Jobs

```bash
# All jobs
curl -X GET http://localhost:8080/api/jobs \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Filter by status
curl -X GET "http://localhost:8080/api/jobs?status=completed" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Pagination
curl -X GET "http://localhost:8080/api/jobs?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Job Details

```bash
curl -X GET http://localhost:8080/api/jobs/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Cancel Job

```bash
curl -X POST http://localhost:8080/api/jobs/1/cancel \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Rerun Job

```bash
curl -X POST http://localhost:8080/api/jobs/1/rerun \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## S3 File Upload

### Using AWS CLI

#### Configure AWS CLI

```bash
aws configure set aws_access_key_id YOUR_ACCESS_KEY --profile mediaconvert
aws configure set aws_secret_access_key YOUR_SECRET_KEY --profile mediaconvert
aws configure set region us-east-1 --profile mediaconvert
```

#### Upload File (with automatic job creation)

```bash
aws s3 cp video.mp4 s3://YOUR_BUCKET_NAME/video.mp4 \
  --endpoint-url http://localhost:8080 \
  --metadata Pipeline=video-compress \
  --profile mediaconvert
```

The `Pipeline` metadata triggers automatic job creation!

#### List Files

```bash
aws s3 ls s3://YOUR_BUCKET_NAME/ \
  --endpoint-url http://localhost:8080 \
  --profile mediaconvert
```

#### Download File

```bash
aws s3 cp s3://YOUR_BUCKET_NAME/video.mp4 downloaded.mp4 \
  --endpoint-url http://localhost:8080 \
  --profile mediaconvert
```

#### Delete File

```bash
aws s3 rm s3://YOUR_BUCKET_NAME/video.mp4 \
  --endpoint-url http://localhost:8080 \
  --profile mediaconvert
```

### Using Python (boto3)

```python
import boto3

# Configure S3 client
s3 = boto3.client(
    's3',
    endpoint_url='http://localhost:8080',
    aws_access_key_id='YOUR_ACCESS_KEY',
    aws_secret_access_key='YOUR_SECRET_KEY',
    region_name='us-east-1'
)

# Upload file with pipeline metadata
s3.upload_file(
    'video.mp4',
    'YOUR_BUCKET_NAME',
    'video.mp4',
    ExtraArgs={
        'Metadata': {
            'Pipeline': 'video-compress'
        }
    }
)

# List files
response = s3.list_objects_v2(Bucket='YOUR_BUCKET_NAME')
for obj in response.get('Contents', []):
    print(obj['Key'])
```

## Testing Workflow

### Complete End-to-End Test

1. **Setup**

   ```bash
   ./test/test-api.sh
   ```

   Save the credentials output.

2. **Create Pipeline**

   ```bash
   curl -X POST http://localhost:8080/api/pipelines \
     -H "Authorization: Bearer YOUR_JWT" \
     -H "Content-Type: application/json" \
     -d @test/examples/create-pipeline.json
   ```

3. **Upload Test File**

   ```bash
   # Create a test video (requires ffmpeg)
   ffmpeg -f lavfi -i testsrc=duration=10:size=1280x720:rate=30 \
     -f lavfi -i sine=frequency=1000:duration=10 \
     -pix_fmt yuv420p test-video.mp4

   # Upload with pipeline
   aws s3 cp test-video.mp4 s3://YOUR_BUCKET/test.mp4 \
     --endpoint-url http://localhost:8080 \
     --metadata Pipeline=video-compress \
     --profile mediaconvert
   ```

4. **Monitor Job**

   ```bash
   # Watch job status
   watch -n 2 'curl -s -H "Authorization: Bearer YOUR_JWT" \
     http://localhost:8080/api/jobs | jq'
   ```

5. **Check Results**

   ```bash
   # List processed files
   aws s3 ls s3://YOUR_BUCKET/results/ --recursive \
     --endpoint-url http://localhost:8080 \
     --profile mediaconvert
   ```

## Troubleshooting

### Check Logs

```bash
# API server logs
# (visible in terminal 1)

# Worker logs
# (visible in terminal 2)

# Database
docker-compose logs postgres

# Redis
docker-compose logs redis

# MinIO
docker-compose logs minio
```

### Common Issues

1. **"Failed to connect to Redis"**
   - Ensure `make docker-up` was run
   - Check Redis is running: `docker-compose ps`

2. **"Signature mismatch"**
   - Verify access key and secret key
   - Check endpoint URL format (must include http://)

3. **"Pipeline not found"**
   - Ensure pipeline name matches exactly
   - Check pipeline exists: `GET /api/pipelines`

4. **Job stuck in "pending"**
   - Check worker is running
   - Check worker logs for errors
   - Verify Redis connection

## Example Pipelines

See `test/pipelines/` for example pipeline definitions:

- `video-compress.yaml` - Video compression
- `video-thumbnail.yaml` - Thumbnail generation
- `image-resize.yaml` - Image resizing
- `video-complete.yaml` - Multi-step processing
- `pdf-extract.yaml` - PDF text extraction + thumbnail
