package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/Breeze0806/go-etl/cmd/web/middleware"
	"github.com/Breeze0806/go-etl/cmd/web/models"
	"github.com/gin-gonic/gin"
)

type SyncTaskHandler struct {
	db  *sql.DB
	jwt *middleware.JWTMiddleware
}

func NewSyncTaskHandler(db *sql.DB, jwt *middleware.JWTMiddleware) *SyncTaskHandler {
	return &SyncTaskHandler{
		db:  db,
		jwt: jwt,
	}
}

func (h *SyncTaskHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	name := c.Query("name")
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

	if name != "" && status != "" {
		query = `
			SELECT id, user_id, name, reader_config, writer_config, status, created_at, updated_at 
			FROM sync_tasks 
			WHERE user_id = ? AND name LIKE ? AND status = ?
			ORDER BY id DESC 
			LIMIT ? OFFSET ?
		`
		args = []interface{}{userIDInt, "%" + name + "%", status, limit, offset}
	} else if name != "" {
		query = `
			SELECT id, user_id, name, reader_config, writer_config, status, created_at, updated_at 
			FROM sync_tasks 
			WHERE user_id = ? AND name LIKE ?
			ORDER BY id DESC 
			LIMIT ? OFFSET ?
		`
		args = []interface{}{userIDInt, "%" + name + "%", limit, offset}
	} else if status != "" {
		query = `
			SELECT id, user_id, name, reader_config, writer_config, status, created_at, updated_at 
			FROM sync_tasks 
			WHERE user_id = ? AND status = ?
			ORDER BY id DESC 
			LIMIT ? OFFSET ?
		`
		args = []interface{}{userIDInt, status, limit, offset}
	} else {
		query = `
			SELECT id, user_id, name, reader_config, writer_config, status, created_at, updated_at 
			FROM sync_tasks 
			WHERE user_id = ?
			ORDER BY id DESC 
			LIMIT ? OFFSET ?
		`
		args = []interface{}{userIDInt, limit, offset}
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query sync tasks"})
		return
	}
	defer rows.Close()

	var syncTasks []models.SyncTaskResponse
	for rows.Next() {
		var st models.SyncTask
		err := rows.Scan(&st.ID, &st.UserID, &st.Name, &st.ReaderConfig,
			&st.WriterConfig, &st.Status, &st.CreatedAt, &st.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan sync task"})
			return
		}
		syncTasks = append(syncTasks, st.ToResponse())
	}

	if syncTasks == nil {
		syncTasks = []models.SyncTaskResponse{}
	}

	var total int
	if name != "" && status != "" {
		err = h.db.QueryRow("SELECT COUNT(*) FROM sync_tasks WHERE user_id = ? AND name LIKE ? AND status = ?",
			userIDInt, "%"+name+"%", status).Scan(&total)
	} else if name != "" {
		err = h.db.QueryRow("SELECT COUNT(*) FROM sync_tasks WHERE user_id = ? AND name LIKE ?",
			userIDInt, "%"+name+"%").Scan(&total)
	} else if status != "" {
		err = h.db.QueryRow("SELECT COUNT(*) FROM sync_tasks WHERE user_id = ? AND status = ?",
			userIDInt, status).Scan(&total)
	} else {
		err = h.db.QueryRow("SELECT COUNT(*) FROM sync_tasks WHERE user_id = ?", userIDInt).Scan(&total)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count sync tasks"})
		return
	}

	c.JSON(http.StatusOK, models.SyncTaskListResponse{
		SyncTasks: syncTasks,
		Total:     total,
		Page:      page,
		Limit:     limit,
	})
}

