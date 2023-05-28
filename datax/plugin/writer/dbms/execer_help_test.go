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
	"strconv"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

type MockFieldType struct {
	*database.BaseFieldType
	goType database.GoType
}

func NewMockFieldType(goType database.GoType) *MockFieldType {
	return &MockFieldType{
		BaseFieldType: database.NewBaseFieldType(&sql.ColumnType{}),
		goType:        goType,
	}
}

func (m *MockFieldType) DatabaseTypeName() string {
	return strconv.Itoa(int(m.goType))
}

func (m *MockFieldType) GoType() database.GoType {
	return m.goType
}

type MockField struct {
	*database.BaseField

	typ database.FieldType
}

func NewMockField(bf *database.BaseField, typ database.FieldType) *MockField {
	return &MockField{
		BaseField: bf,
		typ:       typ,
	}
}

func (m *MockField) Type() database.FieldType {
	return m.typ
}

func (m *MockField) Quoted() string {
	return m.Name()
}

func (m *MockField) BindVar(i int) string {
	return "$" + strconv.Itoa(i)
}

func (m *MockField) Select() string {
	return m.Name()
}

func (m *MockField) Scanner() database.Scanner {
	return nil
}

func (m *MockField) Valuer(c element.Column) database.Valuer {
	return database.NewGoValuer(m, c)
}

type MockTable struct {
	*database.BaseTable
	conf *config.JSON
}

func NewMockTable(bt *database.BaseTable) *MockTable {
	return &MockTable{
		BaseTable: bt,
	}
}

func (m *MockTable) Quoted() string {
	return m.Instance() + "." + m.Schema() + "." + m.Name()
}

func (m *MockTable) AddField(bf *database.BaseField) {
	i, _ := strconv.Atoi(bf.FieldType().DatabaseTypeName())
	m.AppendField(NewMockField(bf, NewMockFieldType(database.GoType(i))))
}

func (m *MockTable) SetConfig(conf *config.JSON) {
	m.conf = conf
}

type MockTableWithJudger struct {
	*MockTable

	retry    bool
	oneByOne bool
}

func NewMockTableWithJudger(bt *database.BaseTable,
	retry bool, oneByOne bool) *MockTableWithJudger {
	return &MockTableWithJudger{
		MockTable: NewMockTable(bt),
		retry:     retry,
		oneByOne:  oneByOne,
	}
}

func (m *MockTableWithJudger) ShouldRetry(err error) bool {
	return m.retry
}

func (m *MockTableWithJudger) ShouldOneByOne(err error) bool {
	return m.oneByOne
}

type MockExecer struct {
	PingErr  error
	QueryErr error
	FetchErr error
	BatchN   int
	BatchErr error
	ExecErr  error
}

func (m *MockExecer) Table(bt *database.BaseTable) database.Table {
	return NewMockTable(bt)
}

func (m *MockExecer) PingContext(ctx context.Context) error {
	return m.PingErr
}

func (m *MockExecer) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, m.QueryErr
}

func (m *MockExecer) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if query == "wait" {
		time.Sleep(100 * time.Millisecond)
	}
	return nil, m.ExecErr
}

func (m *MockExecer) FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error) {
	return NewMockTable(nil), m.FetchErr
}

func (m *MockExecer) BatchExec(ctx context.Context, opts *database.ParameterOptions) (err error) {
	m.BatchN--
	if m.BatchN <= 0 {
		return m.BatchErr
	}
	return nil
}

func (m *MockExecer) BatchExecStmt(ctx context.Context, opts *database.ParameterOptions) (err error) {
	return
}

func (m *MockExecer) BatchExecWithTx(ctx context.Context, opts *database.ParameterOptions) (err error) {
	return
}

func (m *MockExecer) BatchExecStmtWithTx(ctx context.Context, opts *database.ParameterOptions) (err error) {
	return
}

func (m *MockExecer) Close() error {
	return nil
}

func testJSON() *config.JSON {
	return testJSONFromString(`{
		"name" : "dbmswriter",
		"developer":"Breeze0806",
		"dialect":"dbms",
		"description":"dbms is base package for relational database"
	}`)
}

func testJSONFromString(json string) *config.JSON {
	conf, err := config.NewJSONFromString(json)
	if err != nil {
		panic(err)
	}
	return conf
}

type MockReceiver struct {
	err    error
	n      int
	ticker *time.Ticker
}

func NewMockReceiver(n int, err error, wait time.Duration) *MockReceiver {
	return &MockReceiver{
		err:    err,
		n:      n,
		ticker: time.NewTicker(wait),
	}
}

func NewMockReceiverWithoutWait(n int, err error) *MockReceiver {
	return &MockReceiver{
		err: err,
		n:   n,
	}
}

func (m *MockReceiver) GetFromReader() (element.Record, error) {
	m.n--
	if m.n <= 0 {
		return nil, m.err
	}
	if m.ticker != nil {
		<-m.ticker.C
		return element.NewDefaultRecord(), nil
	}
	return element.NewDefaultRecord(), nil
}

func (m *MockReceiver) Shutdown() error {
	m.ticker.Stop()
	return nil
}
