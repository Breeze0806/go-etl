package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Breeze0806/go-etl/cmd/web/middleware"
	"github.com/Breeze0806/go-etl/cmd/web/models"
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/core/job"
	"github.com/gin-gonic/gin"
)

type SyncJobHandler struct {
	db          *sql.DB
	jwt         *middleware.JWTMiddleware
	runningJobs sync.Map
}

func NewSyncJobHandler(db *sql.DB, jwt *middleware.JWTMiddleware) *SyncJobHandler {
	return &SyncJobHandler{
		db:  db,
		jwt: jwt,
	}
}

func (h *SyncJobHandler) Start(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sync task id"})
		return
	}

	var st models.SyncTask
	err = h.db.QueryRow(
		`SELECT id, user_id, name, reader_config, writer_config, status, created_at, updated_at 
		FROM sync_tasks WHERE id = ?`,
		taskID,
	).Scan(&st.ID, &st.UserID, &st.Name, &st.ReaderConfig, &st.WriterConfig, &st.Status, &st.CreatedAt, &st.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "sync task not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if st.UserID != userIDInt {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	if st.Status != models.SyncTaskStatusReady {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sync task is not ready to run"})
		return
	}

	var existingJob models.SyncJob
	err = h.db.QueryRow(
		`SELECT id, task_id, user_id, status, progress, error_message, started_at, finished_at, created_at, updated_at 
		FROM sync_jobs WHERE task_id = ? AND status = ?`,
		taskID, models.SyncJobStatusRunning,
	).Scan(&existingJob.ID, &existingJob.TaskID, &existingJob.UserID, &existingJob.Status, &existingJob.Progress,
		&existingJob.ErrorMessage, &existingJob.StartedAt, &existingJob.FinishedAt, &existingJob.CreatedAt, &existingJob.UpdatedAt)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "a sync job is already running for this task"})
		return
	}
	if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	now := time.Now()
	result, err := h.db.Exec(
		`INSERT INTO sync_jobs (task_id, user_id, status, progress, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		taskID, userIDInt, models.SyncJobStatusPending, 0, now, now,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create sync job"})
		return
	}

	jobID, _ := result.LastInsertId()

	go h.runSyncJob(jobID, taskID, userIDInt, st.ReaderConfig.String, st.WriterConfig.String)

	c.JSON(http.StatusAccepted, gin.H{"id": jobID, "message": "sync job started"})
}

func (h *SyncJobHandler) runSyncJob(jobID, taskID, userID int64, readerConfig, writerConfig string) {
	ctx, cancel := context.WithCancel(context.Background())
	h.runningJobs.Store(jobID, cancel)

	defer func() {
		h.runningJobs.Delete(jobID)
	}()

	now := time.Now()
	_, err := h.db.Exec(
		`UPDATE sync_jobs SET status = ?, started_at = ?, updated_at = ? WHERE id = ?`,
		models.SyncJobStatusRunning, now, now, jobID,
	)
	if err != nil {
		h.updateJobError(jobID, fmt.Sprintf("failed to update job status: %v", err))
		return
	}

	cfg, err := h.buildJobConfig(readerConfig, writerConfig, jobID)
	if err != nil {
		h.updateJobError(jobID, fmt.Sprintf("failed to build job config: %v", err))
		return
	}

	container, err := job.NewContainer(ctx, cfg)
	if err != nil {
		h.updateJobError(jobID, fmt.Sprintf("failed to create job container: %v", err))
		return
	}

	err = container.Start()
	if err != nil {
		h.updateJobError(jobID, fmt.Sprintf("job execution failed: %v", err))
		return
	}

	finishTime := time.Now()
	_, err = h.db.Exec(
		`UPDATE sync_jobs SET status = ?, progress = 100, finished_at = ?, updated_at = ? WHERE id = ?`,
		models.SyncJobStatusCompleted, finishTime, finishTime, jobID,
	)
	if err != nil {
		return
	}

	_, err = h.db.Exec(
		`UPDATE sync_tasks SET status = ?, updated_at = ? WHERE id = ?`,
		models.SyncTaskStatusCompleted, finishTime, taskID,
	)
}

func (h *SyncJobHandler) buildJobConfig(readerConfig, writerConfig string, jobID int64) (*config.JSON, error) {
	jobConfig := map[string]interface{}{
		"core": map[string]interface{}{
			"container": map[string]interface{}{
				"job": map[string]interface{}{
					"id":            jobID,
					"sleepInterval": 100,
				},
				"taskGroup": map[string]interface{}{
					"id": 1,
					"failover": map[string]interface{}{
						"retryIntervalInMsec": 0,
					},
				},
			},
			"transport": map[string]interface{}{
				"channel": map[string]interface{}{
					"speed": map[string]interface{}{
						"byte":   1048576,
						"record": 500,
					},
				},
			},
		},
		"job": map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"reader":      json.RawMessage(readerConfig),
					"writer":      json.RawMessage(writerConfig),
					"transformer": []interface{}{},
				},
			},
			"setting": map[string]interface{}{
				"speed": map[string]interface{}{
					"byte":    1048576,
					"record":  500,
					"channel": 1,
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(jobConfig)
	if err != nil {
		return nil, err
	}

	return config.NewJSONFromString(string(jsonBytes))
}

func (h *SyncJobHandler) updateJobError(jobID int64, errMsg string) {
	finishTime := time.Now()
	h.db.Exec(
		`UPDATE sync_jobs SET status = ?, error_message = ?, finished_at = ?, updated_at = ? WHERE id = ?`,
		models.SyncJobStatusFailed, errMsg, finishTime, finishTime, jobID,
	)
}

func (h *SyncJobHandler) GetByID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	jobIDStr := c.Param("id")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sync job id"})
		return
	}

	var sj models.SyncJob
	err = h.db.QueryRow(
		`SELECT id, task_id, user_id, status, progress, error_message, started_at, finished_at, created_at, updated_at 
		FROM sync_jobs WHERE id = ?`,
		jobID,
	).Scan(&sj.ID, &sj.TaskID, &sj.UserID, &sj.Status, &sj.Progress, &sj.ErrorMessage,
		&sj.StartedAt, &sj.FinishedAt, &sj.CreatedAt, &sj.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "sync job not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if sj.UserID != userIDInt {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	c.JSON(http.StatusOK, sj.ToResponse())
}

func (h *SyncJobHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	taskIDStr := c.Query("task_id")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var query string
	var args []interface{}

	if taskIDStr != "" && status != "" {
		taskID, _ := strconv.ParseInt(taskIDStr, 10, 64)
		query = `
			SELECT id, task_id, user_id, status, progress, error_message, started_at, finished_at, created_at, updated_at 
			FROM sync_jobs 
			WHERE user_id = ? AND task_id = ? AND status = ?
			ORDER BY id DESC 
			LIMIT ? OFFSET ?
		`
		args = []interface{}{userIDInt, taskID, status, limit, offset}
	} else if taskIDStr != "" {
		taskID, _ := strconv.ParseInt(taskIDStr, 10, 64)
		query = `
			SELECT id, task_id, user_id, status, progress, error_message, started_at, finished_at, created_at, updated_at 
			FROM sync_jobs 
			WHERE user_id = ? AND task_id = ?
			ORDER BY id DESC 
			LIMIT ? OFFSET ?
		`
		args = []interface{}{userIDInt, taskID, limit, offset}
	} else if status != "" {
		query = `
			SELECT id, task_id, user_id, status, progress, error_message, started_at, finished_at, created_at, updated_at 
			FROM sync_jobs 
			WHERE user_id = ? AND status = ?
			ORDER BY id DESC 
			LIMIT ? OFFSET ?
		`
		args = []interface{}{userIDInt, status, limit, offset}
	} else {
		query = `
			SELECT id, task_id, user_id, status, progress, error_message, started_at, finished_at, created_at, updated_at 
			FROM sync_jobs 
			WHERE user_id = ?
			ORDER BY id DESC 
			LIMIT ? OFFSET ?
		`
		args = []interface{}{userIDInt, limit, offset}
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query sync jobs"})
		return
	}
	defer rows.Close()

	var syncJobs []models.SyncJobResponse
	for rows.Next() {
		var sj models.SyncJob
		err := rows.Scan(&sj.ID, &sj.TaskID, &sj.UserID, &sj.Status, &sj.Progress, &sj.ErrorMessage,
			&sj.StartedAt, &sj.FinishedAt, &sj.CreatedAt, &sj.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan sync job"})
			return
		}
		syncJobs = append(syncJobs, sj.ToResponse())
	}

	if syncJobs == nil {
		syncJobs = []models.SyncJobResponse{}
	}

	var total int
	if taskIDStr != "" && status != "" {
		taskID, _ := strconv.ParseInt(taskIDStr, 10, 64)
		err = h.db.QueryRow("SELECT COUNT(*) FROM sync_jobs WHERE user_id = ? AND task_id = ? AND status = ?",
			userIDInt, taskID, status).Scan(&total)
	} else if taskIDStr != "" {
		taskID, _ := strconv.ParseInt(taskIDStr, 10, 64)
		err = h.db.QueryRow("SELECT COUNT(*) FROM sync_jobs WHERE user_id = ? AND task_id = ?",
			userIDInt, taskID).Scan(&total)
	} else if status != "" {
		err = h.db.QueryRow("SELECT COUNT(*) FROM sync_jobs WHERE user_id = ? AND status = ?",
			userIDInt, status).Scan(&total)
	} else {
		err = h.db.QueryRow("SELECT COUNT(*) FROM sync_jobs WHERE user_id = ?", userIDInt).Scan(&total)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count sync jobs"})
		return
	}

	c.JSON(http.StatusOK, models.SyncJobListResponse{
		SyncJobs: syncJobs,
		Total:    total,
		Page:     page,
		Limit:    limit,
	})
}

func (h *SyncJobHandler) Stop(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	jobIDStr := c.Param("id")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sync job id"})
		return
	}

	var sj models.SyncJob
	err = h.db.QueryRow(
		`SELECT id, task_id, user_id, status, progress, error_message, started_at, finished_at, created_at, updated_at 
		FROM sync_jobs WHERE id = ?`,
		jobID,
	).Scan(&sj.ID, &sj.TaskID, &sj.UserID, &sj.Status, &sj.Progress, &sj.ErrorMessage,
		&sj.StartedAt, &sj.FinishedAt, &sj.CreatedAt, &sj.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "sync job not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if sj.UserID != userIDInt {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	if sj.Status != models.SyncJobStatusRunning {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sync job is not running"})
		return
	}

	if cancel, ok := h.runningJobs.Load(jobID); ok {
		cancel.(context.CancelFunc)()
	}

	finishTime := time.Now()
	_, err = h.db.Exec(
		`UPDATE sync_jobs SET status = ?, finished_at = ?, updated_at = ? WHERE id = ?`,
		models.SyncJobStatusStopped, finishTime, finishTime, jobID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to stop sync job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sync job stopped successfully"})
}
