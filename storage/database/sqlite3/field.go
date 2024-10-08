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

package sqlite3

import (
	"fmt"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

var (
	dateLayout      = element.DefaultTimeFormat[:10]
	timestampLayout = element.DefaultTimeFormat[:26]
)

// Field - Represents a field in a database table.
type Field struct {
	*database.BaseField
	database.BaseConfigSetter
}

// NewField - Generates a field based on basic column attributes.
func NewField(bf *database.BaseField) *Field {
	return &Field{
		BaseField: bf,
	}
}

// Quoted - Used for quoting in SQL statements.
func (f *Field) Quoted() string {
	return Quoted(f.Name())
}

// BindVar - SQL placeholder used in SQL statements.
func (f *Field) BindVar(i int) string {
	return "?"
}

// Select - Represents a field for querying purposes in SQL query statements.
func (f *Field) Select() string {
	return f.Quoted()
}

// Type - Represents the type of the field.
func (f *Field) Type() database.FieldType {
	return NewFieldType(f.FieldType())
}

// Scanner - Used for reading data from a field.
func (f *Field) Scanner() database.Scanner {
	return NewScanner(f)
}

// Valuer - Handles data processing using GoValuer.
func (f *Field) Valuer(c element.Column) database.Valuer {
	return database.NewGoValuer(f, c)
}

// FieldType - Represents the type of a field.
type FieldType struct {
	*database.BaseFieldType

	goType database.GoType
}

// NewFieldType - Creates a new field type.
func NewFieldType(typ database.ColumnType) *FieldType {
	f := &FieldType{
		BaseFieldType: database.NewBaseFieldType(typ),
	}
	switch f.DatabaseTypeName() {
	case "INTEGER", "BLOB", "NUMERIC", "REAL", "TEXT":
		f.goType = database.GoTypeString
	}
	return f
}

// IsSupported - Indicates whether parsing is supported for a specific type.
func (f *FieldType) IsSupported() bool {
	return f.GoType() != database.GoTypeUnknown
}

// GoType - Returns the Golang type used when processing numerical values.
func (f *FieldType) GoType() database.GoType {
	return f.goType
}

// Scanner - A scanner used for reading data based on the column type.
type Scanner struct {
	database.BaseScanner

	f *Field
}

// NewScanner - Generates a scanner based on the column type.
func NewScanner(f *Field) *Scanner {
	return &Scanner{
		f: f,
	}
}

// Scan - Reads data from a column based on its type.
func (s *Scanner) Scan(src any) (err error) {
	defer s.f.SetError(&err)
	var cv element.ColumnValue
	byteSize := element.ByteSize(src)
	switch s.f.Type().DatabaseTypeName() {
	case "INTEGER":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBigIntColumnValue()
		case int64:
			cv = element.NewBigIntColumnValueFromInt64(data)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeBigInt)
		}
	case "BLOB":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBytesColumnValue()
		case []byte:
			cv = element.NewBytesColumnValue(data)
		default:
			return fmt.Errorf("src is %v(%T),but not %v", src, src, element.TypeBytes)
		}
	case "NUMERIC", "REAL":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilTimeColumnValue()
		case int64:
			cv = element.NewBigIntColumnValueFromInt64(data)
		case float64:
			cv = element.NewDecimalColumnValueFromFloat(data)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeDecimal)
		}
	case "TEXT":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilStringColumnValue()
		case string:
			cv = element.NewStringColumnValue(data)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeString)
		}
	default:
		return fmt.Errorf("src is %v(%T), but db type is %v", src, src, s.f.Type().DatabaseTypeName())
	}
	s.SetColumn(element.NewDefaultColumn(cv, s.f.Name(), byteSize))
	return
}
