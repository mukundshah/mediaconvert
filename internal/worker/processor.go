package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/mukund/mediaconvert/internal/analytics"
	"github.com/mukund/mediaconvert/internal/config"
	"github.com/mukund/mediaconvert/internal/models"
	"github.com/mukund/mediaconvert/internal/pipeline"
	"gorm.io/gorm"
)

// JobProcessor handles job processing
type JobProcessor struct {
	db          *gorm.DB
	minioClient *minio.Client
	config      *config.Config
	redis       *RedisClient
	analytics   *analytics.Client
}

// NewJobProcessor creates a new job processor
func NewJobProcessor(db *gorm.DB, minioClient *minio.Client, cfg *config.Config, redis *RedisClient, analyticsClient *analytics.Client) *JobProcessor {
	return &JobProcessor{
		db:          db,
		minioClient: minioClient,
		config:      cfg,
		redis:       redis,
		analytics:   analyticsClient,
	}
}

// ProcessJob processes a single job
func (p *JobProcessor) ProcessJob(jobID uint) error {
	startTime := time.Now()
	ctx := context.Background()

	// Load job with relationships
	var job models.Job
	if err := p.db.Preload("File").Preload("Pipeline").First(&job, jobID).Error; err != nil {
		return fmt.Errorf("failed to load job: %w", err)
	}

	// Update status to processing
	job.Status = models.JobStatusProcessing
	p.db.Save(&job)
	recordStatusChange(p.db, job.ID, models.JobStatusPending, models.JobStatusProcessing, "Worker started processing", "worker")

	// Record status transition in analytics
	if p.analytics != nil {
		p.analytics.RecordJobStatusTransition(ctx, analytics.JobStatusTransition{
			Timestamp:   time.Now(),
			JobID:       uint64(job.ID),
			UserID:      uint64(job.File.UserID),
			FromStatus:  string(models.JobStatusPending),
			ToStatus:    string(models.JobStatusProcessing),
			TriggeredBy: "worker",
			Message:     "Worker started processing",
		})
	}

	// Create work directory
	workDir := filepath.Join(os.TempDir(), fmt.Sprintf("job-%d", job.ID))
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return p.failJob(&job, fmt.Errorf("failed to create work directory: %w", err))
	}
	defer os.RemoveAll(workDir) // Cleanup

	// Download input file from S3
	inputFile := filepath.Join(workDir, "input"+filepath.Ext(job.File.OriginalName))
	if err := p.downloadFile(job.File.S3Key, inputFile); err != nil {
		return p.failJob(&job, fmt.Errorf("failed to download file: %w", err))
	}

	// Parse pipeline
	var pipelineObj *pipeline.Pipeline
	var err error
	if job.Pipeline != nil {
		if job.Pipeline.Format == models.PipelineFormatYAML {
			pipelineObj, err = pipeline.ParseYAML([]byte(job.Pipeline.Content))
		} else {
			pipelineObj, err = pipeline.ParseJSON([]byte(job.Pipeline.Content))
		}
		if err != nil {
			return p.failJob(&job, fmt.Errorf("failed to parse pipeline: %w", err))
		}
	} else {
		return p.failJob(&job, fmt.Errorf("no pipeline specified"))
	}

	// Execute pipeline
	outputFiles, err := ExecutePipeline(pipelineObj, inputFile, workDir)
	if err != nil {
		return p.failJob(&job, fmt.Errorf("pipeline execution failed: %w", err))
	}

	// Upload results to S3
	resultPaths, err := p.uploadResults(job.File.UserID, job.ID, outputFiles)
	if err != nil {
		return p.failJob(&job, fmt.Errorf("failed to upload results: %w", err))
	}

	// Update job as completed
	now := time.Now()
	job.Status = models.JobStatusCompleted
	job.FinishedAt = &now

	// Convert result info to JSON
	resultData := map[string]interface{}{
		"output_files": resultPaths,
		"processed_at": now,
	}
	resultJSON, _ := json.Marshal(resultData)
	job.ResultInfo = resultJSON

	p.db.Save(&job)
	recordStatusChange(p.db, job.ID, models.JobStatusProcessing, models.JobStatusCompleted, "Job completed successfully", "worker")

	// Record metrics and status transition in analytics
	if p.analytics != nil {
		processingTime := time.Since(startTime)
		var pipelineID *uint64
		if job.PipelineID != nil {
			id := uint64(*job.PipelineID)
			pipelineID = &id
		}

		// Record job metric
		p.analytics.RecordJobMetric(ctx, analytics.JobMetric{
			Timestamp:      time.Now(),
			JobID:          uint64(job.ID),
			UserID:         uint64(job.File.UserID),
			PipelineID:     pipelineID,
			FileID:         uint64(job.FileID),
			ContentType:    job.File.ContentType,
			FileSize:       job.File.Size,
			Status:         string(job.Status),
			ProcessingTime: processingTime,
			CreatedAt:      job.CreatedAt,
			FinishedAt:     job.FinishedAt,
			ErrorMessage:   nil,
		})

		// Record status transition
		p.analytics.RecordJobStatusTransition(ctx, analytics.JobStatusTransition{
			Timestamp:   time.Now(),
			JobID:       uint64(job.ID),
			UserID:      uint64(job.File.UserID),
			FromStatus:  string(models.JobStatusProcessing),
			ToStatus:    string(models.JobStatusCompleted),
			TriggeredBy: "worker",
			Message:     "Job completed successfully",
		})
	}

	fmt.Printf("Job %d completed successfully\n", job.ID)
	return nil
}

