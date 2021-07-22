package rdbm

import (
	"context"
	"database/sql"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

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
