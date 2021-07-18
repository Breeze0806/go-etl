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
	return m.Instance() + "." + m.Name()
}

//AddField 新增列
func (m *MockTable) AddField(bf *database.BaseField) {
	i, _ := strconv.Atoi(bf.FieldType().DatabaseTypeName())
	m.AppendField(NewMockField(bf, NewMockFieldType(database.GoType(i))))
}

//MockQuerier 模拟查询器
type MockQuerier struct {
	QueryErr error
	FetchErr error
}

//Table 新建表
func (m *MockQuerier) Table(bt *database.BaseTable) database.Table {
	return NewMockTable(bt)
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

//TestJSONFromFile 从文件获取JSON配置
func TestJSONFromFile(filename string) *config.JSON {
	conf, err := config.NewJSONFromFile(filename)
	if err != nil {
		panic(err)
	}
	return conf
}

//TestJSONFromString 从字符串获取JSON配置
func TestJSONFromString(json string) *config.JSON {
	conf, err := config.NewJSONFromString(json)
	if err != nil {
		panic(err)
	}
	return conf
}
