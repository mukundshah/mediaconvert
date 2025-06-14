package s3compat

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// GenerateAccessKey generates a random access key (20 characters, AWS-style)
func GenerateAccessKey() (string, error) {
	// Generate 15 random bytes (will encode to 24 base32 chars, we'll take 20)
	bytes := make([]byte, 15)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Encode to base32 and take first 20 characters
	encoded := base32.StdEncoding.EncodeToString(bytes)
	accessKey := "AKIA" + encoded[:16] // AWS-style prefix + 16 chars = 20 total
	return accessKey, nil
}

// GenerateSecretKey generates a random secret key (40 characters)
func GenerateSecretKey() (string, error) {
	// Generate 30 random bytes (will encode to 48 base32 chars, we'll take 40)
	bytes := make([]byte, 30)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Encode to base32 and take first 40 characters
	encoded := base32.StdEncoding.EncodeToString(bytes)
	return encoded[:40], nil
}

// GenerateBucketName generates a unique bucket name for a user
func GenerateBucketName(userID uint) (string, error) {
	// Generate random suffix
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	suffix := base32.StdEncoding.EncodeToString(bytes)
	suffix = strings.ToLower(suffix[:8])

	return fmt.Sprintf("user-%d-%s", userID, suffix), nil
}

// HashSecretKey hashes a secret key using bcrypt
func HashSecretKey(secretKey string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(secretKey), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckSecretKey verifies a secret key against a hash
func CheckSecretKey(hashedSecretKey, secretKey string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedSecretKey), []byte(secretKey))
}
