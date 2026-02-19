package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

// SyncTaskStatus represents the status of a sync task
type SyncTaskStatus string

// Supported sync task statuses
const (
	SyncTaskStatusDraft     SyncTaskStatus = "draft"
	SyncTaskStatusReady     SyncTaskStatus = "ready"
	SyncTaskStatusRunning   SyncTaskStatus = "running"
	SyncTaskStatusPaused    SyncTaskStatus = "paused"
	SyncTaskStatusCompleted SyncTaskStatus = "completed"
	SyncTaskStatusFailed    SyncTaskStatus = "failed"
)

// IsValid checks if the sync task status is valid
func (s SyncTaskStatus) IsValid() bool {
	switch s {
	case SyncTaskStatusDraft, SyncTaskStatusReady, SyncTaskStatusRunning,
		SyncTaskStatusPaused, SyncTaskStatusCompleted, SyncTaskStatusFailed:
		return true
	}
	return false
}

// String returns the string representation of SyncTaskStatus
func (s SyncTaskStatus) String() string {
	return string(s)
}

// SyncTask represents a sync task in the database
type SyncTask struct {
	ID           int64          `json:"id" db:"id"`
	UserID       int64          `json:"user_id" db:"user_id"`
	Name         string         `json:"name" db:"name"`
	ReaderConfig sql.NullString `json:"reader_config" db:"reader_config"`
	WriterConfig sql.NullString `json:"writer_config" db:"writer_config"`
	Status       SyncTaskStatus `json:"status" db:"status"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
}

// SyncTaskResponse is the API response for sync task
type SyncTaskResponse struct {
	ID           int64           `json:"id"`
	UserID       int64           `json:"user_id"`
	Name         string          `json:"name"`
	ReaderConfig json.RawMessage `json:"reader_config,omitempty"`
	WriterConfig json.RawMessage `json:"writer_config,omitempty"`
	Status       SyncTaskStatus  `json:"status"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// ToResponse converts SyncTask to SyncTaskResponse
func (st *SyncTask) ToResponse() SyncTaskResponse {
	resp := SyncTaskResponse{
		ID:        st.ID,
		UserID:    st.UserID,
		Name:      st.Name,
		Status:    st.Status,
		CreatedAt: st.CreatedAt,
		UpdatedAt: st.UpdatedAt,
	}

	if st.ReaderConfig.Valid {
		resp.ReaderConfig = json.RawMessage(st.ReaderConfig.String)
	}
	if st.WriterConfig.Valid {
		resp.WriterConfig = json.RawMessage(st.WriterConfig.String)
	}

	return resp
}

// CreateSyncTaskRequest represents the request to create a sync task
type CreateSyncTaskRequest struct {
	Name         string          `json:"name" binding:"required,min=1,max=100"`
	ReaderConfig json.RawMessage `json:"reader_config" binding:"required"`
	WriterConfig json.RawMessage `json:"writer_config" binding:"required"`
}

// UpdateSyncTaskRequest represents the request to update a sync task
type UpdateSyncTaskRequest struct {
	Name         string          `json:"name" binding:"omitempty,min=1,max=100"`
	ReaderConfig json.RawMessage `json:"reader_config"`
	WriterConfig json.RawMessage `json:"writer_config"`
	Status       SyncTaskStatus  `json:"status"`
}

// SyncTaskListResponse represents the response for listing sync tasks
type SyncTaskListResponse struct {
	SyncTasks []SyncTaskResponse `json:"synctasks"`
	Total     int                `json:"total"`
	Page      int                `json:"page"`
	Limit     int                `json:"limit"`
}

// ErrSyncTaskNotFound is returned when a sync task is not found
var ErrSyncTaskNotFound = errors.New("sync task not found")

// ErrSyncTaskExists is returned when a sync task with the same name already exists
var ErrSyncTaskExists = errors.New("sync task with this name already exists")

// ErrInvalidSyncTaskStatus is returned when an invalid sync task status is provided
var ErrInvalidSyncTaskStatus = errors.New("invalid sync task status")
