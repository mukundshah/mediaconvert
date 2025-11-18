package analytics

import (
	"context"
	"fmt"
	"time"
)

// JobStats represents aggregated job statistics
type JobStats struct {
	TotalJobs        int64   `json:"total_jobs"`
	CompletedJobs    int64   `json:"completed_jobs"`
	FailedJobs       int64   `json:"failed_jobs"`
	PendingJobs      int64   `json:"pending_jobs"`
	ProcessingJobs   int64   `json:"processing_jobs"`
	SuccessRate      float64 `json:"success_rate"`
	AvgProcessingTime float64 `json:"avg_processing_time_ms"`
	TotalDataProcessed int64  `json:"total_data_processed_bytes"`
}

// TimelinePoint represents a data point in a timeline
type TimelinePoint struct {
	Time              time.Time `json:"time"`
	JobCount          int64     `json:"job_count"`
	AvgProcessingTime float64   `json:"avg_processing_time_ms"`
	SuccessRate       float64   `json:"success_rate"`
}

// PipelineStat represents statistics for a pipeline
type PipelineStat struct {
	PipelineID       *uint64  `json:"pipeline_id"`
	TotalJobs        int64    `json:"total_jobs"`
	SuccessfulJobs   int64    `json:"successful_jobs"`
	FailedJobs       int64    `json:"failed_jobs"`
	SuccessRate      float64  `json:"success_rate"`
	AvgProcessingTime float64 `json:"avg_processing_time_ms"`
}

// GetJobStats returns aggregated job statistics for a user
func (c *Client) GetJobStats(ctx context.Context, userID uint64, days int) (*JobStats, error) {
	query := `
		SELECT
			count() as total_jobs,
			sumIf(1, status = 'completed') as completed_jobs,
			sumIf(1, status = 'failed') as failed_jobs,
			sumIf(1, status = 'pending') as pending_jobs,
			sumIf(1, status = 'processing') as processing_jobs,
			(completed_jobs * 100.0 / total_jobs) as success_rate,
			avgIf(processing_time_ms, status = 'completed') as avg_processing_time,
			sumIf(file_size, status = 'completed') as total_data_processed
		FROM job_metrics
		WHERE user_id = ? AND timestamp >= now() - INTERVAL ? DAY
	`

	var stats JobStats
	row := c.conn.QueryRow(ctx, query, userID, days)

	if err := row.Scan(
		&stats.TotalJobs,
		&stats.CompletedJobs,
		&stats.FailedJobs,
		&stats.PendingJobs,
		&stats.ProcessingJobs,
		&stats.SuccessRate,
		&stats.AvgProcessingTime,
		&stats.TotalDataProcessed,
	); err != nil {
		return nil, fmt.Errorf("failed to scan stats: %w", err)
	}

	return &stats, nil
}

// GetJobTimeline returns job metrics over time
func (c *Client) GetJobTimeline(ctx context.Context, userID uint64, days int, interval string) ([]TimelinePoint, error) {
	var timeFunc string
	switch interval {
	case "day":
		timeFunc = "toStartOfDay"
	case "hour":
		timeFunc = "toStartOfHour"
	default:
		timeFunc = "toStartOfHour"
	}

	query := fmt.Sprintf(`
		SELECT
			%s(timestamp) as time,
			count() as job_count,
			avgIf(processing_time_ms, status = 'completed') as avg_processing_time,
			(sumIf(1, status = 'completed') * 100.0 / count()) as success_rate
		FROM job_metrics
		WHERE user_id = ? AND timestamp >= now() - INTERVAL ? DAY
		GROUP BY time
		ORDER BY time
	`, timeFunc)

	rows, err := c.conn.Query(ctx, query, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to query timeline: %w", err)
	}
	defer rows.Close()

	var timeline []TimelinePoint
	for rows.Next() {
		var point TimelinePoint
		if err := rows.Scan(
			&point.Time,
			&point.JobCount,
			&point.AvgProcessingTime,
			&point.SuccessRate,
		); err != nil {
			return nil, fmt.Errorf("failed to scan timeline point: %w", err)
		}
		timeline = append(timeline, point)
	}

	return timeline, nil
}

// GetPipelineStats returns statistics grouped by pipeline
func (c *Client) GetPipelineStats(ctx context.Context, userID uint64, days int) ([]PipelineStat, error) {
	query := `
		SELECT
			pipeline_id,
			count() as total_jobs,
			sumIf(1, status = 'completed') as successful_jobs,
			sumIf(1, status = 'failed') as failed_jobs,
			(successful_jobs * 100.0 / total_jobs) as success_rate,
			avgIf(processing_time_ms, status = 'completed') as avg_processing_time
		FROM job_metrics
		WHERE user_id = ? AND timestamp >= now() - INTERVAL ? DAY
		GROUP BY pipeline_id
		ORDER BY total_jobs DESC
	`

	rows, err := c.conn.Query(ctx, query, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to query pipeline stats: %w", err)
	}
	defer rows.Close()

	var stats []PipelineStat
	for rows.Next() {
		var stat PipelineStat
		if err := rows.Scan(
			&stat.PipelineID,
			&stat.TotalJobs,
			&stat.SuccessfulJobs,
			&stat.FailedJobs,
			&stat.SuccessRate,
			&stat.AvgProcessingTime,
		); err != nil {
			return nil, fmt.Errorf("failed to scan pipeline stat: %w", err)
		}
		stats = append(stats, stat)
	}

	return stats, nil
}
