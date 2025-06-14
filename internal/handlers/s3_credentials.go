package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mukund/mediaconvert/internal/auth"
	"github.com/mukund/mediaconvert/internal/models"
	"github.com/mukund/mediaconvert/internal/s3compat"
	"gorm.io/gorm"
)

type S3CredentialHandler struct {
	db *gorm.DB
}

func NewS3CredentialHandler(db *gorm.DB) *S3CredentialHandler {
	return &S3CredentialHandler{db: db}
}

type S3CredentialResponse struct {
	ID         uint   `json:"id"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key,omitempty"` // Only returned on creation
	BucketName string `json:"bucket_name"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  string `json:"created_at"`
}

type CreateCredentialRequest struct {
	BucketName string `json:"bucket_name,omitempty"` // Optional custom bucket name
}

// CreateCredentials generates new S3 credentials for the user
func (h *S3CredentialHandler) CreateCredentials(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse request body
	var req CreateCredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if user already has active credentials
	var existingCount int64
	h.db.Model(&models.S3Credential{}).Where("user_id = ? AND is_active = ?", userID, true).Count(&existingCount)
	if existingCount >= 5 { // Limit to 5 active credentials per user
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum number of active credentials reached (5)"})
		return
	}

	// Generate credentials
	accessKey, err := s3compat.GenerateAccessKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access key"})
		return
	}

	secretKey, err := s3compat.GenerateSecretKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate secret key"})
		return
	}

	// Handle bucket name
	var bucketName string
	if req.BucketName != "" {
		// Validate custom bucket name
		if !isValidBucketName(req.BucketName) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bucket name. Must be 3-63 characters, lowercase letters, numbers, and hyphens only"})
			return
		}

		// Check if bucket name is already taken (globally unique)
		var existingBucket models.S3Credential
		if err := h.db.Where("bucket_name = ?", req.BucketName).First(&existingBucket).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Bucket name already taken"})
			return
		}

		bucketName = req.BucketName
	} else {
		// Generate bucket name
		bucketName, err = s3compat.GenerateBucketName(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate bucket name"})
			return
		}
	}

	// Hash secret key
	// hashedSecretKey, err := s3compat.HashSecretKey(secretKey)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash secret key"})
	// 	return
	// }

	// Create credential record
	credential := models.S3Credential{
		UserID:     userID,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		BucketName: bucketName,
		IsActive:   true,
	}

	if err := h.db.Create(&credential).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create credentials"})
		return
	}

	c.JSON(http.StatusCreated, S3CredentialResponse{
		ID:         credential.ID,
		AccessKey:  credential.AccessKey,
		SecretKey:  secretKey, // Return plain secret key only on creation
		BucketName: credential.BucketName,
		IsActive:   credential.IsActive,
		CreatedAt:  credential.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// ListCredentials returns all credentials for the user
func (h *S3CredentialHandler) ListCredentials(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var credentials []models.S3Credential
	if err := h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&credentials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch credentials"})
		return
	}

	response := make([]S3CredentialResponse, len(credentials))
	for i, cred := range credentials {
		response[i] = S3CredentialResponse{
			ID:         cred.ID,
			AccessKey:  cred.AccessKey,
			BucketName: cred.BucketName,
			IsActive:   cred.IsActive,
			CreatedAt:  cred.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"credentials": response})
}

// RevokeCredentials deactivates credentials
func (h *S3CredentialHandler) RevokeCredentials(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	credentialID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credential ID"})
		return
	}

	// Find and update credential
	var credential models.S3Credential
	if err := h.db.Where("id = ? AND user_id = ?", credentialID, userID).First(&credential).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Credential not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch credential"})
		}
		return
	}

	credential.IsActive = false
	if err := h.db.Save(&credential).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credentials revoked successfully"})
}

// CheckBucketAvailability checks if a bucket name is available
func (h *S3CredentialHandler) CheckBucketAvailability(c *gin.Context) {
	bucketName := c.Query("name")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bucket name is required"})
		return
	}

	// Validate bucket name format
	if !isValidBucketName(bucketName) {
		c.JSON(http.StatusOK, gin.H{
			"available": false,
			"reason":    "Invalid bucket name format. Must be 3-63 characters, lowercase letters, numbers, and hyphens only",
		})
		return
	}

	// Check if bucket name is taken
	var existingBucket models.S3Credential
	if err := h.db.Where("bucket_name = ?", bucketName).First(&existingBucket).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{
			"available": false,
			"reason":    "Bucket name already taken",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"available": true,
		"name":      bucketName,
	})
}

// isValidBucketName validates bucket name according to S3 naming rules
func isValidBucketName(name string) bool {
	// Length check
	if len(name) < 3 || len(name) > 63 {
		return false
	}

	// Must start and end with lowercase letter or number
	if !isLowerAlphaNum(rune(name[0])) || !isLowerAlphaNum(rune(name[len(name)-1])) {
		return false
	}

	// Check each character
	for i, c := range name {
		if !isLowerAlphaNum(c) && c != '-' {
			return false
		}
		// No consecutive hyphens
		if c == '-' && i > 0 && name[i-1] == '-' {
			return false
		}
	}

	return true
}

func isLowerAlphaNum(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}
