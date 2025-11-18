package analytics

import (
	"context"
	"fmt"
	"time"
)

// RecordJobMetric records a job metric to ClickHouse
func (c *Client) RecordJobMetric(ctx context.Context, metric JobMetric) error {
	var pipelineID *uint64 = metric.PipelineID
	var finishedAt *time.Time = metric.FinishedAt
	var errorMsg *string = metric.ErrorMessage

	batch, err := c.conn.PrepareBatch(ctx, "INSERT INTO job_metrics")
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	if err := batch.Append(
		metric.Timestamp,
		metric.JobID,
		metric.UserID,
		pipelineID,
		metric.FileID,
		metric.ContentType,
		uint64(metric.FileSize),
		metric.Status,
		uint64(metric.ProcessingTime.Milliseconds()),
		metric.CreatedAt,
		finishedAt,
		errorMsg,
	); err != nil {
		return fmt.Errorf("failed to append to batch: %w", err)
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send batch: %w", err)
	}

	return nil
}

// RecordJobStatusTransition records a job status change to ClickHouse
func (c *Client) RecordJobStatusTransition(ctx context.Context, transition JobStatusTransition) error {
	batch, err := c.conn.PrepareBatch(ctx, "INSERT INTO job_status_history_analytics")
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	if err := batch.Append(
		transition.Timestamp,
		transition.JobID,
		transition.UserID,
		transition.FromStatus,
		transition.ToStatus,
		transition.TriggeredBy,
		transition.Message,
	); err != nil {
		return fmt.Errorf("failed to append to batch: %w", err)
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send batch: %w", err)
	}

	return nil
}

// RecordPipelineExecutionLog records a pipeline operation execution log
func (c *Client) RecordPipelineExecutionLog(ctx context.Context, log PipelineExecutionLog) error {
	var pipelineID *uint64 = log.PipelineID
	var errorMsg *string = log.ErrorMessage
	var success uint8 = 0
	if log.Success {
		success = 1
	}

	batch, err := c.conn.PrepareBatch(ctx, "INSERT INTO pipeline_execution_logs")
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	if err := batch.Append(
		log.Timestamp,
		log.JobID,
		log.UserID,
		pipelineID,
		log.OperationName,
		log.OperationType,
		uint64(log.Duration.Milliseconds()),
		uint64(log.InputSize),
		uint64(log.OutputSize),
		success,
		errorMsg,
	); err != nil {
		return fmt.Errorf("failed to append to batch: %w", err)
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send batch: %w", err)
	}

	return nil
}
