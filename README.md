# MediaConvert

A scalable media processing service with S3-compatible API, pipeline-based processing, and asynchronous job execution.

## Features

- **S3-Compatible API**: Full AWS S3 API compatibility for seamless integration with existing tools
- **Pipeline-Based Processing**: Define reusable processing pipelines in YAML or JSON
- **Asynchronous Processing**: Background worker service with Redis-based job notifications
- **JWT Authentication**: Secure API access with JWT tokens
- **Multiple Media Operations**: Support for video transcoding, image resizing, PDF text extraction, thumbnail generation, and more
- **Custom S3 Credentials**: Per-user S3 credentials with bucket isolation
- **Job Management**: Track, cancel, and rerun processing jobs
- **Status History**: Complete audit trail of job status changes
- **Analytics & Metrics**: Real-time job analytics powered by ClickHouse for performance monitoring and insights

## Architecture

```mermaid
flowchart LR
    Client["Client<br/>(AWS CLI)"] -->|"HTTP Requests"| APIServer["API Server<br/>(Gin)"]
    APIServer -->|"Job Queue"| Worker["Worker<br/>(Processor)"]
    APIServer -->|"Data Storage"| PostgreSQL["PostgreSQL"]
    APIServer -->|"Job Notifications"| Redis["Redis"]
    APIServer -->|"Analytics Queries"| ClickHouse["ClickHouse<br/>(Analytics)"]
    Worker -->|"Read/Write Files"| MinIO["MinIO<br/>(S3)"]
    Worker -->|"Metrics & Logs"| ClickHouse
```

## Prerequisites

- **Go 1.24+**
- **PostgreSQL 15+**
- **Redis 7+**
- **ClickHouse** (optional, for analytics)
- **MinIO** (or compatible S3 storage)
- **System Tools**:
  - `ffmpeg` (for video/audio processing)
  - `imagemagick` (for image processing)
  - `pdftotext` (for PDF text extraction)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/mukund/mediaconvert.git
