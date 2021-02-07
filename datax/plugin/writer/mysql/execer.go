package mysql

import (
	"context"
	"database/sql"

	"github.com/Breeze0806/go-etl/storage/database"
)

//Execer 执行器
type Execer interface {
	Table(*database.BaseTable) database.Table
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error)
	BatchExec(ctx context.Context, opts *database.ParameterOptions) (err error)
	BatchExecWithTx(ctx context.Context, opts *database.ParameterOptions) (err error)
	BatchExecStmtWithTx(ctx context.Context, opts *database.ParameterOptions) (err error)
	Close() error
}
