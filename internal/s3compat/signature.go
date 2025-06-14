package s3compat

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// SignatureInfo contains parsed signature information from a request
type SignatureInfo struct {
	AccessKey        string
	SignedHeaders    string
	Signature        string
	Date             string
	Region           string
	Service          string
	Algorithm        string
	CredentialScope  string
	IsPresigned      bool
	ExpiresIn        int64
}

// ParseAuthorizationHeader parses the AWS Signature V4 Authorization header
func ParseAuthorizationHeader(authHeader string) (*SignatureInfo, error) {
	// Format: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request, SignedHeaders=host;range;x-amz-date, Signature=fe5f80f77d5fa3beca038a248ff027d0445342fe2855ddc963176630326f1024

	if !strings.HasPrefix(authHeader, "AWS4-HMAC-SHA256 ") {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	info := &SignatureInfo{
		Algorithm: "AWS4-HMAC-SHA256",
	}

	parts := strings.Split(strings.TrimPrefix(authHeader, "AWS4-HMAC-SHA256 "), ", ")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "Credential":
			credParts := strings.Split(value, "/")
			if len(credParts) >= 5 {
				info.AccessKey = credParts[0]
				info.Date = credParts[1]
				info.Region = credParts[2]
				info.Service = credParts[3]
				info.CredentialScope = strings.Join(credParts[1:], "/")
			}
		case "SignedHeaders":
			info.SignedHeaders = value
		case "Signature":
			info.Signature = value
		}
	}

	return info, nil
}

// ParsePresignedURL parses signature information from a presigned URL
func ParsePresignedURL(r *http.Request) (*SignatureInfo, error) {
	query := r.URL.Query()

	algorithm := query.Get("X-Amz-Algorithm")
	if algorithm != "AWS4-HMAC-SHA256" {
		return nil, fmt.Errorf("invalid or missing algorithm")
	}

	credential := query.Get("X-Amz-Credential")
	if credential == "" {
		return nil, fmt.Errorf("missing credential")
	}

	credParts := strings.Split(credential, "/")
	if len(credParts) < 5 {
		return nil, fmt.Errorf("invalid credential format")
	}

	expiresStr := query.Get("X-Amz-Expires")
	if expiresStr == "" {
		return nil, fmt.Errorf("missing expires")
	}

	var expiresIn int64
	fmt.Sscanf(expiresStr, "%d", &expiresIn)

	info := &SignatureInfo{
		Algorithm:       algorithm,
		AccessKey:       credParts[0],
		Date:            credParts[1],
		Region:          credParts[2],
		Service:         credParts[3],
		CredentialScope: strings.Join(credParts[1:], "/"),
		SignedHeaders:   query.Get("X-Amz-SignedHeaders"),
		Signature:       query.Get("X-Amz-Signature"),
		IsPresigned:     true,
		ExpiresIn:       expiresIn,
	}

	return info, nil
}

// ValidateSignature validates an AWS Signature V4 signature
func ValidateSignature(r *http.Request, secretKey string, signatureInfo *SignatureInfo) error {
	// Calculate the expected signature
	expectedSignature, err := CalculateSignature(r, secretKey, signatureInfo)
	if err != nil {
		return fmt.Errorf("failed to calculate signature: %w", err)
	}

	// Compare signatures
	if !hmac.Equal([]byte(expectedSignature), []byte(signatureInfo.Signature)) {
		return fmt.Errorf("signature mismatch")
	}

	// For presigned URLs, check expiration
	if signatureInfo.IsPresigned {
		dateHeader := r.URL.Query().Get("X-Amz-Date")
		if dateHeader == "" {
			return fmt.Errorf("missing X-Amz-Date in presigned URL")
		}

		requestTime, err := time.Parse("20060102T150405Z", dateHeader)
		if err != nil {
			return fmt.Errorf("invalid date format: %w", err)
		}

		expirationTime := requestTime.Add(time.Duration(signatureInfo.ExpiresIn) * time.Second)
		if time.Now().After(expirationTime) {
			return fmt.Errorf("presigned URL has expired")
		}
	}

	return nil
}

