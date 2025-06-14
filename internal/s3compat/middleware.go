package s3compat

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mukund/mediaconvert/internal/models"
	"gorm.io/gorm"
)

// S3AuthMiddleware validates S3 API requests using AWS Signature V4
func S3AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var signatureInfo *SignatureInfo
		var err error

		// Check if it's a presigned URL request
		if c.Request.URL.Query().Get("X-Amz-Algorithm") != "" {
			signatureInfo, err = ParsePresignedURL(c.Request)
		} else {
			// Check for Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization"})
				c.Abort()
				return
			}
			signatureInfo, err = ParseAuthorizationHeader(authHeader)
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization format: " + err.Error()})
			c.Abort()
			return
		}

		// Find credential by access key
		var credential models.S3Credential
		if err := db.Where("access_key = ? AND is_active = ?", signatureInfo.AccessKey, true).
			Preload("User").
			First(&credential).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{"error": "Invalid access key"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed"})
			}
			c.Abort()
			return
		}

		// Validate signature
		if err := ValidateSignature(c.Request, credential.SecretKey, signatureInfo); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Signature validation failed: " + err.Error()})
			c.Abort()
			return
		}

		// Determine the requested bucket name
		requestedBucket := c.Param("bucket")
		if requestedBucket == "" {
			// Try to get from path if not in param (e.g., path-style access)
			path := c.Request.URL.Path
			parts := strings.Split(strings.Trim(path, "/"), "/")
			if len(parts) > 0 {
				requestedBucket = parts[0]
			}
		}

		// Verify bucket access
		if requestedBucket != "" && requestedBucket != credential.BucketName {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this bucket"})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", credential.UserID)
		c.Set("s3_credential_id", credential.ID)
		c.Set("bucket_name", credential.BucketName)

		c.Next()
	}
}

// GetBucketName retrieves the bucket name from context
func GetBucketName(c *gin.Context) (string, bool) {
	bucketName, exists := c.Get("bucket_name")
	if !exists {
		return "", false
	}
	return bucketName.(string), true
}

// GetUserIDFromS3Context retrieves user ID from S3 context
func GetUserIDFromS3Context(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// ValidateBucketAccess checks if the requested bucket matches the user's credential
func ValidateBucketAccess(c *gin.Context) bool {
	requestedBucket := c.Param("bucket")
	if requestedBucket == "" {
		// Try to get from path
		path := c.Request.URL.Path
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) > 0 {
			requestedBucket = parts[0]
		}
	}

	credentialBucket, exists := GetBucketName(c)
	if !exists {
		return false
	}

	return requestedBucket == credentialBucket
}
