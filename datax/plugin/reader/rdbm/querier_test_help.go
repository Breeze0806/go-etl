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

//MockQuerier 模拟查询器
type MockQuerier struct {
	PingErr  error
	QueryErr error
	FetchErr error
}

//Table 新建表
func (m *MockQuerier) Table(bt *database.BaseTable) database.Table {
	return NewMockTable(bt)
}

//PingContext 测试关系型数据库的连接性
func (m *MockQuerier) PingContext(ctx context.Context) error {
	return m.PingErr
}

//QueryContext 查询
func (m *MockQuerier) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, m.QueryErr
}

//FetchTableWithParam 获取表参数
func (m *MockQuerier) FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error) {
	return nil, m.FetchErr
}

//FetchRecord 获取记录
func (m *MockQuerier) FetchRecord(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error) {
	_, err = handler.CreateRecord()
	if err != nil {
		return
	}
	return handler.OnRecord(element.NewDefaultRecord())
}

//FetchRecordWithTx 通过事务获取记录
func (m *MockQuerier) FetchRecordWithTx(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error) {
	_, err = handler.CreateRecord()
	if err != nil {
		return
	}
	return handler.OnRecord(element.NewDefaultRecord())
}

//Close 关闭资源
func (m *MockQuerier) Close() error {
	return nil
}

//TestJSON 从文件获取JSON配置
func TestJSON() *config.JSON {
	return TestJSONFromString(`{
		"name" : "rdbmreader",
		"developer":"Breeze0806",
		"dialect":"rdbm",
		"description":"rdbm is base package for relational database"
	}`)
}

//TestJSONFromString 从字符串获取JSON配置
func TestJSONFromString(json string) *config.JSON {
	conf, err := config.NewJSONFromString(json)
	if err != nil {
		panic(err)
	}
	return conf
}

//MockSender 模拟发送器
type MockSender struct {
	CreateErr error
	SendErr   error
}

//CreateRecord 创建记录
func (m *MockSender) CreateRecord() (element.Record, error) {
	return element.NewDefaultRecord(), m.CreateErr
}

//SendWriter 发往写入器
func (m *MockSender) SendWriter(record element.Record) error {
	return m.SendErr
}

//Flush 刷新至写入器
func (m *MockSender) Flush() error {
	return nil
}

//Terminate 终止发送数据
func (m *MockSender) Terminate() error {
	return nil
}

//Shutdown 关闭
func (m *MockSender) Shutdown() error {
	return nil
}
