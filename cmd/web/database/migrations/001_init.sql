-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email TEXT,
    role TEXT NOT NULL DEFAULT 'user' CHECK(role IN ('admin', 'user', 'readonly')),
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Data sources table
CREATE TABLE IF NOT EXISTS data_sources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('mysql', 'postgresql', 'oracle', 'db2', 'sqlite', 'sqlserver', 'dameng', 'csv', 'xlsx')),
    host TEXT,
    port INTEGER,
    username TEXT,
    password TEXT,
    database TEXT,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_data_sources_user_id ON data_sources(user_id);
CREATE INDEX IF NOT EXISTS idx_data_sources_type ON data_sources(type);

-- Sync tasks table
CREATE TABLE IF NOT EXISTS sync_tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    reader_config TEXT NOT NULL DEFAULT '{}',
    writer_config TEXT NOT NULL DEFAULT '{}',
    status TEXT NOT NULL DEFAULT 'draft' CHECK(status IN ('draft', 'active', 'paused', 'deleted')),
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sync_tasks_user_id ON sync_tasks(user_id);
CREATE INDEX IF NOT EXISTS idx_sync_tasks_status ON sync_tasks(status);

-- Sync jobs table
CREATE TABLE IF NOT EXISTS sync_jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'running', 'completed', 'failed', 'cancelled')),
    started_at INTEGER,
    finished_at INTEGER,
    error_message TEXT,
    progress REAL DEFAULT 0.0 CHECK(progress >= 0 AND progress <= 100),
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (task_id) REFERENCES sync_tasks(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sync_jobs_task_id ON sync_jobs(task_id);
CREATE INDEX IF NOT EXISTS idx_sync_jobs_status ON sync_jobs(status);

CREATE TRIGGER IF NOT EXISTS update_users_updated_at
AFTER UPDATE ON users
BEGIN
    UPDATE users SET updated_at = strftime('%s', 'now') WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_data_sources_updated_at
AFTER UPDATE ON data_sources
BEGIN
    UPDATE data_sources SET updated_at = strftime('%s', 'now') WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_sync_tasks_updated_at
AFTER UPDATE ON sync_tasks
BEGIN
    UPDATE sync_tasks SET updated_at = strftime('%s', 'now') WHERE id = NEW.id;
END;