func (h *SyncTaskHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	var req models.CreateSyncTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()

	result, err := h.db.Exec(
		`INSERT INTO sync_tasks (user_id, name, reader_config, writer_config, status, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userIDInt, req.Name, nullString(string(req.ReaderConfig)), nullString(string(req.WriterConfig)),
		models.SyncTaskStatusDraft, now, now,
	)
	if err != nil {
		if isDuplicateErrorSyncTask(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "sync task with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create sync task"})
		return
	}

	id, _ := result.LastInsertId()

	var st models.SyncTask
	err = h.db.QueryRow(
		`SELECT id, user_id, name, reader_config, writer_config, status, created_at, updated_at 
		FROM sync_tasks WHERE id = ?`,
		id,
	).Scan(&st.ID, &st.UserID, &st.Name, &st.ReaderConfig, &st.WriterConfig, &st.Status, &st.CreatedAt, &st.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch created sync task"})
		return
	}

	c.JSON(http.StatusCreated, st.ToResponse())
}

func (h *SyncTaskHandler) Update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	stIDStr := c.Param("id")
	stID, err := strconv.ParseInt(stIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sync task id"})
		return
	}

	var existingST models.SyncTask
	err = h.db.QueryRow(
		`SELECT id, user_id, name, reader_config, writer_config, status, created_at, updated_at 
		FROM sync_tasks WHERE id = ?`,
		stID,
	).Scan(&existingST.ID, &existingST.UserID, &existingST.Name, &existingST.ReaderConfig,
		&existingST.WriterConfig, &existingST.Status, &existingST.CreatedAt, &existingST.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "sync task not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if existingST.UserID != userIDInt {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	var req models.UpdateSyncTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status != "" && !req.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sync task status"})
		return
	}

	name := existingST.Name
	if req.Name != "" {
		name = req.Name
	}

	readerConfig := existingST.ReaderConfig
	if len(req.ReaderConfig) > 0 {
		readerConfig = nullString(string(req.ReaderConfig))
	}

	writerConfig := existingST.WriterConfig
	if len(req.WriterConfig) > 0 {
		writerConfig = nullString(string(req.WriterConfig))
	}

	status := existingST.Status
	if req.Status != "" {
		status = req.Status
	}

	now := time.Now()

	_, err = h.db.Exec(
		`UPDATE sync_tasks SET name = ?, reader_config = ?, writer_config = ?, status = ?, updated_at = ? WHERE id = ?`,
		name, readerConfig, writerConfig, status, now, stID,
	)
	if err != nil {
		if isDuplicateErrorSyncTask(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "sync task with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update sync task"})
		return
	}

	var updatedST models.SyncTask
	err = h.db.QueryRow(
		`SELECT id, user_id, name, reader_config, writer_config, status, created_at, updated_at 
		FROM sync_tasks WHERE id = ?`,
		stID,
	).Scan(&updatedST.ID, &updatedST.UserID, &updatedST.Name, &updatedST.ReaderConfig,
		&updatedST.WriterConfig, &updatedST.Status, &updatedST.CreatedAt, &updatedST.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated sync task"})
		return
	}

	c.JSON(http.StatusOK, updatedST.ToResponse())
}

func (h *SyncTaskHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	stIDStr := c.Param("id")
	stID, err := strconv.ParseInt(stIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sync task id"})
		return
	}

	var existingST models.SyncTask
	err = h.db.QueryRow(
		`SELECT id, user_id FROM sync_tasks WHERE id = ?`,
		stID,
	).Scan(&existingST.ID, &existingST.UserID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "sync task not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if existingST.UserID != userIDInt {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	result, err := h.db.Exec("DELETE FROM sync_tasks WHERE id = ?", stID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete sync task"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "sync task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sync task deleted successfully"})
}

func (h *SyncTaskHandler) GetByID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	stIDStr := c.Param("id")
	stID, err := strconv.ParseInt(stIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sync task id"})
		return
	}

	var st models.SyncTask
	err = h.db.QueryRow(
		`SELECT id, user_id, name, reader_config, writer_config, status, created_at, updated_at 
		FROM sync_tasks WHERE id = ?`,
		stID,
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

	c.JSON(http.StatusOK, st.ToResponse())
}

func isDuplicateErrorSyncTask(err error) bool {
	errStr := err.Error()
	return errStr == "UNIQUE constraint failed: sync_tasks.name" ||
		errStr == "UNIQUE constraint failed: sync_tasks.user_id, sync_tasks.name"
}