// CalculateSignature calculates the AWS Signature V4 signature
func CalculateSignature(r *http.Request, secretKey string, signatureInfo *SignatureInfo) (string, error) {
	// Step 1: Create canonical request
	canonicalRequest := createCanonicalRequest(r, signatureInfo)

	// Step 2: Create string to sign
	stringToSign := createStringToSign(r, canonicalRequest, signatureInfo)

	// Step 3: Calculate signing key
	signingKey := getSigningKey(secretKey, signatureInfo.Date, signatureInfo.Region, signatureInfo.Service)

	// Step 4: Calculate signature
	signature := hmacSHA256(signingKey, stringToSign)

	return hex.EncodeToString(signature), nil
}

func createCanonicalRequest(r *http.Request, signatureInfo *SignatureInfo) string {
	// Canonical URI
	canonicalURI := r.URL.Path
	if canonicalURI == "" {
		canonicalURI = "/"
	}

	// Canonical query string
	canonicalQueryString := getCanonicalQueryString(r, signatureInfo.IsPresigned)

	// Canonical headers
	canonicalHeaders, signedHeadersList := getCanonicalHeaders(r, signatureInfo.SignedHeaders)

	// Payload hash
	payloadHash := getPayloadHash(r)

	// Combine into canonical request
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		r.Method,
		canonicalURI,
		canonicalQueryString,
		canonicalHeaders,
		signedHeadersList,
		payloadHash,
	)
}

func getCanonicalQueryString(r *http.Request, isPresigned bool) string {
	query := r.URL.Query()

	// For presigned URLs, exclude the signature parameter
	if isPresigned {
		query.Del("X-Amz-Signature")
	}

	var keys []string
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		for _, v := range query[k] {
			parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(v))
		}
	}

	return strings.Join(parts, "&")
}

func getCanonicalHeaders(r *http.Request, signedHeaders string) (string, string) {
	headers := make(map[string][]string)
	headerNames := strings.Split(signedHeaders, ";")

	for _, name := range headerNames {
		name = strings.ToLower(strings.TrimSpace(name))

		// Special case for 'host' header
		if name == "host" {
			headers[name] = []string{r.Host}
		} else if values := r.Header.Values(name); len(values) > 0 {
			headers[name] = values
		}
	}

	var keys []string
	for k := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var canonicalHeaders []string
	for _, k := range keys {
		values := headers[k]
		canonicalHeaders = append(canonicalHeaders, k+":"+strings.Join(values, ",")+"\n")
	}

	return strings.Join(canonicalHeaders, ""), strings.Join(keys, ";")
}

func getPayloadHash(r *http.Request) string {
	// For presigned URLs or unsigned payload
	if payloadHash := r.Header.Get("X-Amz-Content-Sha256"); payloadHash != "" {
		return payloadHash
	}
	return "UNSIGNED-PAYLOAD"
}

func createStringToSign(r *http.Request, canonicalRequest string, signatureInfo *SignatureInfo) string {
	hashedCanonicalRequest := sha256Hash(canonicalRequest)

	// Get timestamp from X-Amz-Date header or query parameter
	timestamp := r.Header.Get("X-Amz-Date")
	if timestamp == "" && signatureInfo.IsPresigned {
		timestamp = r.URL.Query().Get("X-Amz-Date")
	}
	if timestamp == "" {
		// Fallback to date with midnight time
		timestamp = signatureInfo.Date + "T000000Z"
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s",
		signatureInfo.Algorithm,
		timestamp,
		signatureInfo.CredentialScope,
		hashedCanonicalRequest,
	)
}

func getSigningKey(secretKey, date, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secretKey), date)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	kSigning := hmacSHA256(kService, "aws4_request")
	return kSigning
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func sha256Hash(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
