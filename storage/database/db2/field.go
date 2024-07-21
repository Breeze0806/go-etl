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

package db2

import (
	"fmt"
	"time"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

var (
	dateLayout      = element.DefaultTimeFormat[:10]
	timestampLayout = element.DefaultTimeFormat[:26]
	timeLayout      = timestampLayout
)

// Field represents a database field.
type Field struct {
	*database.BaseField
	database.BaseConfigSetter
}

// NewField creates a new field based on basic column attributes.
func NewField(bf *database.BaseField) *Field {
	return &Field{
		BaseField: bf,
	}
}

// Quoted is used for quoting in SQL statements.
func (f *Field) Quoted() string {
	return Quoted(f.Name())
}

// BindVar represents an SQL placeholder used in SQL statements.
func (f *Field) BindVar(_ int) string {
	return "?"
}

// Select is the field used during queries for SQL SELECT statements.
func (f *Field) Select() string {
	return Quoted(f.Name())
}

// Type represents the type of the field.
func (f *Field) Type() database.FieldType {
	return NewFieldType(f.FieldType())
}

// Scanner is used for reading data.
func (f *Field) Scanner() database.Scanner {
	return NewScanner(f)
}

// Valuer adopts GoValuer for processing data.
func (f *Field) Valuer(c element.Column) database.Valuer {
	return database.NewGoValuer(f, c)
}

// FieldType represents the type of a field.
type FieldType struct {
	*database.BaseFieldType

	goType database.GoType
}

// NewFieldType creates a new field type.
func NewFieldType(typ database.ColumnType) *FieldType {
	f := &FieldType{
		BaseFieldType: database.NewBaseFieldType(typ),
	}
	switch f.DatabaseTypeName() {
	case "BIGINT", "INTEGER", "SMALLINT":
		f.goType = database.GoTypeInt64
	case "BLOB", "CLOB":
		f.goType = database.GoTypeBytes
	case "DOUBLE", "REAL":
		f.goType = database.GoTypeFloat64
	case "DATE", "TIME", "TIMESTAMP":
		f.goType = database.GoTypeTime
	case "BOOLEAN":
		f.goType = database.GoTypeBool
	case "VARCHAR", "CHAR", "DECIMAL":
		f.goType = database.GoTypeString
	}
	return f
}

// IsSupported determines whether parsing is supported.
func (f *FieldType) IsSupported() bool {
	return f.goType != database.GoTypeUnknown
}

// GoType returns the Golang type used when processing numeric values.
func (f *FieldType) GoType() database.GoType {
	return f.goType
}

// Scanner is used for scanning data based on column types.
type Scanner struct {
	f *Field
	database.BaseScanner
}

// NewScanner generates a scanner based on the column type.
func NewScanner(f *Field) *Scanner {
	return &Scanner{
		f: f,
	}
}

// Scan reads data based on the column type.
// INTEGER, BIGINT, and SMALLINT are treated as integers.
// DOUBLE, REAL, and DECIMAL are treated as high-precision real numbers.
// DATE, TIME, and TIMESTAMP are treated as time values.
// CHAR and VARCHAR are treated as strings.
// BLOB is treated as a byte stream.
// BOOLEAN is treated as a boolean value.
func (s *Scanner) Scan(src interface{}) (err error) {
	defer s.f.SetError(&err)
	var cv element.ColumnValue
	byteSize := element.ByteSize(src)
	switch s.f.Type().DatabaseTypeName() {
	case "BIGINT", "INTEGER", "SMALLINT":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBigIntColumnValue()
		case int64:
			v := data
			cv = element.NewBigIntColumnValueFromInt64(v)
		case int32:
			v := int64(data)
			cv = element.NewBigIntColumnValueFromInt64(v)
		case int16:
			v := int64(data)
			cv = element.NewBigIntColumnValueFromInt64(v)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeBigInt)
		}
	case "BLOB", "CLOB":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBytesColumnValue()
		case []byte:
			cv = element.NewBytesColumnValue(data)
		default:
			return fmt.Errorf("src is %v(%T),but not %v", src, src, element.TypeBytes)
		}
	case "DATE":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilTimeColumnValue()
		case time.Time:
			cv = element.NewTimeColumnValueWithDecoder(data, element.NewStringTimeDecoder(dateLayout))
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeTime)
		}
	case "TIME":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilTimeColumnValue()
		case time.Time:
			cv = element.NewTimeColumnValueWithDecoder(data, element.NewStringTimeDecoder(timeLayout))
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeTime)
		}
	case "TIMESTAMP":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilTimeColumnValue()
		case time.Time:
			cv = element.NewTimeColumnValueWithDecoder(data, element.NewStringTimeDecoder(timestampLayout))
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeTime)
		}
	case "CHAR", "VARCHAR":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilStringColumnValue()
		case []byte:
			switch s.f.Type().DatabaseTypeName() {
			case "CHAR":
				data = s.f.TrimByteChar(data)
			}

			var buf []byte
			buf, err = decodeChinese(data)
			if err != nil {
				return err
			}
			cv = element.NewStringColumnValue(string(buf))
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeString)
		}
	case "DOUBLE", "REAL", "DECIMAL":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilDecimalColumnValue()
		case []byte:
			if cv, err = element.NewDecimalColumnValueFromString(string(data)); err != nil {
				return
			}
		case float64:
			cv = element.NewDecimalColumnValueFromFloat(data)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeDecimal)
		}
	case "BOOLEAN":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBoolColumnValue()
		case bool:
			cv = element.NewBoolColumnValue(data)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeDecimal)
		}
	default:
		return fmt.Errorf("src is %v(%T), but db type is %v", src, src, s.f.Type().DatabaseTypeName())
	}
	s.SetColumn(element.NewDefaultColumn(cv, s.f.Name(), byteSize))
	return
}
