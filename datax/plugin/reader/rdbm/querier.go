package rdbm

import (
	"context"
	"database/sql"

	"github.com/Breeze0806/go-etl/storage/database"
)

//Querier 询问器
type Querier interface {
	//通过基础表信息获取具体表
	Table(*database.BaseTable) database.Table
	//通过query查询语句进行查询
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	//通过参数param获取具体表
	FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error)
	//通过参数param，处理句柄handler获取记录
	FetchRecord(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error)
	//通过参数param，处理句柄handler使用事务获取记录
	FetchRecordWithTx(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error)
	//关闭资源
	Close() error
}
