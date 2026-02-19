package models

import (
	"database/sql"
	"errors"
	"time"
)

// DataSourceType represents the type of data source
type DataSourceType string

// Supported data source types
const (
	DataSourceTypeMySQL     DataSourceType = "mysql"
	DataSourceTypePostgres  DataSourceType = "postgres"
	DataSourceTypeOracle    DataSourceType = "oracle"
	DataSourceTypeDB2       DataSourceType = "db2"
	DataSourceTypeSQLite3   DataSourceType = "sqlite3"
	DataSourceTypeSQLServer DataSourceType = "sqlserver"
	DataSourceTypeDameng    DataSourceType = "dameng"
	DataSourceTypeCSV       DataSourceType = "csv"
	DataSourceTypeXLSX      DataSourceType = "xlsx"
)

// IsValid checks if the data source type is valid
func (t DataSourceType) IsValid() bool {
	switch t {
	case DataSourceTypeMySQL, DataSourceTypePostgres, DataSourceTypeOracle,
		DataSourceTypeDB2, DataSourceTypeSQLite3, DataSourceTypeSQLServer,
		DataSourceTypeDameng, DataSourceTypeCSV, DataSourceTypeXLSX:
		return true
	}
	return false
}

// String returns the string representation of DataSourceType
func (t DataSourceType) String() string {
	return string(t)
}

// DataSource represents a data source in the database
type DataSource struct {
	ID        int64          `json:"id" db:"id"`
	UserID    int64          `json:"user_id" db:"user_id"`
	Name      string         `json:"name" db:"name"`
	Type      DataSourceType `json:"type" db:"type"`
	Host      sql.NullString `json:"host" db:"host"`
	Port      sql.NullInt64  `json:"port" db:"port"`
	Username  sql.NullString `json:"username" db:"username"`
	Password  sql.NullString `json:"-" db:"password"` // Never expose in JSON
	Database  sql.NullString `json:"database" db:"database"`
	FilePath  sql.NullString `json:"file_path" db:"file_path"` // For CSV/XLSX
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// DataSourceResponse is the API response for data source (without password)
type DataSourceResponse struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	Name      string         `json:"name"`
	Type      DataSourceType `json:"type"`
	Host      string         `json:"host,omitempty"`
	Port      int64          `json:"port,omitempty"`
	Username  string         `json:"username,omitempty"`
	Database  string         `json:"database,omitempty"`
	FilePath  string         `json:"file_path,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// ToResponse converts DataSource to DataSourceResponse
func (ds *DataSource) ToResponse() DataSourceResponse {
	resp := DataSourceResponse{
		ID:        ds.ID,
		UserID:    ds.UserID,
		Name:      ds.Name,
		Type:      ds.Type,
		CreatedAt: ds.CreatedAt,
		UpdatedAt: ds.UpdatedAt,
	}

	if ds.Host.Valid {
		resp.Host = ds.Host.String
	}
	if ds.Port.Valid {
		resp.Port = ds.Port.Int64
	}
	if ds.Username.Valid {
		resp.Username = ds.Username.String
	}
	if ds.Database.Valid {
		resp.Database = ds.Database.String
	}
	if ds.FilePath.Valid {
		resp.FilePath = ds.FilePath.String
	}

	return resp
}

// CreateDataSourceRequest represents the request to create a data source
type CreateDataSourceRequest struct {
	Name     string         `json:"name" binding:"required,min=1,max=100"`
	Type     DataSourceType `json:"type" binding:"required"`
	Host     string         `json:"host"`
	Port     int            `json:"port"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Database string         `json:"database"`
	FilePath string         `json:"file_path"` // For CSV/XLSX
}

// UpdateDataSourceRequest represents the request to update a data source
type UpdateDataSourceRequest struct {
	Name     string         `json:"name" binding:"omitempty,min=1,max=100"`
	Type     DataSourceType `json:"type"`
	Host     string         `json:"host"`
	Port     int            `json:"port"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Database string         `json:"database"`
	FilePath string         `json:"file_path"`
}

// DataSourceListResponse represents the response for listing data sources
type DataSourceListResponse struct {
	DataSources []DataSourceResponse `json:"datasources"`
	Total       int                  `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
}

// TestConnectionRequest represents the request to test a data source connection
type TestConnectionRequest struct {
	Type     DataSourceType `json:"type" binding:"required"`
	Host     string         `json:"host"`
	Port     int            `json:"port"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Database string         `json:"database"`
	FilePath string         `json:"file_path"`
}

// TestConnectionResponse represents the response for testing a connection
type TestConnectionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrDataSourceNotFound is returned when a data source is not found
var ErrDataSourceNotFound = errors.New("data source not found")

// ErrDataSourceExists is returned when a data source with the same name already exists
var ErrDataSourceExists = errors.New("data source with this name already exists")

// ErrInvalidDataSourceType is returned when an invalid data source type is provided
var ErrInvalidDataSourceType = errors.New("invalid data source type")
