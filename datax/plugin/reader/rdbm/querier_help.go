package rdbm

import (
	"context"
	"database/sql"
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
	return m.Instance() + "." + m.Name()
}

func (m *MockTable) AddField(bf *database.BaseField) {
	i, _ := strconv.Atoi(bf.FieldType().DatabaseTypeName())
	m.AppendField(NewMockField(bf, NewMockFieldType(database.GoType(i))))
}

type MockQuerier struct {
	QueryErr error
	FetchErr error
}

func (m *MockQuerier) Table(bt *database.BaseTable) database.Table {
	return NewMockTable(bt)
}

func (m *MockQuerier) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, m.QueryErr
}

func (m *MockQuerier) FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error) {
	return nil, m.FetchErr
}

func (m *MockQuerier) FetchRecord(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error) {
	_, err = handler.CreateRecord()
	if err != nil {
		return
	}
	return handler.OnRecord(element.NewDefaultRecord())
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

func TestJSONFromFile(filename string) *config.JSON {
	conf, err := config.NewJSONFromFile(filename)
	if err != nil {
		panic(err)
	}
	return conf
}

func TestJSONFromString(json string) *config.JSON {
	conf, err := config.NewJSONFromString(json)
	if err != nil {
		panic(err)
	}
	return conf
}
