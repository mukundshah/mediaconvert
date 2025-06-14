package s3compat

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/mukund/mediaconvert/internal/config"
	"github.com/mukund/mediaconvert/internal/models"
	"github.com/mukund/mediaconvert/internal/worker"
	"gorm.io/gorm"
)

type S3Handler struct {
	db       *gorm.DB
	s3Client *s3.Client
	config   *config.Config
	redis    *worker.RedisClient
}

func NewS3Handler(db *gorm.DB, s3Client *s3.Client, cfg *config.Config, redis *worker.RedisClient) *S3Handler {
	return &S3Handler{
		db:       db,
		s3Client: s3Client,
		config:   cfg,
		redis:    redis,
	}
}

// PutObject handles S3 PUT object requests (upload)
func (h *S3Handler) PutObject(c *gin.Context) {
	userID, _ := GetUserIDFromS3Context(c)
	bucket := c.Param("bucket")
	key := c.Param("key")

	if key == "" {
		key = strings.TrimPrefix(c.Request.URL.Path, "/"+bucket+"/")
	}

	// Validate bucket access
	if !ValidateBucketAccess(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to bucket"})
		return
	}

	// Build S3 key with user prefix
	s3Key := fmt.Sprintf("users/%d/%s", userID, key)

	// Upload to S3
	_, err := h.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(h.config.S3Bucket),
		Key:         aws.String(s3Key),
		Body:        c.Request.Body,
		ContentType: aws.String(c.GetHeader("Content-Type")),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Create File record
	fileRecord := models.File{
		UserID:       userID,
		OriginalName: path.Base(key),
		S3Key:        s3Key,
		Size:         c.Request.ContentLength,
		ContentType:  c.GetHeader("Content-Type"),
	}

	if err := h.db.Create(&fileRecord).Error; err != nil {
		// File uploaded but DB record failed - log but continue
		fmt.Printf("Warning: Failed to create file record: %v\n", err)
		c.Header("ETag", fmt.Sprintf("\"%s\"", s3Key))
		c.Status(http.StatusOK)
		return
	}

	// Check for pipeline name in metadata (X-Amz-Meta-Pipeline header)
	pipelineName := c.GetHeader("X-Amz-Meta-Pipeline")
	if pipelineName != "" {
		// Look up pipeline by name
		var pipelineRecord models.Pipeline
		if err := h.db.Where("user_id = ? AND name = ?", userID, pipelineName).First(&pipelineRecord).Error; err == nil {
			// Create job with pipeline
			job := models.Job{
				FileID:     fileRecord.ID,
				PipelineID: &pipelineRecord.ID,
				Status:     models.JobStatusPending,
			}

			if err := h.db.Create(&job).Error; err != nil {
				fmt.Printf("Warning: Failed to create job: %v\n", err)
			} else {
				// Publish job notification to Redis
				if h.redis != nil {
					if err := h.redis.PublishJobNotification(job.ID); err != nil {
						fmt.Printf("Warning: Failed to publish job notification: %v\n", err)
					}
				}
			}
		} else {
			fmt.Printf("Warning: Pipeline '%s' not found for user %d\n", pipelineName, userID)
		}
	}

	// Return S3-compatible response
	c.Header("ETag", fmt.Sprintf("\"%s\"", s3Key))
	c.Status(http.StatusOK)
}

// GetObject handles S3 GET object requests (download)
func (h *S3Handler) GetObject(c *gin.Context) {
	userID, _ := GetUserIDFromS3Context(c)
	bucket := c.Param("bucket")
	key := c.Param("key")

	if key == "" {
		key = strings.TrimPrefix(c.Request.URL.Path, "/"+bucket+"/")
	}

	if !ValidateBucketAccess(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to bucket"})
		return
	}

	// Build S3 key
	s3Key := fmt.Sprintf("users/%d/%s", userID, key)

	// Get object from S3
	result, err := h.s3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(h.config.S3Bucket),
		Key:    aws.String(s3Key),
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}
	defer result.Body.Close()

	// Set headers
	if result.ContentType != nil {
		c.Header("Content-Type", *result.ContentType)
	}
	if result.ContentLength != nil {
		c.Header("Content-Length", fmt.Sprintf("%d", *result.ContentLength))
	}
	if result.ETag != nil {
		c.Header("ETag", *result.ETag)
	}
	if result.LastModified != nil {
		c.Header("Last-Modified", result.LastModified.Format(http.TimeFormat))
	}

	// Stream response
	c.Status(http.StatusOK)
	io.Copy(c.Writer, result.Body)
}

