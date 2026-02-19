package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "gitee.com/chunanyong/dm"
	"github.com/Breeze0806/go-etl/cmd/web/middleware"
	"github.com/Breeze0806/go-etl/cmd/web/models"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/ibmdb/go_ibm_db"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"

	"github.com/gin-gonic/gin"
)

type DataSourceHandler struct {
	db  *sql.DB
	jwt *middleware.JWTMiddleware
}

func NewDataSourceHandler(db *sql.DB, jwt *middleware.JWTMiddleware) *DataSourceHandler {
	return &DataSourceHandler{
		db:  db,
		jwt: jwt,
	}
}

func (h *DataSourceHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	name := c.Query("name")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var query string
	var args []interface{}

	if name != "" {
		query = `
			SELECT id, user_id, name, type, host, port, username, password, database, file_path, created_at, updated_at 
			FROM data_sources 
			WHERE user_id = ? AND name LIKE ? 
			ORDER BY id DESC 
			LIMIT ? OFFSET ?
		`
		args = []interface{}{userIDInt, "%" + name + "%", limit, offset}
	} else {
		query = `
			SELECT id, user_id, name, type, host, port, username, password, database, file_path, created_at, updated_at 
			FROM data_sources 
			WHERE user_id = ?
			ORDER BY id DESC 
			LIMIT ? OFFSET ?
		`
		args = []interface{}{userIDInt, limit, offset}
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query data sources"})
		return
	}
	defer rows.Close()

	var dataSources []models.DataSourceResponse
	for rows.Next() {
		var ds models.DataSource
		err := rows.Scan(&ds.ID, &ds.UserID, &ds.Name, &ds.Type, &ds.Host, &ds.Port,
			&ds.Username, &ds.Password, &ds.Database, &ds.FilePath, &ds.CreatedAt, &ds.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan data source"})
			return
		}
		dataSources = append(dataSources, ds.ToResponse())
	}

	if dataSources == nil {
		dataSources = []models.DataSourceResponse{}
	}

	var total int
	if name != "" {
		err = h.db.QueryRow("SELECT COUNT(*) FROM data_sources WHERE user_id = ? AND name LIKE ?",
			userIDInt, "%"+name+"%").Scan(&total)
	} else {
		err = h.db.QueryRow("SELECT COUNT(*) FROM data_sources WHERE user_id = ?", userIDInt).Scan(&total)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count data sources"})
		return
	}

	c.JSON(http.StatusOK, models.DataSourceListResponse{
		DataSources: dataSources,
		Total:       total,
		Page:        page,
		Limit:       limit,
	})
}

