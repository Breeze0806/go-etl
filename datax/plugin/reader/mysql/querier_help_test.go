package mysql

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

type mockFieldType struct {
	*database.BaseFieldType
	goType database.GoType
}

func newMockFieldType(goType database.GoType) *mockFieldType {
	return &mockFieldType{
		BaseFieldType: database.NewBaseFieldType(&sql.ColumnType{}),
		goType:        goType,
	}
}

func (m *mockFieldType) DatabaseTypeName() string {
	return strconv.Itoa(int(m.goType))
}

func (m *mockFieldType) GoType() database.GoType {
	return m.goType
}

type mockField struct {
	*database.BaseField

	typ database.FieldType
}

func newMockField(bf *database.BaseField, typ database.FieldType) *mockField {
	return &mockField{
		BaseField: bf,
		typ:       typ,
	}
}

func (m *mockField) Type() database.FieldType {
	return m.typ
}

func (m *mockField) Quoted() string {
	return m.Name()
}

func (m *mockField) BindVar(i int) string {
	return "$" + strconv.Itoa(i)
}

func (m *mockField) Select() string {
	return m.Name()
}

func (m *mockField) Scanner() database.Scanner {
	return nil
}

func (m *mockField) Valuer(c element.Column) database.Valuer {
	return database.NewGoValuer(m, c)
}

type mockTable struct {
	*database.BaseTable
}

func newMockTable(bt *database.BaseTable) *mockTable {
	return &mockTable{
		BaseTable: bt,
	}
}

func (m *mockTable) Quoted() string {
	return m.Instance() + "." + m.Name()
}

func (m *mockTable) AddField(bf *database.BaseField) {
	i, _ := strconv.Atoi(bf.FieldType().DatabaseTypeName())
	m.AppendField(newMockField(bf, newMockFieldType(database.GoType(i))))
}

type mockQuerier struct {
	queryErr error
	fetchErr error
}

func (m *mockQuerier) Table(bt *database.BaseTable) database.Table {
	return newMockTable(bt)
}

func (m *mockQuerier) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, m.queryErr
}

func (m *mockQuerier) FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error) {
	return nil, m.fetchErr
}

func (m *mockQuerier) FetchRecord(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error) {
	_, err = handler.CreateRecord()
	if err != nil {
		return
	}
	return handler.OnRecord(element.NewDefaultRecord())
}

func (m *mockQuerier) FetchRecordWithTx(ctx context.Context, param database.Parameter, handler database.FetchHandler) (err error) {
	_, err = handler.CreateRecord()
	if err != nil {
		return
	}
	return handler.OnRecord(element.NewDefaultRecord())
}

func (m *mockQuerier) Close() error {
	return nil
}

func testJSONFromFile(filename string) *config.JSON {
	conf, err := config.NewJSONFromFile(filename)
	if err != nil {
		panic(err)
	}
	return conf
}

func testJSONFromString(json string) *config.JSON {
	conf, err := config.NewJSONFromString(json)
	if err != nil {
		panic(err)
	}
	return conf
}