// HeadObject handles S3 HEAD object requests (metadata)
func (h *S3Handler) HeadObject(c *gin.Context) {
	userID, _ := GetUserIDFromS3Context(c)
	bucket := c.Param("bucket")
	key := c.Param("key")

	if key == "" {
		key = strings.TrimPrefix(c.Request.URL.Path, "/"+bucket+"/")
	}

	if !ValidateBucketAccess(c) {
		c.Status(http.StatusForbidden)
		return
	}

	s3Key := fmt.Sprintf("users/%d/%s", userID, key)

	// Head object from S3
	result, err := h.s3Client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(h.config.S3Bucket),
		Key:    aws.String(s3Key),
	})

	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Set headers
	if result.ContentType != nil {
		c.Header("Content-Type", *result.ContentType)
	}
	if result.ContentLength != nil {
		c.Header("Content-Length", fmt.Sprintf("%d", *result.ContentLength))
	}
	if result.ETag != nil {
		c.Header("ETag", *result.ETag)
	}
	if result.LastModified != nil {
		c.Header("Last-Modified", result.LastModified.Format(http.TimeFormat))
	}

	c.Status(http.StatusOK)
}

// DeleteObject handles S3 DELETE object requests
func (h *S3Handler) DeleteObject(c *gin.Context) {
	userID, _ := GetUserIDFromS3Context(c)
	bucket := c.Param("bucket")
	key := c.Param("key")

	if key == "" {
		key = strings.TrimPrefix(c.Request.URL.Path, "/"+bucket+"/")
	}

	if !ValidateBucketAccess(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to bucket"})
		return
	}

	s3Key := fmt.Sprintf("users/%d/%s", userID, key)

	// Delete from S3
	_, err := h.s3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(h.config.S3Bucket),
		Key:    aws.String(s3Key),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete object"})
		return
	}

	// Delete file record
	h.db.Where("user_id = ? AND s3_key = ?", userID, s3Key).Delete(&models.File{})

	c.Status(http.StatusNoContent)
}

// ListObjects handles S3 LIST objects requests
func (h *S3Handler) ListObjects(c *gin.Context) {
	userID, _ := GetUserIDFromS3Context(c)
	bucket := c.Param("bucket")

	if !ValidateBucketAccess(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to bucket"})
		return
	}

	prefix := fmt.Sprintf("users/%d/", userID)
	if reqPrefix := c.Query("prefix"); reqPrefix != "" {
		prefix = fmt.Sprintf("users/%d/%s", userID, reqPrefix)
	}

	// List objects from S3
	result, err := h.s3Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(h.config.S3Bucket),
		Prefix: aws.String(prefix),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list objects"})
		return
	}

	// Build S3-compatible XML response
	type S3Object struct {
		Key          string    `xml:"Key"`
		LastModified time.Time `xml:"LastModified"`
		ETag         string    `xml:"ETag"`
		Size         int64     `xml:"Size"`
		StorageClass string    `xml:"StorageClass"`
	}

	type ListBucketResult struct {
		Name     string     `xml:"Name"`
		Prefix   string     `xml:"Prefix"`
		Contents []S3Object `xml:"Contents"`
	}

	var contents []S3Object
	for _, obj := range result.Contents {
		// Strip user prefix from key
		displayKey := strings.TrimPrefix(*obj.Key, prefix)
		contents = append(contents, S3Object{
			Key:          displayKey,
			LastModified: *obj.LastModified,
			ETag:         *obj.ETag,
			Size:         *obj.Size,
			StorageClass: "STANDARD",
		})
	}

	response := ListBucketResult{
		Name:     bucket,
		Prefix:   c.Query("prefix"),
		Contents: contents,
	}

	c.XML(http.StatusOK, response)
}
