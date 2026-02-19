package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) Migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		email TEXT,
		role TEXT NOT NULL DEFAULT 'user',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

	CREATE TABLE IF NOT EXISTS data_sources (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		host TEXT,
		port INTEGER,
		username TEXT,
		password TEXT,
		database TEXT,
		file_path TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		UNIQUE(user_id, name)
	);
	CREATE INDEX IF NOT EXISTS idx_data_sources_user_id ON data_sources(user_id);

	CREATE TABLE IF NOT EXISTS sync_tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		reader_config TEXT,
		writer_config TEXT,
		status TEXT NOT NULL DEFAULT 'draft',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		UNIQUE(user_id, name)
	);
	CREATE INDEX IF NOT EXISTS idx_sync_tasks_user_id ON sync_tasks(user_id);

	CREATE TABLE IF NOT EXISTS sync_jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		progress INTEGER NOT NULL DEFAULT 0,
		error_message TEXT,
		started_at DATETIME,
		finished_at DATETIME,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_sync_jobs_task_id ON sync_jobs(task_id);
	CREATE INDEX IF NOT EXISTS idx_sync_jobs_user_id ON sync_jobs(user_id);
	`

	_, err := db.Exec(schema)
	return err
}
