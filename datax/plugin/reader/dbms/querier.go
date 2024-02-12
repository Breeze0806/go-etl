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

// Querier - Query Executor
type Querier interface {
	// Obtain a specific table based on basic table information.
	Table(*database.BaseTable) database.Table
	// Check connectivity.
	PingContext(ctx context.Context) error
	// Perform a query using the specified query statement.
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	// Obtain a specific table based on the provided parameters.
	FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error)
	// Retrieve records using the provided parameters and the handler.
	FetchRecord(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error)
	// Retrieve records using the provided parameters, the handler, and within a transaction.
	FetchRecordWithTx(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error)
	// Close resources.
	Close() error
}
