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

package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"reflect"
	"sync"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

var once sync.Once

type mockTableWithOther struct {
	*mockTable
	err        error
	execParams map[string]func(t Table, txOpts *sql.TxOptions) Parameter
}

func (m *mockTableWithOther) FetchFields(ctx context.Context, db *DB) error {
	if m.err != nil {
		return m.err
	}
	db.FetchTableWithParam(ctx, NewTableQueryParam(m.mockTable))
	return nil
}

func (m *mockTableWithOther) ExecParam(mode string, txOpts *sql.TxOptions) (p Parameter, ok bool) {
	var fn func(t Table, txOpts *sql.TxOptions) Parameter
	if fn, ok = m.execParams[mode]; ok {
		p = fn(m, txOpts)
		return
	}
	return
}

type mockTableWithNoAdder struct {
	*BaseTable
}

func (m *mockTableWithNoAdder) Quoted() string {
	return m.Instance() + "." + m.Schema() + "." + m.Name()
}

type mockParameter struct {
	*BaseParam
	queryErr error
	agrsErr  error
}

func (m *mockParameter) Query([]element.Record) (string, error) {
	if m.queryErr != nil {
		return "", m.queryErr
	}
	return "mock", nil
}

func (m *mockParameter) Agrs([]element.Record) ([]interface{}, error) {
	return nil, m.agrsErr
}

type mockDriver struct {
	rows *mockRows
}

func (m *mockDriver) Open(dsn string) (driver.Conn, error) {
	return nil, nil
}

func (m *mockDriver) OpenConnector(dsn string) (driver.Connector, error) {
	return &mockConnector{
		rows: m.rows,
	}, nil
}

type mockConnector struct {
	rows *mockRows
}

func (m *mockConnector) Connect(ctx context.Context) (driver.Conn, error) {
	return &mockConn{
		rows: m.rows,
	}, nil
}

func (m *mockConnector) Driver() driver.Driver {
	return &mockDriver{
		rows: m.rows,
	}
}

type mockConn struct {
	rows *mockRows
}

func (m *mockConn) Begin() (driver.Tx, error) {
	return &mockTx{}, nil
}

func (m *mockConn) Close() (err error) {
	return
}

func (m *mockConn) Prepare(query string) (driver.Stmt, error) {
	return &mockStmt{
		rows: m.rows,
	}, nil
}

func (m *mockConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return &mockStmt{
		rows: m.rows,
	}, nil
}

func (m *mockConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return &mockTx{}, nil
}

func (m *mockConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	m.rows.readCnt = 0
	return m.rows, nil
}

func (m *mockConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return nil, nil
}

func (m *mockConn) ResetSession(ctx context.Context) error {
	return nil
}

type mockStmt struct {
	rows *mockRows
}

func (m *mockStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	m.rows.readCnt = 0
	return m.rows, nil
}

func (m *mockStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return nil, nil
}

func (m *mockStmt) Close() error {
	return nil
}

func (m *mockStmt) NumInput() int {
	return -1
}

func (m *mockStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, nil
}

func (m *mockStmt) Query(args []driver.Value) (driver.Rows, error) {
	return m.rows, nil
}

type mockTx struct{}

func (m *mockTx) Commit() (err error) {
	return
}

func (m *mockTx) Rollback() (err error) {
	return
}

type mockRows struct {
	columns      []string
	types        []*mockFieldType
	columnValues [][]driver.Value
	readCnt      int
}

func (m *mockRows) Columns() []string {
	return m.columns
}

func (m *mockRows) ColumnTypeDatabaseTypeName(i int) string {
	return m.types[i].DatabaseTypeName()
}

func (m *mockRows) ColumnTypeLength(i int) (length int64, ok bool) {
	return m.types[i].Length()
}

func (m *mockRows) ColumnTypeNullable(i int) (nullable, ok bool) {
	return m.types[i].Nullable()
}

func (m *mockRows) ColumnTypePrecisionScale(i int) (int64, int64, bool) {
	return m.types[i].DecimalSize()
}

func (m *mockRows) ColumnTypeScanType(i int) reflect.Type {
	return m.types[i].ScanType()
}

func (m *mockRows) Close() (err error) {
	return
}

func (m *mockRows) Next(dest []driver.Value) error {
	if m.readCnt >= len(m.columnValues) {
		return io.EOF
	}
	for i, v := range m.columnValues[m.readCnt] {
		dest[i] = v
	}
	m.readCnt++
	return nil
}

func testDB(name string, conf *config.JSON) (db *DB, err error) {
	var source Source
	if source, err = NewSource(name, conf); err != nil {
		return
	}

	return NewDB(source)
}

func testMustDB(name string, conf *config.JSON) *DB {
	db, err := testDB(name, conf)
	if err != nil {
		panic(err)
	}
	return db
}

func registerMock() {
	UnregisterAllDialects()
	RegisterDialect("mock", &mockDialect{
		name: "mock",
	})
	RegisterDialect("mockErr", &mockDialect{
		name: "",
		err:  errors.New("mock error"),
	})
	RegisterDialect("test", &mockDialect{
		name: "test",
	})
	once.Do(func() {
		sql.Register("mock", &mockDriver{
			rows: &mockRows{
				columns: []string{
					"f1", "f2", "f3", "f4",
				},
				types: []*mockFieldType{
					newMockFieldType(GoTypeBool),
					newMockFieldType(GoTypeInt64),
					newMockFieldType(GoTypeFloat64),
					newMockFieldType(GoTypeString),
				},
				columnValues: [][]driver.Value{
					{false, int64(1), float64(1), string("1")},
					{true, int64(2), float64(2), string("2")},
				},
			},
		})
	})
}