func (h *DataSourceHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	var req models.CreateDataSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !req.Type.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data source type"})
		return
	}

	now := time.Now()

	result, err := h.db.Exec(
		`INSERT INTO data_sources (user_id, name, type, host, port, username, password, database, file_path, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userIDInt, req.Name, req.Type, nullString(req.Host), nullInt64(req.Port),
		nullString(req.Username), nullString(req.Password), nullString(req.Database),
		nullString(req.FilePath), now, now,
	)
	if err != nil {
		if isDuplicateErrorDataSource(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "data source with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create data source"})
		return
	}

	id, _ := result.LastInsertId()

	var ds models.DataSource
	err = h.db.QueryRow(
		`SELECT id, user_id, name, type, host, port, username, password, database, file_path, created_at, updated_at 
		FROM data_sources WHERE id = ?`,
		id,
	).Scan(&ds.ID, &ds.UserID, &ds.Name, &ds.Type, &ds.Host, &ds.Port,
		&ds.Username, &ds.Password, &ds.Database, &ds.FilePath, &ds.CreatedAt, &ds.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch created data source"})
		return
	}

	c.JSON(http.StatusCreated, ds.ToResponse())
}

func (h *DataSourceHandler) Update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	dsIDStr := c.Param("id")
	dsID, err := strconv.ParseInt(dsIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data source id"})
		return
	}

	var existingDS models.DataSource
	err = h.db.QueryRow(
		`SELECT id, user_id, name, type, host, port, username, password, database, file_path, created_at, updated_at 
		FROM data_sources WHERE id = ?`,
		dsID,
	).Scan(&existingDS.ID, &existingDS.UserID, &existingDS.Name, &existingDS.Type,
		&existingDS.Host, &existingDS.Port, &existingDS.Username, &existingDS.Password,
		&existingDS.Database, &existingDS.FilePath, &existingDS.CreatedAt, &existingDS.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "data source not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if existingDS.UserID != userIDInt {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	var req models.UpdateDataSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Type != "" && !req.Type.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data source type"})
		return
	}

	name := existingDS.Name
	if req.Name != "" {
		name = req.Name
	}

	dsType := existingDS.Type
	if req.Type != "" {
		dsType = req.Type
	}

	host := existingDS.Host
	if req.Host != "" {
		host = sql.NullString{String: req.Host, Valid: true}
	}

	port := existingDS.Port
	if req.Port > 0 {
		port = sql.NullInt64{Int64: int64(req.Port), Valid: true}
	}

	username := existingDS.Username
	if req.Username != "" {
		username = sql.NullString{String: req.Username, Valid: true}
	}

	password := existingDS.Password
	if req.Password != "" {
		password = sql.NullString{String: req.Password, Valid: true}
	}

	database := existingDS.Database
	if req.Database != "" {
		database = sql.NullString{String: req.Database, Valid: true}
	}

	filePath := existingDS.FilePath
	if req.FilePath != "" {
		filePath = sql.NullString{String: req.FilePath, Valid: true}
	}

	now := time.Now()

	_, err = h.db.Exec(
		`UPDATE data_sources SET name = ?, type = ?, host = ?, port = ?, username = ?, 
		password = ?, database = ?, file_path = ?, updated_at = ? WHERE id = ?`,
		name, dsType, host, port, username, password, database, filePath, now, dsID,
	)
	if err != nil {
		if isDuplicateErrorDataSource(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "data source with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update data source"})
		return
	}

	var updatedDS models.DataSource
	err = h.db.QueryRow(
		`SELECT id, user_id, name, type, host, port, username, password, database, file_path, created_at, updated_at 
		FROM data_sources WHERE id = ?`,
		dsID,
	).Scan(&updatedDS.ID, &updatedDS.UserID, &updatedDS.Name, &updatedDS.Type,
		&updatedDS.Host, &updatedDS.Port, &updatedDS.Username, &updatedDS.Password,
		&updatedDS.Database, &updatedDS.FilePath, &updatedDS.CreatedAt, &updatedDS.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated data source"})
		return
	}

	c.JSON(http.StatusOK, updatedDS.ToResponse())
}

func (h *DataSourceHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	dsIDStr := c.Param("id")
	dsID, err := strconv.ParseInt(dsIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data source id"})
		return
	}

	var existingDS models.DataSource
	err = h.db.QueryRow(
		`SELECT id, user_id FROM data_sources WHERE id = ?`,
		dsID,
	).Scan(&existingDS.ID, &existingDS.UserID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "data source not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if existingDS.UserID != userIDInt {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	result, err := h.db.Exec("DELETE FROM data_sources WHERE id = ?", dsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete data source"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "data source not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "data source deleted successfully"})
}

func (h *DataSourceHandler) GetByID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt := userID.(int64)

	dsIDStr := c.Param("id")
	dsID, err := strconv.ParseInt(dsIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data source id"})
		return
	}

	var ds models.DataSource
	err = h.db.QueryRow(
		`SELECT id, user_id, name, type, host, port, username, password, database, file_path, created_at, updated_at 
		FROM data_sources WHERE id = ?`,
		dsID,
	).Scan(&ds.ID, &ds.UserID, &ds.Name, &ds.Type, &ds.Host, &ds.Port,
		&ds.Username, &ds.Password, &ds.Database, &ds.FilePath, &ds.CreatedAt, &ds.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "data source not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if ds.UserID != userIDInt {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	c.JSON(http.StatusOK, ds.ToResponse())
}

func (h *DataSourceHandler) TestConnection(c *gin.Context) {
	var req models.TestConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !req.Type.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data source type"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dsn, driverName, err := buildConnectionString(req.Type, req.Host, req.Port, req.Username, req.Password, req.Database, req.FilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		c.JSON(http.StatusOK, models.TestConnectionResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to connect: %v", err),
		})
		return
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		c.JSON(http.StatusOK, models.TestConnectionResponse{
			Success: false,
			Message: fmt.Sprintf("Connection failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, models.TestConnectionResponse{
		Success: true,
		Message: "Connection successful",
	})
}

func buildConnectionString(dsType models.DataSourceType, host string, port int, username, password, database, filePath string) (string, string, error) {
	switch dsType {
	case models.DataSourceTypeMySQL:
		if port == 0 {
			port = 3306
		}
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=5s",
			username, password, host, port, database)
		return dsn, "mysql", nil

	case models.DataSourceTypePostgres:
		if port == 0 {
			port = 5432
		}
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=5",
			host, port, username, password, database)
		return dsn, "postgres", nil

	case models.DataSourceTypeOracle:
		if port == 0 {
			port = 1521
		}
		dsn := fmt.Sprintf("%s/%s@%s:%d/%s",
			username, password, host, port, database)
		return dsn, "oracle", nil

	case models.DataSourceTypeDB2:
		if port == 0 {
			port = 50000
		}
		dsn := fmt.Sprintf("HOSTNAME=%s;PORT=%d;DATABASE=%s;UID=%s;PWD=%s;",
			host, port, database, username, password)
		return dsn, "db2", nil

	case models.DataSourceTypeSQLite3:
		dsn := filePath
		if dsn == "" {
			dsn = database
		}
		return dsn, "sqlite3", nil

	case models.DataSourceTypeSQLServer:
		if port == 0 {
			port = 1433
		}
		dsn := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s;encrypt=true;connection timeout=5",
			host, port, database, username, password)
		return dsn, "mssql", nil

	case models.DataSourceTypeDameng:
		if port == 0 {
			port = 5236
		}
		dsn := fmt.Sprintf("dm://%s:%s@%s:%d",
			username, password, host, port)
		return dsn, "dm", nil

	case models.DataSourceTypeCSV, models.DataSourceTypeXLSX:
		return filePath, "sqlite3", nil

	default:
		return "", "", fmt.Errorf("unsupported data source type: %s", dsType)
	}
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullInt64(i int) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(i), Valid: true}
}

func isDuplicateErrorDataSource(err error) bool {
	errStr := err.Error()
	return errStr == "UNIQUE constraint failed: data_sources.name" ||
		errStr == "UNIQUE constraint failed: data_sources.user_id, data_sources.name"
}
