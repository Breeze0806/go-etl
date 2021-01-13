package database

import (
	"database/sql"
	"strconv"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

type mockNilDialect struct {
}

func (m *mockNilDialect) Source(*BaseSource) (Source, error) {
	return nil, nil
}

type mockFieldType struct {
	*BaseFieldType
	goType GoType
}

func newMockFieldType(goType GoType) *mockFieldType {
	return &mockFieldType{
		BaseFieldType: NewBaseFieldType(&sql.ColumnType{}),
		goType:        goType,
	}
}

func (m *mockFieldType) GoType() GoType {
	return m.goType
}

type mockField struct {
	*BaseField

	typ FieldType
}

func newMockField(bf *BaseField, typ FieldType) *mockField {
	return &mockField{
		BaseField: bf,
		typ:       typ,
	}
}

func (m *mockField) Type() FieldType {
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

func (m *mockField) Scanner() Scanner {
	return nil
}

func (m *mockField) Valuer(c element.Column) Valuer {
	return NewGoValuer(m, c)
}

type mockTable struct {
	*BaseTable
}

func newMockTable(bt *BaseTable) *mockTable {
	return &mockTable{
		BaseTable: bt,
	}
}

func (m *mockTable) Quoted() string {
	return m.Instance() + "." + m.Schema() + "." + m.Name()
}

func (m *mockTable) String() string {
	return m.Instance() + "." + m.Schema() + "." + m.Name()
}

func testJsonFromString(s string) *config.Json {
	json, err := config.NewJsonFromString(s)
	if err != nil {
		panic(err)
	}
	return json
}
