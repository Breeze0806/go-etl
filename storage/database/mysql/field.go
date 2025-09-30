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

package mysql

import (
	"fmt"
	"time"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

var (
	dateLayout     = element.DefaultTimeFormat[:10]
	datetimeLayout = element.DefaultTimeFormat[:26]
)

// Field Field
type Field struct {
	*database.BaseField
	database.BaseConfigSetter
}

// NewField Generate a field based on basic column attributes
func NewField(bf *database.BaseField) *Field {
	return &Field{
		BaseField: bf,
	}
}

// Quoted Quotation, used in SQL statements
func (f *Field) Quoted() string {
	return Quoted(f.Name())
}

// BindVar SQL placeholder, used in SQL statements
func (f *Field) BindVar(_ int) string {
	return "?"
}

// Select Field for querying, used in SQL query statements
func (f *Field) Select() string {
	return Quoted(f.Name())
}

// Type Field type
func (f *Field) Type() database.FieldType {
	return NewFieldType(f.FieldType())
}

// Scanner Scanner, used for reading data
func (f *Field) Scanner() database.Scanner {
	return NewScanner(f)
}

// Valuer Valuer, using GoValuer to process data
func (f *Field) Valuer(c element.Column) database.Valuer {
	return database.NewGoValuer(f, c)
}

// FieldType Field type
type FieldType struct {
	*database.BaseFieldType

	goType database.GoType
}

// NewFieldType Create a new field type
func NewFieldType(typ database.ColumnType) *FieldType {
	f := &FieldType{
		BaseFieldType: database.NewBaseFieldType(typ),
	}
	switch f.DatabaseTypeName() {
	// Due to the existence of non-negative integers, directly converting them to the corresponding int type would result in conversion errors.
	// TIME has negative values and cannot be converted normally, while YEAR is TINYINT.
	// todo: test YEAR
	case "MEDIUMINT", "INT", "BIGINT", "SMALLINT", "TINYINT",
		"UNSIGNED INT", "UNSIGNED BIGINT", "UNSIGNED SMALLINT", "UNSIGNED TINYINT",
		"TEXT", "LONGTEXT", "MEDIUMTEXT", "TINYTEXT", "CHAR", "VARCHAR",
		"TIME", "YEAR",
		"DECIMAL":
		f.goType = database.GoTypeString
	case "BLOB", "LONGBLOB", "MEDIUMBLOB", "BINARY", "TINYBLOB", "VARBINARY", "BIT":
		f.goType = database.GoTypeBytes
	case "DOUBLE", "FLOAT":
		f.goType = database.GoTypeFloat64
	case "DATE", "DATETIME", "TIMESTAMP":
		f.goType = database.GoTypeTime
	}
	return f
}

// IsSupported Whether it supports parsing
func (f *FieldType) IsSupported() bool {
	return f.GoType() != database.GoTypeUnknown
}

// GoType Returns the Golang type when processing numeric values
func (f *FieldType) GoType() database.GoType {
	return f.goType
}

// Scanner Scanner
type Scanner struct {
	f *Field
	database.BaseScanner
}

// NewScanner Generate a scanner based on the column type
func NewScanner(f *Field) *Scanner {
	return &Scanner{
		f: f,
	}
}

// Scan Read data based on the column type
// MEDIUMINT, INT, BIGINT, SMALLINT, TINYINT, YEAR, UNSIGNED INT, UNSIGNED BIGINT, UNSIGNED SMALLINT, UNSIGNED TINYINT are treated as integers.
// DOUBLE, FLOAT, DECIMAL are treated as high-precision real numbers.
// DATE, DATETIME, TIMESTAMP are treated as time.
// TEXT, LONGTEXT, MEDIUMTEXT, TINYTEXT, CHAR, VARCHAR, TIME are treated as strings.
// BLOB, LONGBLOB, MEDIUMBLOB, BINARY, TINYBLOB, VARBINARY are treated as byte streams.
func (s *Scanner) Scan(src any) (err error) {
	defer s.f.SetError(&err)
	var cv element.ColumnValue
	byteSize := element.ByteSize(src)

	switch s.f.Type().DatabaseTypeName() {
	// todo: test year
	case "MEDIUMINT", "INT", "BIGINT", "SMALLINT", "TINYINT", "YEAR",
		"UNSIGNED INT", "UNSIGNED BIGINT", "UNSIGNED SMALLINT", "UNSIGNED TINYINT":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBigIntColumnValue()
		case uint64:
			cv = element.NewBigIntColumnValueFromUint64(data)
		case int64:
			cv = element.NewBigIntColumnValueFromInt64(data)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeBigInt)
		}
	case "BLOB", "LONGBLOB", "MEDIUMBLOB", "BINARY", "TINYBLOB", "VARBINARY", "BIT":
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
	case "DATETIME", "TIMESTAMP":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilTimeColumnValue()
		case time.Time:
			cv = element.NewTimeColumnValueWithDecoder(data, element.NewStringTimeDecoder(datetimeLayout))
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeTime)
		}
	case "TEXT", "LONGTEXT", "MEDIUMTEXT", "TINYTEXT", "CHAR", "VARCHAR", "TIME":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilStringColumnValue()
		case []byte:
			switch s.f.Type().DatabaseTypeName() {
			case "CHAR":
				data = s.f.TrimByteChar(data)
			}
			cv = element.NewStringColumnValue(string(data))
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeString)
		}
	case "DOUBLE", "FLOAT", "DECIMAL":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilDecimalColumnValue()
		case []byte:
			if cv, err = element.NewDecimalColumnValueFromString(string(data)); err != nil {
				return
			}
		case float32:
			cv = element.NewDecimalColumnValue(element.NewFromFloat32(data))
		case float64:
			cv = element.NewDecimalColumnValueFromFloat(data)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeDecimal)
		}
	default:
		return fmt.Errorf("src is %v(%T), but db type is %v", src, src, s.f.Type().DatabaseTypeName())
	}
	s.SetColumn(element.NewDefaultColumn(cv, s.f.Name(), byteSize))
	return
}
