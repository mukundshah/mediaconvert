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
)

type Job struct {
	gorm.Model
	FileID     uint
	File       File
	Status     JobStatus `gorm:"default:'pending'"`
	Pipeline   datatypes.JSON // JSON definition of the processing pipeline
	ResultInfo datatypes.JSON // JSON storing result details (e.g., output paths)
	Error      string
	FinishedAt *time.Time
}
