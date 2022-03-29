package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Breeze0806/go-etl/element"
)

//写入数据库模式
const (
	WriteModeInsert = "insert"
)

//FetchHandler 获取记录句柄
type FetchHandler interface {
	OnRecord(element.Record) error
	CreateRecord() (element.Record, error)
}

//BaseFetchHandler 基础获取记录句柄
type BaseFetchHandler struct {
	onRecord     func(element.Record) error
	createRecord func() (element.Record, error)
}

//NewBaseFetchHandler 创建基础获取记录句柄
func NewBaseFetchHandler(createRecord func() (element.Record, error),
	onRecord func(element.Record) error) *BaseFetchHandler {
	return &BaseFetchHandler{
		onRecord:     onRecord,
		createRecord: createRecord,
	}
}

//OnRecord 处理记录r
func (b *BaseFetchHandler) OnRecord(r element.Record) error {
	return b.onRecord(r)
}

//CreateRecord 创建记录
func (b *BaseFetchHandler) CreateRecord() (element.Record, error) {
	return b.createRecord()
}

//DB 用户维护数据库连接池
type DB struct {
	Source

	db *sql.DB
}

//NewDB 从数据源source中获取数据库连接池
func NewDB(source Source) (d *DB, err error) {
	d = &DB{
		Source: source,
	}

	d.db, err = sql.Open(d.Source.DriverName(), d.Source.ConnectName())
	if err != nil {
		return nil, fmt.Errorf("Open(%v, %v) error: %v", d.Source.DriverName(), d.Source.ConnectName(), err)
	}

	var c *Config
	c, err = NewConfig(d.Config())
	if err != nil {
		return nil, err
	}

	d.db.SetMaxOpenConns(c.Pool.GetMaxOpenConns())
	d.db.SetMaxIdleConns(c.Pool.GetMaxIdleConns())
	if c.Pool.ConnMaxIdleTime.Duration != 0 {
		d.db.SetConnMaxIdleTime(c.Pool.ConnMaxIdleTime.Duration)
	}
	if c.Pool.ConnMaxLifetime.Duration != 0 {
		d.db.SetConnMaxLifetime(c.Pool.ConnMaxLifetime.Duration)
	}

	return
}

//FetchTable 通过上下文ctx和基础表数据t，获取对应的表并会返回错误
func (d *DB) FetchTable(ctx context.Context, t *BaseTable) (Table, error) {
	return d.FetchTableWithParam(ctx, NewTableQueryParam(d.Table(t)))
}

//FetchTableWithParam 通过上下文ctx和sql参数param，获取对应的表并会返回错误
func (d *DB) FetchTableWithParam(ctx context.Context, param Parameter) (Table, error) {
	table := param.Table()
	if fetcher, ok := table.(FieldsFetcher); ok {
		if err := fetcher.FetchFields(ctx, d); err != nil {
			return nil, err
		}
		return table, nil
	}
	adder, ok := table.(FieldAdder)
	if !ok {
		return nil, fmt.Errorf("Table is not FieldAdder")
	}
	query, agrs, err := getQueryAndAgrs(param, nil)
	if err != nil {
		return nil, err
	}
	rows, err := d.QueryContext(ctx, query, agrs...)
	if err != nil {
		return nil, fmt.Errorf("QueryContext(%v) err: %v", query, err)
	}
	defer rows.Close()
	names, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("rows.Columns() err: %v", err)
	}
	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("rows.ColumnTypes() err: %v", err)
	}
	for i := range names {
		adder.AddField(NewBaseField(i, names[i], NewBaseFieldType(types[i])))
	}

	for _, v := range table.Fields() {
		if !v.Type().IsSupportted() {
			return nil, fmt.Errorf("table: %v filed:%v type(%v) is not supportted", table.Quoted(), v.Name(), v.Type().DatabaseTypeName())
		}
	}
	return table, nil
}

//FetchRecord 通过上下文ctx，sql参数param以及记录处理函数onRecord
//获取多行记录返回错误
func (d *DB) FetchRecord(ctx context.Context, param Parameter, handler FetchHandler) (err error) {
	var query string
	var agrs []interface{}

	if query, agrs, err = getQueryAndAgrs(param, nil); err != nil {
		return
	}

	var rows *sql.Rows
	if rows, err = d.QueryContext(ctx, query, agrs...); err != nil {
		return fmt.Errorf("QueryContext(%v) err: %v", query, err)
	}
	defer rows.Close()
	return readRowsToRecord(rows, param, handler)
}

