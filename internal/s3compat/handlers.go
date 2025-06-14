package s3compat

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/mukund/mediaconvert/internal/config"
	"github.com/mukund/mediaconvert/internal/models"
	"github.com/mukund/mediaconvert/internal/worker"
	"gorm.io/gorm"
)

type S3Handler struct {
	db           *gorm.DB
	minioClient  *minio.Client
	config       *config.Config
	redis        *worker.RedisClient
}

func NewS3Handler(db *gorm.DB, minioClient *minio.Client, cfg *config.Config, redis *worker.RedisClient) *S3Handler {
	return &S3Handler{
		db:          db,
		minioClient: minioClient,
		config:      cfg,
		redis:       redis,
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

	// Remove leading slash from key if present (Gin's *key param includes it)
	key = strings.TrimPrefix(key, "/")

	// Validate bucket access
	if !ValidateBucketAccess(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to bucket"})
		return
	}

	// Build S3 key with user prefix
	s3Key := fmt.Sprintf("users/%d/%s", userID, key)

	// Read request body into memory
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("Failed to read request body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file"})
		return
	}

	// Upload to MinIO
	contentType := c.GetHeader("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = h.minioClient.PutObject(
		context.Background(),
		h.config.S3Bucket,
		s3Key,
		bytes.NewReader(bodyBytes),
		int64(len(bodyBytes)),
		minio.PutObjectOptions{ContentType: contentType},
	)

	if err != nil {
		fmt.Printf("MinIO Upload Error: %v\n", err)
		fmt.Printf("  Bucket: %s\n", h.config.S3Bucket)
		fmt.Printf("  Key: %s\n", s3Key)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
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
	key := strings.TrimPrefix(c.Param("key"), "/")

	// Validate bucket access
	if !ValidateBucketAccess(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to bucket"})
		return
	}

	// Build S3 key
	s3Key := fmt.Sprintf("users/%d/%s", userID, key)

	// Get object from MinIO
	object, err := h.minioClient.GetObject(context.Background(), h.config.S3Bucket, s3Key, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file"})
		return
	}
	defer object.Close()

	// Get object info for content type and size
	stat, err := object.Stat()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}

	// Set headers
	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size))
	c.Header("ETag", fmt.Sprintf("\"%s\"", stat.ETag))
	c.Header("Last-Modified", stat.LastModified.Format(http.TimeFormat))

	// Stream the object
	c.DataFromReader(http.StatusOK, stat.Size, stat.ContentType, object, nil)
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

	// Build S3 key
	s3Key := fmt.Sprintf("users/%d/%s", userID, key)

	// Get object metadata from MinIO
	stat, err := h.minioClient.StatObject(context.Background(), h.config.S3Bucket, s3Key, minio.StatObjectOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}

	// Set headers
	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size))
	c.Header("ETag", fmt.Sprintf("\"%s\"", stat.ETag))
	c.Header("Last-Modified", stat.LastModified.Format(http.TimeFormat))

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

	// Build S3 key
	s3Key := fmt.Sprintf("users/%d/%s", userID, key)

	// Delete from MinIO
	err := h.minioClient.RemoveObject(context.Background(), h.config.S3Bucket, s3Key, minio.RemoveObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
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

	// List objects from MinIO
	objectCh := h.minioClient.ListObjects(context.Background(), h.config.S3Bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

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
	for object := range objectCh {
		if object.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list objects"})
			return
		}

		// Strip user prefix from key
		displayKey := strings.TrimPrefix(object.Key, prefix)
		contents = append(contents, S3Object{
			Key:          displayKey,
			LastModified: object.LastModified,
			ETag:         object.ETag,
			Size:         object.Size,
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
