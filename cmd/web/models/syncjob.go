package models

import (
	"database/sql"
	"errors"
	"time"
)

// SyncJobStatus represents the status of a sync job
type SyncJobStatus string

// Supported sync job statuses
const (
	SyncJobStatusPending   SyncJobStatus = "pending"
	SyncJobStatusRunning   SyncJobStatus = "running"
	SyncJobStatusCompleted SyncJobStatus = "completed"
	SyncJobStatusFailed    SyncJobStatus = "failed"
	SyncJobStatusStopped   SyncJobStatus = "stopped"
)

// IsValid checks if the sync job status is valid
func (s SyncJobStatus) IsValid() bool {
	switch s {
	case SyncJobStatusPending, SyncJobStatusRunning, SyncJobStatusCompleted,
		SyncJobStatusFailed, SyncJobStatusStopped:
		return true
	}
	return false
}

// String returns the string representation of SyncJobStatus
func (s SyncJobStatus) String() string {
	return string(s)
}

// SyncJob represents a sync job execution in the database
type SyncJob struct {
	ID           int64          `json:"id" db:"id"`
	TaskID       int64          `json:"task_id" db:"task_id"`
	UserID       int64          `json:"user_id" db:"user_id"`
	Status       SyncJobStatus  `json:"status" db:"status"`
	Progress     int            `json:"progress" db:"progress"`
	ErrorMessage sql.NullString `json:"error_message" db:"error_message"`
	StartedAt    sql.NullTime   `json:"started_at" db:"started_at"`
	FinishedAt   sql.NullTime   `json:"finished_at" db:"finished_at"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
}

// SyncJobResponse is the API response for sync job
type SyncJobResponse struct {
	ID           int64         `json:"id"`
	TaskID       int64         `json:"task_id"`
	UserID       int64         `json:"user_id"`
	Status       SyncJobStatus `json:"status"`
	Progress     int           `json:"progress"`
	ErrorMessage string        `json:"error_message,omitempty"`
	StartedAt    *time.Time    `json:"started_at,omitempty"`
	FinishedAt   *time.Time    `json:"finished_at,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// ToResponse converts SyncJob to SyncJobResponse
func (sj *SyncJob) ToResponse() SyncJobResponse {
	resp := SyncJobResponse{
		ID:        sj.ID,
		TaskID:    sj.TaskID,
		UserID:    sj.UserID,
		Status:    sj.Status,
		Progress:  sj.Progress,
		CreatedAt: sj.CreatedAt,
		UpdatedAt: sj.UpdatedAt,
	}

	if sj.ErrorMessage.Valid {
		resp.ErrorMessage = sj.ErrorMessage.String
	}
	if sj.StartedAt.Valid {
		resp.StartedAt = &sj.StartedAt.Time
	}
	if sj.FinishedAt.Valid {
		resp.FinishedAt = &sj.FinishedAt.Time
	}

	return resp
}

// SyncJobListResponse represents the response for listing sync jobs
type SyncJobListResponse struct {
	SyncJobs []SyncJobResponse `json:"syncjobs"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

// ErrSyncJobNotFound is returned when a sync job is not found
var ErrSyncJobNotFound = errors.New("sync job not found")

// ErrSyncJobAlreadyRunning is returned when a sync job is already running for the task
var ErrSyncJobAlreadyRunning = errors.New("a sync job is already running for this task")

// ErrSyncTaskNotReady is returned when the sync task is not ready to run
var ErrSyncTaskNotReady = errors.New("sync task is not ready to run")