//FetchRecordWithTx 通过上下文ctx，sql参数param以及记录处理函数onRecord
//使用事务获取多行记录并返回错误
func (d *DB) FetchRecordWithTx(ctx context.Context, param Parameter, handler FetchHandler) (err error) {
	var query string
	var agrs []interface{}

	if query, agrs, err = getQueryAndAgrs(param, nil); err != nil {
		return
	}

	var tx *sql.Tx

	if tx, err = d.BeginTx(ctx, param.TxOptions()); err != nil {
		return fmt.Errorf("BeginTx(%+v) err: %v", param.TxOptions(), err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var rows *sql.Rows
	if rows, err = tx.QueryContext(ctx, query, agrs...); err != nil {
		return fmt.Errorf("QueryContext(%v) error: %v", query, err)
	}
	defer rows.Close()
	return readRowsToRecord(rows, param, handler)
}

//BatchExec 批量执行sql并处理多行记录
func (d *DB) BatchExec(ctx context.Context, opts *ParameterOptions) (err error) {
	var param Parameter
	if param, err = execParam(opts); err != nil {
		return
	}
	return d.batchExec(ctx, param, opts.Records)
}

//BatchExecWithTx 批量事务执行sql并处理多行记录
func (d *DB) BatchExecWithTx(ctx context.Context, opts *ParameterOptions) (err error) {
	var param Parameter
	if param, err = execParam(opts); err != nil {
		return
	}
	return d.batchExecWithTx(ctx, param, opts.Records)
}

//BatchExecStmtWithTx 批量事务prepare执行sql并处理多行记录
func (d *DB) BatchExecStmtWithTx(ctx context.Context, opts *ParameterOptions) (err error) {
	var param Parameter
	if param, err = execParam(opts); err != nil {
		return
	}
	return d.batchExecStmtWithTx(ctx, param, opts.Records)
}

//BeginTx 获取事务
func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return d.db.BeginTx(ctx, opts)
}

//PingContext 通过query查询多行数据
func (d *DB) PingContext(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

//QueryContext 通过query查询多行数据
func (d *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

//ExecContext 执行query并获取结果
func (d *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

//Close 关闭数据连接池
func (d *DB) Close() (err error) {
	if d.db != nil {
		return d.db.Close()
	}
	return
}

func (d *DB) batchExec(ctx context.Context, param Parameter, records []element.Record) (err error) {
	var query string
	var agrs []interface{}

	if query, agrs, err = getQueryAndAgrs(param, records); err != nil {
		return
	}

	if _, err = d.ExecContext(ctx, query, agrs...); err != nil {
		return fmt.Errorf("ExecContext(%v) err: %v", query, err)
	}
	return nil
}

func (d *DB) batchExecWithTx(ctx context.Context, param Parameter, records []element.Record) (err error) {
	var query string
	var agrs []interface{}

	if query, agrs, err = getQueryAndAgrs(param, records); err != nil {
		return
	}

	var tx *sql.Tx
	if tx, err = d.db.BeginTx(ctx, param.TxOptions()); err != nil {
		return fmt.Errorf("BeginTx(%+v) err: %v", param.TxOptions(), err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if _, err = tx.ExecContext(ctx, query, agrs...); err != nil {
		return fmt.Errorf("ExecContext(%v) err: %v", query, err)
	}
	return nil
}

func (d *DB) batchExecStmtWithTx(ctx context.Context, param Parameter, records []element.Record) (err error) {
	var query string
	if query, err = param.Query(records); err != nil {
		return fmt.Errorf("param.Query() err: %v", err)
	}

	var tx *sql.Tx
	if tx, err = d.db.BeginTx(ctx, param.TxOptions()); err != nil {
		return fmt.Errorf("BeginTx() err: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var stmt *sql.Stmt
	if stmt, err = tx.PrepareContext(ctx, query); err != nil {
		return fmt.Errorf("tx.PrepareContext(%v) err: %v", query, err)
	}
	defer func() {
		stmt.Close()
	}()

	for _, r := range records {
		var valuers []interface{}
		if valuers, err = param.Agrs([]element.Record{
			r,
		}); err != nil {
			return fmt.Errorf("param.Args() err: %v", err)
		}
		if _, err = stmt.ExecContext(ctx, valuers...); err != nil {
			return fmt.Errorf("stmt.ExecContext err: %v", err)
		}
	}
	if _, err = stmt.ExecContext(ctx); err != nil {
		return fmt.Errorf("stmt.ExecContext err: %v", err)
	}
	return
}

func execParam(opts *ParameterOptions) (param Parameter, err error) {
	execParams, ok := opts.Table.(ExecParameter)
	if !ok {
		if opts.Mode != WriteModeInsert {
			return nil, fmt.Errorf("table is not ExecParameter and mode is not insert")
		}
		param = NewInsertParam(opts.Table, opts.TxOptions)
	} else {
		if param, ok = execParams.ExecParam(opts.Mode, opts.TxOptions); !ok {
			if opts.Mode != WriteModeInsert {
				return nil, fmt.Errorf("ExecParam is not exist and mode is not insert")
			}
			param = NewInsertParam(opts.Table, opts.TxOptions)
		}
	}
	return
}

func getQueryAndAgrs(param Parameter, records []element.Record) (query string, agrs []interface{}, err error) {
	if query, err = param.Query(records); err != nil {
		err = fmt.Errorf("param.Query() err: %v", err)
		return
	}
	if agrs, err = param.Agrs(records); err != nil {
		query = ""
		err = fmt.Errorf("param.Agrs() err: %v", err)
		return
	}
	return
}

func readRowsToRecord(rows *sql.Rows, param Parameter, handler FetchHandler) (err error) {
	var scanners []interface{}
	for _, f := range param.Table().Fields() {
		scanners = append(scanners, f.Scanner())
	}

	for rows.Next() {
		if err = rows.Scan(scanners...); err != nil {
			return fmt.Errorf("rows.Scan err: %v", err)
		}
		var record element.Record
		if record, err = handler.CreateRecord(); err != nil {
			return fmt.Errorf("CreateRecord err: %v", err)
		}
		for _, v := range scanners {
			record.Add(v.(Scanner).Column())
		}
		if err = handler.OnRecord(record); err != nil {
			return fmt.Errorf("OnRecord err: %v", err)
		}
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows.Err() err: %v", err)
	}
	return nil
}
