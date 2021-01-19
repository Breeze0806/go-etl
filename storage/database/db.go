package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

type DB struct {
	*sql.DB
	Source
}

func Open(name string, conf *config.Json) (db *DB, err error) {
	db = &DB{}

	d, ok := dialects.dialect(name)
	if !ok {
		return nil, fmt.Errorf("dialect %v does not exsit", name)
	}

	db.Source, err = d.Source(NewBaseSource(conf))
	if err != nil {
		return nil, fmt.Errorf("dialect %v Source() err: %v", name, err)
	}

	db.DB, err = sql.Open(db.Source.DriverName(), db.Source.ConnectName())
	if err != nil {
		return nil, fmt.Errorf("Open(%v, %v) error: %v", db.Source.DriverName(), db.Source.ConnectName(), err)
	}

	var c Config
	err = json.Unmarshal([]byte(conf.String()), &c)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal(%v) error: %v", conf.String(), err)
	}

	db.SetMaxOpenConns(c.GetMaxOpenConns())
	db.SetMaxIdleConns(c.GetMaxIdleConns())
	if c.ConnMaxIdleTime.Duration != 0 {
		db.SetConnMaxIdleTime(c.ConnMaxIdleTime.Duration)
	}
	if c.ConnMaxLifetime.Duration != 0 {
		db.SetConnMaxLifetime(c.ConnMaxLifetime.Duration)
	}

	return
}

func (d *DB) FetchTable(ctx context.Context, t *BaseTable) (Table, error) {
	return d.FetchTableWithParam(ctx, NewTableQueryParam(d.Table(t)))
}

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
		adder.AddField(NewBaseField(names[i], types[i]))
	}
	return table, nil
}

func (d *DB) FetchRecord(ctx context.Context, param Parameter, onRecord func(element.Record) error) (err error) {
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
	return readRowsToRecord(rows, param, onRecord)
}

func (d *DB) FetchRecordWithTx(ctx context.Context, param Parameter, onRecord func(element.Record) error) (err error) {
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
	return readRowsToRecord(rows, param, onRecord)
}

func (d *DB) BatchExec(ctx context.Context, opts *ParameterOptions) (err error) {
	var param Parameter
	if param, err = execParam(opts); err != nil {
		return
	}
	return d.batchExec(ctx, param, opts.Records)
}

func (d *DB) BatchExecWithTx(ctx context.Context, opts *ParameterOptions) (err error) {
	var param Parameter
	if param, err = execParam(opts); err != nil {
		return
	}
	return d.batchExecWithTx(ctx, param, opts.Records)
}

func (d *DB) BatchExecStmtWithTx(ctx context.Context, opts *ParameterOptions) (err error) {
	var param Parameter
	if param, err = execParam(opts); err != nil {
		return
	}
	return d.batchExecStmtWithTx(ctx, param, opts.Records)
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
	if tx, err = d.BeginTx(ctx, param.TxOptions()); err != nil {
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
		if opts.Mode != "insert" {
			return nil, fmt.Errorf("table is not ExecParameter and mode is not insert")
		}
		param = NewInsertParam(opts.Table, opts.TxOptions)
	} else {
		if param, ok = execParams.ExecParam(opts.Mode, opts.TxOptions); !ok {
			if opts.Mode != "insert" {
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

func readRowsToRecord(rows *sql.Rows, param Parameter, onRecord func(element.Record) error) (err error) {
	var scanners []interface{}
	for _, f := range param.Table().Fields() {
		scanners = append(scanners, f.Scanner())
	}

	for rows.Next() {
		if err = rows.Scan(scanners...); err != nil {
			return fmt.Errorf("rows.Scan err: %v", err)
		}
		record := element.NewDefaultRecord()
		for _, v := range scanners {
			record.Add(v.(Scanner).Column())
		}
		if err = onRecord(record); err != nil {
			return fmt.Errorf("onRecord err: %v", err)
		}
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows.Err() err: %v", err)
	}
	return nil
}
