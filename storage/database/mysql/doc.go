// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package mysql implements the Dialect for MySQL databases, supporting MySQL 5.6+ corresponding databases.
// The driver used is github.com/go-sql-driver/mysql.
// The data source Source uses BaseSource to simplify its implementation, wrapping the github.com/go-sql-driver/mysql driver. For database configuration, it needs to be consistent with Config.
// The Table implementation uses BaseTable to simplify its implementation, also based on github.com/go-sql-driver/mysql.
// Table implements the FieldAdder approach to acquire columns. In ExecParameter, it implements the replace mode for the replace into bulk data processing mode, and reuses the existing database.InsertParam for the insert mode.
// The Field uses BaseField to simplify its implementation, where FieldType adopts the original sql.ColumnType and implements ValuerGoType.
// The Scanner uses BaseScanner to simplify its implementation.
// The Valuer uses the implementation approach of GoValuer.

package mysql