func (p *JobProcessor) downloadFile(s3Key, destPath string) error {
	return p.minioClient.FGetObject(context.Background(), p.config.S3Bucket, s3Key, destPath, minio.GetObjectOptions{})
}

func (p *JobProcessor) uploadResults(userID, jobID uint, files []string) ([]string, error) {
	var s3Keys []string

	for _, filePath := range files {
		// Build S3 key
		fileName := filepath.Base(filePath)
		s3Key := fmt.Sprintf("users/%d/results/job-%d/%s", userID, jobID, fileName)

		// Upload to S3
		_, err := p.minioClient.FPutObject(context.Background(), p.config.S3Bucket, s3Key, filePath, minio.PutObjectOptions{})
		if err != nil {
			return nil, err
		}

		s3Keys = append(s3Keys, s3Key)
	}

	return s3Keys, nil
}

func (p *JobProcessor) failJob(job *models.Job, err error) error {
	ctx := context.Background()
	now := time.Now()
	job.Status = models.JobStatusFailed
	job.Error = err.Error()
	job.FinishedAt = &now
	p.db.Save(job)
	recordStatusChange(p.db, job.ID, models.JobStatusProcessing, models.JobStatusFailed, err.Error(), "worker")

	// Record metrics and status transition in analytics
	if p.analytics != nil {
		errorMsg := err.Error()
		var pipelineID *uint64
		if job.PipelineID != nil {
			id := uint64(*job.PipelineID)
			pipelineID = &id
		}

		// Record job metric
		p.analytics.RecordJobMetric(ctx, analytics.JobMetric{
			Timestamp:      time.Now(),
			JobID:          uint64(job.ID),
			UserID:         uint64(job.File.UserID),
			PipelineID:     pipelineID,
			FileID:         uint64(job.FileID),
			ContentType:    job.File.ContentType,
			FileSize:       job.File.Size,
			Status:         string(job.Status),
			ProcessingTime: 0, // Failed jobs don't have meaningful processing time
			CreatedAt:      job.CreatedAt,
			FinishedAt:     job.FinishedAt,
			ErrorMessage:   &errorMsg,
		})

		// Record status transition
		p.analytics.RecordJobStatusTransition(ctx, analytics.JobStatusTransition{
			Timestamp:   time.Now(),
			JobID:       uint64(job.ID),
			UserID:      uint64(job.File.UserID),
			FromStatus:  string(models.JobStatusProcessing),
			ToStatus:    string(models.JobStatusFailed),
			TriggeredBy: "worker",
			Message:     err.Error(),
		})
	}

	return err
}

func recordStatusChange(db *gorm.DB, jobID uint, fromStatus, toStatus models.JobStatus, message, triggeredBy string) {
	history := models.JobStatusHistory{
		JobID:       jobID,
		FromStatus:  fromStatus,
		ToStatus:    toStatus,
		Message:     message,
		TriggeredBy: triggeredBy,
	}
	db.Create(&history)
}
