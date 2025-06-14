package s3compat

import (
	"fmt"
	"net/url"
	"time"
)

// PresignedURLParams contains parameters for generating a presigned URL
type PresignedURLParams struct {
	Method      string
	Bucket      string
	Key         string
	AccessKey   string
	SecretKey   string
	Region      string
	Expires     time.Duration
	Endpoint    string
}

// GeneratePresignedURL generates an AWS Signature V4 presigned URL
func GeneratePresignedURL(params PresignedURLParams) (string, error) {
	// Default region if not specified
	if params.Region == "" {
		params.Region = "us-east-1"
	}

	// Build base URL
	baseURL := fmt.Sprintf("%s/%s/%s", params.Endpoint, params.Bucket, params.Key)
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// Get current time
	now := time.Now().UTC()
	dateStamp := now.Format("20060102")
	amzDate := now.Format("20060102T150405Z")

	// Build credential scope
	credentialScope := fmt.Sprintf("%s/%s/s3/aws4_request", dateStamp, params.Region)
	credential := fmt.Sprintf("%s/%s", params.AccessKey, credentialScope)

	// Build query parameters
	query := parsedURL.Query()
	query.Set("X-Amz-Algorithm", "AWS4-HMAC-SHA256")
	query.Set("X-Amz-Credential", credential)
	query.Set("X-Amz-Date", amzDate)
	query.Set("X-Amz-Expires", fmt.Sprintf("%d", int(params.Expires.Seconds())))
	query.Set("X-Amz-SignedHeaders", "host")

	parsedURL.RawQuery = query.Encode()

	// Create canonical request for signing
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\nhost:%s\n\nhost\nUNSIGNED-PAYLOAD",
		params.Method,
		parsedURL.Path,
		parsedURL.RawQuery,
		parsedURL.Host,
	)

	// Create string to sign
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		amzDate,
		credentialScope,
		sha256Hash(canonicalRequest),
	)

	// Calculate signature
	signingKey := getSigningKey(params.SecretKey, dateStamp, params.Region, "s3")
	signature := hmacSHA256(signingKey, stringToSign)
	signatureHex := fmt.Sprintf("%x", signature)

	// Add signature to query
	query.Set("X-Amz-Signature", signatureHex)
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

// GeneratePresignedPutURL generates a presigned URL for uploading
func GeneratePresignedPutURL(bucket, key, accessKey, secretKey, endpoint string, expires time.Duration) (string, error) {
	return GeneratePresignedURL(PresignedURLParams{
		Method:    "PUT",
		Bucket:    bucket,
		Key:       key,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Endpoint:  endpoint,
		Expires:   expires,
	})
}

// GeneratePresignedGetURL generates a presigned URL for downloading
func GeneratePresignedGetURL(bucket, key, accessKey, secretKey, endpoint string, expires time.Duration) (string, error) {
	return GeneratePresignedURL(PresignedURLParams{
		Method:    "GET",
		Bucket:    bucket,
		Key:       key,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Endpoint:  endpoint,
		Expires:   expires,
	})
}
