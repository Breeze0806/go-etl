package rdbm

import (
	"context"
	"database/sql"

	"github.com/Breeze0806/go-etl/storage/database"
)

//Querier 询问器
type Querier interface {
	Table(*database.BaseTable) database.Table
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error)
	FetchRecord(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error)
	FetchRecordWithTx(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error)
	Close() error
}
