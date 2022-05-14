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

package rdbm

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

//MockFieldType 模拟字段类型测试类
type MockFieldType struct {
	*database.BaseFieldType
	goType database.GoType
}

//NewMockFieldType 新建模拟字段类型测试类
func NewMockFieldType(goType database.GoType) *MockFieldType {
	return &MockFieldType{
		BaseFieldType: database.NewBaseFieldType(&sql.ColumnType{}),
		goType:        goType,
	}
}

//DatabaseTypeName 字段类型名称，如DECIMAL,VARCHAR, BIGINT等数据库类型
func (m *MockFieldType) DatabaseTypeName() string {
	return strconv.Itoa(int(m.goType))
}

//GoType 字段类型对应的golang类型
func (m *MockFieldType) GoType() database.GoType {
	return m.goType
}

//MockField 模拟列字段测试类
type MockField struct {
	*database.BaseField

	typ database.FieldType
}

//NewMockField 新建模拟字段测试类
func NewMockField(bf *database.BaseField, typ database.FieldType) *MockField {
	return &MockField{
		BaseField: bf,
		typ:       typ,
	}
}

//Type 字段类型
func (m *MockField) Type() database.FieldType {
	return m.typ
}

//Quoted 引用
func (m *MockField) Quoted() string {
	return m.Name()
}

//BindVar 占位符
func (m *MockField) BindVar(i int) string {
	return "$" + strconv.Itoa(i)
}

//Select 查询时使用的字段
func (m *MockField) Select() string {
	return m.Name()
}

//Scanner 空值
func (m *MockField) Scanner() database.Scanner {
	return nil
}

//Valuer 类型赋值器
func (m *MockField) Valuer(c element.Column) database.Valuer {
	return database.NewGoValuer(m, c)
}

//MockTable 模拟表测试类
type MockTable struct {
	*database.BaseTable
}

//NewMockTable 新建模拟表测试类
func NewMockTable(bt *database.BaseTable) *MockTable {
	return &MockTable{
		BaseTable: bt,
	}
}

//Quoted 引用
func (m *MockTable) Quoted() string {
	return m.Instance() + "." + m.Schema() + "." + m.Name()
}

//AddField 新增列
func (m *MockTable) AddField(bf *database.BaseField) {
	i, _ := strconv.Atoi(bf.FieldType().DatabaseTypeName())
	m.AppendField(NewMockField(bf, NewMockFieldType(database.GoType(i))))
}

//MockExecer 模拟执行器
type MockExecer struct {
	PingErr  error
	QueryErr error
	FetchErr error
	BatchN   int
	BatchErr error
}

//Table 新建表
func (m *MockExecer) Table(bt *database.BaseTable) database.Table {
	return NewMockTable(bt)
}

//PingContext 测试关系型数据库连接情况
func (m *MockExecer) PingContext(ctx context.Context) error {
	return m.PingErr
}

//QueryContext 查询
func (m *MockExecer) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, m.QueryErr
}

//ExecContext 获取表参数
func (m *MockExecer) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

//FetchTableWithParam 获取表参数
func (m *MockExecer) FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error) {
	return nil, m.FetchErr
}

//BatchExec 批量执行
func (m *MockExecer) BatchExec(ctx context.Context, opts *database.ParameterOptions) (err error) {
	m.BatchN--
	if m.BatchN <= 0 {
		return m.BatchErr
	}
	return nil
}

//BatchExecWithTx 批量事务执行
func (m *MockExecer) BatchExecWithTx(ctx context.Context, opts *database.ParameterOptions) (err error) {
	return
}

//BatchExecStmtWithTx 批量事务执行
func (m *MockExecer) BatchExecStmtWithTx(ctx context.Context, opts *database.ParameterOptions) (err error) {
	return
}

//Close 关闭
func (m *MockExecer) Close() error {
	return nil
}

//testJSON 从文件获取JSON配置
func testJSON() *config.JSON {
	return testJSONFromString(`{
		"name" : "rdbmwriter",
		"developer":"Breeze0806",
		"dialect":"rdbm",
		"description":"rdbm is base package for relational database"
	}`)
}

//testJSONFromString 从字符串获取JSON配置
func testJSONFromString(json string) *config.JSON {
	conf, err := config.NewJSONFromString(json)
	if err != nil {
		panic(err)
	}
	return conf
}

//MockReceiver 模拟接受器
type MockReceiver struct {
	err    error
	n      int
	ticker *time.Ticker
}

//NewMockReceiver 新建等待模拟接受器
func NewMockReceiver(n int, err error, wait time.Duration) *MockReceiver {
	return &MockReceiver{
		err:    err,
		n:      n,
		ticker: time.NewTicker(wait),
	}
}

//NewMockReceiverWithoutWait 新建无等待模拟接受器
func NewMockReceiverWithoutWait(n int, err error) *MockReceiver {
	return &MockReceiver{
		err: err,
		n:   n,
	}
}

//GetFromReader 从读取器获取记录
func (m *MockReceiver) GetFromReader() (element.Record, error) {
	m.n--
	if m.n <= 0 {
		return nil, m.err
	}
	if m.ticker != nil {
		select {
		case <-m.ticker.C:
			return element.NewDefaultRecord(), nil
		}
	}
	return element.NewDefaultRecord(), nil
}

//Shutdown 关闭
func (m *MockReceiver) Shutdown() error {
	m.ticker.Stop()
	return nil
}
