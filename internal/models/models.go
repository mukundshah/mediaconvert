package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/datatypes"
)

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"` // Hashed
}

type File struct {
	gorm.Model
	UserID      uint
	User        User
	OriginalName string
	S3Key       string `gorm:"uniqueIndex;not null"`
	Size        int64
	ContentType string
}

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
	JobStatusCanceled   JobStatus = "canceled"
)

type PipelineFormat string

const (
	PipelineFormatYAML PipelineFormat = "yaml"
	PipelineFormatJSON PipelineFormat = "json"
)

type Pipeline struct {
	gorm.Model
	UserID  uint
	User    User
	Name    string         `gorm:"not null"`
	Format  PipelineFormat `gorm:"type:varchar(10);not null"`
	Content string         `gorm:"type:text;not null"`
}

type Job struct {
	gorm.Model
	FileID       uint
	File         File
	PipelineID   *uint          // Optional reference to a saved pipeline
	Pipeline     *Pipeline      // Relationship to saved pipeline
	PipelineData datatypes.JSON // Inline pipeline definition (for ad-hoc jobs or snapshot)
	Status       JobStatus      `gorm:"default:'pending'"`
	ResultInfo   datatypes.JSON // JSON storing result details (e.g., output paths)
	Error        string
	FinishedAt   *time.Time
}

type JobStatusHistory struct {
	gorm.Model
	JobID       uint
	Job         Job
	FromStatus  JobStatus `gorm:"type:varchar(20)"`
	ToStatus    JobStatus `gorm:"type:varchar(20);not null"`
	Message     string
	TriggeredBy string `gorm:"type:varchar(20);not null"` // "user", "system", "worker"
}

type S3Credential struct {
	gorm.Model
	UserID     uint
	User       User
	AccessKey  string `gorm:"uniqueIndex;not null;size:20"`
	SecretKey  string `gorm:"not null"` // Hashed
	BucketName string `gorm:"uniqueIndex;not null"`
	IsActive   bool   `gorm:"default:true"`
}
