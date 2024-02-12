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

package dbms

import (
	"context"
	"database/sql"

	"github.com/Breeze0806/go-etl/storage/database"
)

type Execer interface {
	Table(*database.BaseTable) database.Table
	// Obtain relational database configuration through configuration
	PingContext(ctx context.Context) error

	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	// BaseDbHandler Basic Database Execution Handler Encapsulation
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error)

	BatchExec(ctx context.Context, opts *database.ParameterOptions) (err error)
	// NewBaseDbHandler Create a database execution handler encapsulation using the executor function newExecer and database transaction execution options opts
	BatchExecStmt(ctx context.Context, opts *database.ParameterOptions) (err error)

	BatchExecWithTx(ctx context.Context, opts *database.ParameterOptions) (err error)

	BatchExecStmtWithTx(ctx context.Context, opts *database.ParameterOptions) (err error)

	Close() error
}
