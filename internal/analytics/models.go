package analytics

import (
	"time"
)

// JobMetric represents a job execution metric
type JobMetric struct {
	Timestamp       time.Time
	JobID           uint64
	UserID          uint64
	PipelineID      *uint64
	FileID          uint64
	ContentType     string
	FileSize        int64
	Status          string
	ProcessingTime  time.Duration
	CreatedAt       time.Time
	FinishedAt      *time.Time
	ErrorMessage    *string
}

// JobStatusTransition represents a job status change
type JobStatusTransition struct {
	Timestamp   time.Time
	JobID       uint64
	UserID      uint64
	FromStatus  string
	ToStatus    string
	TriggeredBy string
	Message     string
}

// PipelineExecutionLog represents a pipeline operation execution
type PipelineExecutionLog struct {
	Timestamp      time.Time
	JobID          uint64
	UserID         uint64
	PipelineID     *uint64
	OperationName  string
	OperationType  string
	Duration       time.Duration
	InputSize      int64
	OutputSize     int64
	Success        bool
	ErrorMessage   *string
}
