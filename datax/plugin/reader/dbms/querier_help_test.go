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
	"encoding/json"
	"math/big"
	"reflect"
	"strconv"

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

type MockQuerier struct {
	PingErr     error
	QueryErr    error
	FetchErr    error
	FetchMinErr error
	FetchMaxErr error
	isTime      bool
	config      *config.JSON
}

func (m *MockQuerier) Table(bt *database.BaseTable) database.Table {
	return NewMockTable(bt)
}

func (m *MockQuerier) PingContext(ctx context.Context) error {
	return m.PingErr
}

func (m *MockQuerier) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, m.QueryErr
}

func (m *MockQuerier) FetchTableWithParam(ctx context.Context,
	param database.Parameter) (database.Table, error) {
	if _, ok := param.(*SplitParam); ok {
		t := NewMockTable(database.NewBaseTable("db", "schema", "name"))
		typ := database.GoTypeInt64
		if m.isTime {
			typ = database.GoTypeTime
		}
		t.AddField(database.NewBaseField(0, "f1", NewMockFieldType(typ)))
		return t, m.FetchErr
	}
	return nil, m.FetchErr
}

func (m *MockQuerier) FetchRecord(ctx context.Context,
	param database.Parameter, handler database.FetchHandler) (err error) {

	r, err := handler.CreateRecord()
	if err != nil {
		return
	}
	switch param.(type) {
	case *MinParam:
		if m.FetchMinErr != nil {
			return m.FetchMinErr
		}
		r.Add(element.NewDefaultColumn(element.NewBigIntColumnValue(big.NewInt(10000)), "f1", 0))
	case *MaxParam:
		if m.FetchMaxErr != nil {
			return m.FetchMaxErr
		}
		r.Add(element.NewDefaultColumn(element.NewBigIntColumnValue(big.NewInt(30000)), "f1", 0))
	}
	return handler.OnRecord(r)
}

func (m *MockQuerier) FetchRecordWithTx(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error) {
	_, err = handler.CreateRecord()
	if err != nil {
		return
	}
	return handler.OnRecord(element.NewDefaultRecord())
}

func (m *MockQuerier) Close() error {
	return nil
}

func testJSON() *config.JSON {
	return testJSONFromString(`{
		"name" : "dbmsreader",
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

type MockSender struct {
	CreateErr error
	SendErr   error
}

func (m *MockSender) CreateRecord() (element.Record, error) {
	return element.NewDefaultRecord(), m.CreateErr
}

func (m *MockSender) SendWriter(record element.Record) error {
	return m.SendErr
}

func (m *MockSender) Flush() error {
	return nil
}

func (m *MockSender) Terminate() error {
	return nil
}

func (m *MockSender) Shutdown() error {
	return nil
}

func equalConfigJSON(gotConfig, wantConfig *config.JSON) bool {
	var got, want interface{}
	err := json.Unmarshal([]byte(gotConfig.String()), &got)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(wantConfig.String()), &want)
	if err != nil {
		panic(err)
	}
	return reflect.DeepEqual(got, want)
}