cd mediaconvert
```

### 2. Start Infrastructure

```bash
# Start PostgreSQL, Redis, and MinIO
make docker-up
```

This starts:

- PostgreSQL on port `5432`
- Redis on port `6379`
- ClickHouse on ports `8123` (HTTP) and `9001` (Native protocol)
- MinIO on ports `9000` (API) and `9001` (Console)

### 3. Configure Environment

Create a `.env` file:

```env
DATABASE_URL=postgres://user:password@localhost:5432/mediaconvert?sslmode=disable
REDIS_URL=localhost:6379
CLICKHOUSE_DSN=localhost:9001
PORT=8080
S3_ENDPOINT=http://localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=media
S3_REGION=us-east-1
JWT_SECRET=your-secret-key-here
```

> **Note**: `CLICKHOUSE_DSN` is optional. If not provided or if ClickHouse is unavailable, the service will continue to work without analytics.

### 4. Run the Services

**Terminal 1 - API Server:**

```bash
make run
```

**Terminal 2 - Worker:**

```bash
make run-worker
```

### 5. Test the Setup

```bash
# Run automated tests
./test/test-api.sh
```

## API Overview

### Authentication

All API endpoints (except `/auth/*` and `/health`) require JWT authentication.

#### Register

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

#### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

### S3 Credentials

#### Create Credentials

```bash
curl -X POST http://localhost:8080/api/s3-credentials \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### List Credentials

```bash
curl -X GET http://localhost:8080/api/s3-credentials \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Pipeline Management

#### Create Pipeline

```bash
curl -X POST http://localhost:8080/api/pipelines \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "video-compress",
    "format": "yaml",
    "content": "name: video-compress\nsteps:\n  - operation: transcode\n    input: \${input}\n    output: \${output}/output.mp4\n    params:\n      codec: h264\n      quality: 23"
  }'
```

### Job Management

#### List Jobs

```bash
curl -X GET http://localhost:8080/api/jobs \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Job Details

```bash
curl -X GET http://localhost:8080/api/jobs/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Analytics

The service provides analytics endpoints powered by ClickHouse for monitoring job performance and usage patterns.

#### Get Job Statistics

```bash
curl -X GET "http://localhost:8080/api/analytics/jobs/stats?days=7" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Returns aggregated statistics including:

- Total, completed, failed, pending, and processing jobs
- Success rate
- Average processing time
- Total data processed

#### Get Job Timeline

```bash
curl -X GET "http://localhost:8080/api/analytics/jobs/timeline?days=7&interval=hour" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Returns time-series data showing job metrics over time. Supports `hour` and `day` intervals.

#### Get Pipeline Statistics

```bash
curl -X GET "http://localhost:8080/api/analytics/pipelines/stats?days=30" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Returns statistics grouped by pipeline, including:

- Total jobs per pipeline
- Success/failure counts
- Success rate
- Average processing time

## S3-Compatible API

The service provides a fully S3-compatible API that works with standard AWS tools.

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

The `Pipeline` metadata automatically creates a processing job!

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

### Using Python (boto3)

```python
import boto3

s3 = boto3.client(
    's3',
    endpoint_url='http://localhost:8080',
    aws_access_key_id='YOUR_ACCESS_KEY',
    aws_secret_access_key='YOUR_SECRET_KEY',
    region_name='us-east-1'
)

# Upload with pipeline metadata
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
```

## Pipeline Operations

### Supported Operations

- **`transcode`**: Video/audio transcoding (H264, H265, VP9, etc.)
- **`resize`**: Image resizing
- **`extract_text`**: Extract text from PDFs
- **`extract_frame`**: Extract frames from videos
- **`convert`**: Format conversion
- **`generate_thumbnail`**: Generate thumbnails from videos, images, or PDFs

### Pipeline Example

```yaml
name: "video-process-complete"
description: "Complete video processing: compress + thumbnail"
steps:
  - operation: "transcode"
    input: "${input}"
    output: "${output}/video.mp4"
    params:
      codec: "h264"
      quality: 23
      audio_codec: "aac"

  - operation: "generate_thumbnail"
    input: "${input}"
    output: "${output}/thumb.jpg"
    params:
      type: "video"
      timestamp: "00:00:01"
      width: 320
      height: 240
```

### Variables

- `${input}`: Path to the input file
- `${output}`: Path to the output directory

## Example Pipelines

See `test/pipelines/` for example pipeline definitions:

- `video-compress.yaml` - Video compression
- `video-thumbnail.yaml` - Thumbnail generation
- `image-resize.yaml` - Image resizing
- `video-complete.yaml` - Multi-step processing
- `pdf-extract.yaml` - PDF text extraction + thumbnail

## Project Structure

```
mediaconvert/
├── cmd/
│   ├── server/          # API server
│   └── worker/          # Background worker
├── internal/
│   ├── analytics/       # ClickHouse analytics client
│   ├── auth/            # JWT authentication
│   ├── config/          # Configuration management
│   ├── db/              # Database connection & migrations
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   ├── pipeline/        # Pipeline schema & parsing
│   ├── s3compat/        # S3-compatible API implementation
│   ├── system/          # System dependency checks
│   └── worker/          # Job processing logic
├── test/
│   ├── pipelines/       # Example pipelines
│   ├── test-api.sh      # Automated test script
│   └── TESTING.md       # Testing guide
└── docker-compose.yml   # Infrastructure setup
```

## Development

### Build

```bash
# Build server
make build

# Build worker
make build-worker
```

### Run Tests

```bash
make test
```

### Clean Build Artifacts

```bash
make clean
```

## Configuration

All configuration is done via environment variables or a `.env` file:

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://user:password@localhost:5432/mediaconvert?sslmode=disable` | PostgreSQL connection string |
| `REDIS_URL` | `localhost:6379` | Redis connection string |
| `CLICKHOUSE_DSN` | `localhost:9001` | ClickHouse connection string (optional) |
| `PORT` | `8080` | API server port |
| `S3_ENDPOINT` | `http://localhost:9000` | S3-compatible storage endpoint |
| `S3_ACCESS_KEY` | `minioadmin` | S3 access key |
| `S3_SECRET_KEY` | `minioadmin` | S3 secret key |
| `S3_BUCKET` | `media` | Default S3 bucket |
| `S3_REGION` | `us-east-1` | S3 region |
| `JWT_SECRET` | `change-this-secret-in-production` | JWT signing secret |

## Analytics

The service automatically records job metrics to ClickHouse when available:

- **Job Metrics**: Processing times, status, file sizes, content types
- **Status Transitions**: Complete audit trail of job status changes
- **Pipeline Execution Logs**: Detailed operation-level metrics (optional)

Analytics are **optional** - if ClickHouse is not configured or unavailable, the service continues to function normally without analytics.

### Analytics Tables

- `job_metrics`: Aggregated job execution metrics
- `job_status_history_analytics`: Time-series job status transitions
- `pipeline_execution_logs`: Detailed pipeline operation logs

## Security Considerations

- **JWT Secret**: Change the default JWT secret in production
- **S3 Credentials**: Secret keys are stored in plain text (required for signature validation)
- **Bucket Isolation**: Users can only access their own buckets
- **Authentication**: All API endpoints require valid JWT tokens

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests. Whether it's bug fixes, new features, documentation improvements, or examples, your contributions help make this project better.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). See the [LICENSE](LICENSE) file for details.
