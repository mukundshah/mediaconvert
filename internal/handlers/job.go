package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mukund/mediaconvert/internal/auth"
	"github.com/mukund/mediaconvert/internal/models"
	"gorm.io/gorm"
)

type JobHandler struct {
	db *gorm.DB
}

func NewJobHandler(db *gorm.DB) *JobHandler {
	return &JobHandler{db: db}
}

type JobListResponse struct {
	Jobs       []JobDetail        `json:"jobs"`
	Pagination PaginationResponse `json:"pagination"`
}

type PaginationResponse struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type JobDetail struct {
	ID           uint                   `json:"id"`
	FileID       uint                   `json:"file_id"`
	File         *FileInfo              `json:"file,omitempty"`
	PipelineID   *uint                  `json:"pipeline_id,omitempty"`
	Pipeline     *PipelineInfo          `json:"pipeline,omitempty"`
	PipelineData map[string]interface{} `json:"pipeline_data,omitempty"`
	Status       string                 `json:"status"`
	ResultInfo   map[string]interface{} `json:"result_info,omitempty"`
	Error        string                 `json:"error,omitempty"`
	CreatedAt    string                 `json:"created_at"`
	FinishedAt   *string                `json:"finished_at,omitempty"`
}

type FileInfo struct {
	ID           uint   `json:"id"`
	OriginalName string `json:"original_name"`
	S3Key        string `json:"s3_key"`
	Size         int64  `json:"size"`
	ContentType  string `json:"content_type"`
}

type PipelineInfo struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Format  string `json:"format"`
	Content string `json:"content,omitempty"`
}

// ListJobs returns a paginated list of jobs for the authenticated user
func (h *JobHandler) ListJobs(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse query parameters
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Build query
	query := h.db.Model(&models.Job{}).
		Joins("JOIN files ON files.id = jobs.file_id").
		Where("files.user_id = ?", userID)

	if status != "" {
		query = query.Where("jobs.status = ?", status)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get jobs with preloaded relationships
	var jobs []models.Job
	if err := query.
		Preload("File").
		Preload("Pipeline").
		Order("jobs.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	// Convert to response format
	jobDetails := make([]JobDetail, len(jobs))
	for i, job := range jobs {
		jobDetails[i] = convertToJobDetail(job, false)
	}

	c.JSON(http.StatusOK, JobListResponse{
		Jobs: jobDetails,
		Pagination: PaginationResponse{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	})
}

// GetJob returns a single job by ID
func (h *JobHandler) GetJob(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	jobID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	var job models.Job
	if err := h.db.
		Preload("File").
		Preload("Pipeline").
		First(&job, jobID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job"})
		}
		return
	}

	// Verify ownership
	if job.File.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, convertToJobDetail(job, true))
}

// CancelJob cancels a pending or processing job
func (h *JobHandler) CancelJob(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	jobID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	var job models.Job
	if err := h.db.
		Preload("File").
		First(&job, jobID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job"})
		}
		return
	}

	// Verify ownership
	if job.File.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Check if job can be cancelled
	if job.Status != models.JobStatusPending && job.Status != models.JobStatusProcessing {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job cannot be cancelled (already completed or failed)"})
		return
	}

	// Record old status for history
	oldStatus := job.Status

	// Update job status to canceled
	job.Status = models.JobStatusCanceled
	job.Error = "Job cancelled by user"
	now := time.Now()
	job.FinishedAt = &now

	if err := h.db.Save(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel job"})
		return
	}

	// Record status change
	if err := recordStatusChange(h.db, job.ID, oldStatus, models.JobStatusCanceled, "Job cancelled by user", "user"); err != nil {
		log.Printf("Failed to record status change: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Job cancelled successfully",
		"job_id":  job.ID,
	})
}

// RerunJob creates a new job with the same configuration as an existing job
func (h *JobHandler) RerunJob(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	jobID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	var originalJob models.Job
	if err := h.db.
		Preload("File").
		Preload("Pipeline").
		First(&originalJob, jobID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job"})
		}
		return
	}

	// Verify ownership
	if originalJob.File.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Create new job with same configuration
	newJob := models.Job{
		FileID:       originalJob.FileID,
		PipelineID:   originalJob.PipelineID,
		PipelineData: originalJob.PipelineData,
		Status:       models.JobStatusPending,
	}

	if err := h.db.Create(&newJob).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new job"})
		return
	}

	// Load the new job with relationships
	if err := h.db.
		Preload("File").
		Preload("Pipeline").
		First(&newJob, newJob.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch new job"})
		return
	}

	// Record initial status for new job
	if err := recordStatusChange(h.db, newJob.ID, "", models.JobStatusPending, "Job created via rerun", "user"); err != nil {
		log.Printf("Failed to record status change: %v", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Job rerun successfully",
		"original_job": originalJob.ID,
		"new_job":      convertToJobDetail(newJob, false),
	})
}

func convertToJobDetail(job models.Job, includeContent bool) JobDetail {
	detail := JobDetail{
		ID:        job.ID,
		FileID:    job.FileID,
		Status:    string(job.Status),
		Error:     job.Error,
		CreatedAt: job.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if job.File.ID > 0 {
		detail.File = &FileInfo{
			ID:           job.File.ID,
			OriginalName: job.File.OriginalName,
			S3Key:        job.File.S3Key,
			Size:         job.File.Size,
			ContentType:  job.File.ContentType,
		}
	}

	if job.PipelineID != nil && job.Pipeline != nil {
		detail.PipelineID = job.PipelineID
		detail.Pipeline = &PipelineInfo{
			ID:     job.Pipeline.ID,
			Name:   job.Pipeline.Name,
			Format: string(job.Pipeline.Format),
		}
		if includeContent {
			detail.Pipeline.Content = job.Pipeline.Content
		}
	}

	if job.FinishedAt != nil {
		finishedStr := job.FinishedAt.Format("2006-01-02T15:04:05Z07:00")
		detail.FinishedAt = &finishedStr
	}

	return detail
}

// recordStatusChange creates a JobStatusHistory record
func recordStatusChange(db *gorm.DB, jobID uint, fromStatus, toStatus models.JobStatus, message, triggeredBy string) error {
	history := models.JobStatusHistory{
		JobID:       jobID,
		FromStatus:  fromStatus,
		ToStatus:    toStatus,
		Message:     message,
		TriggeredBy: triggeredBy,
	}
	return db.Create(&history).Error
}
