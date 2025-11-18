package analytics

import (
	"context"
	"fmt"
	"log"
)

// InitSchema creates all necessary tables in ClickHouse
func (c *Client) InitSchema(ctx context.Context) error {
	tables := []string{
		createJobMetricsTable,
		createJobStatusHistoryTable,
		createPipelineExecutionLogsTable,
	}

	for _, tableSQL := range tables {
		if err := c.conn.Exec(ctx, tableSQL); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	log.Println("ClickHouse schema initialized successfully")
	return nil
}

const createJobMetricsTable = `
CREATE TABLE IF NOT EXISTS job_metrics (
    timestamp DateTime64(3),
    job_id UInt64,
    user_id UInt64,
    pipeline_id Nullable(UInt64),
    file_id UInt64,
    content_type String,
    file_size UInt64,
    status String,
    processing_time_ms UInt64,
    created_at DateTime,
    finished_at Nullable(DateTime),
    error_message Nullable(String)
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, job_id)
SETTINGS index_granularity = 8192;
`

const createJobStatusHistoryTable = `
CREATE TABLE IF NOT EXISTS job_status_history_analytics (
    timestamp DateTime64(3),
    job_id UInt64,
    user_id UInt64,
    from_status String,
    to_status String,
    triggered_by String,
    message String
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, job_id)
SETTINGS index_granularity = 8192;
`

const createPipelineExecutionLogsTable = `
CREATE TABLE IF NOT EXISTS pipeline_execution_logs (
    timestamp DateTime64(3),
    job_id UInt64,
    user_id UInt64,
    pipeline_id Nullable(UInt64),
    operation_name String,
    operation_type String,
    duration_ms UInt64,
    input_size_bytes UInt64,
    output_size_bytes UInt64,
    success UInt8,
    error_message Nullable(String)
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, job_id, operation_name)
SETTINGS index_granularity = 8192;
`
