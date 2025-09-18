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

package dm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

// Field represents a field in a DM database table.
type Field struct {
	*database.BaseField
	database.BaseConfigSetter
}

var (
	dateLayout     = element.DefaultTimeFormat[:10]
	datetimeLayout = element.DefaultTimeFormat[:26]
)

// NewField generates a field based on basic column attributes.
func NewField(bf *database.BaseField) *Field {
	return &Field{
		BaseField: bf,
	}
}

// Quoted is used for quoting in SQL statements.
func (f *Field) Quoted() string {
	return Quoted(f.Name())
}

// BindVar is the SQL placeholder used in SQL statements.
func (f *Field) BindVar(i int) string {
	return ":" + strconv.Itoa(i)
}

// Select represents a field for querying purposes in SQL query statements.
func (f *Field) Select() string {
	return f.Quoted()
}

// Type represents the type of the field.
func (f *Field) Type() database.FieldType {
	return NewFieldType(f.FieldType())
}

// Scanner is used for reading data from a field.
func (f *Field) Scanner() database.Scanner {
	return NewScanner(f)
}

// Valuer handles data processing using GoValuer.
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
	// DM数据类型映射需要根据实际的DM类型定义进行调整
	switch f.DatabaseTypeName() {
	case "BOOLEAN", "BIT":
		f.goType = database.GoTypeBool
	case "TINYINT", "SMALLINT", "INT", "BIGINT":
		f.goType = database.GoTypeInt64
	case "FLOAT", "DOUBLE", "DECIMAL", "NUMERIC", "NUMBER":
		f.goType = database.GoTypeString // DECIMAL使用字符串以保证精度
	case "CHAR", "VARCHAR", "VARCHAR2", "TEXT", "CLOB":
		f.goType = database.GoTypeString
	case "BLOB", "BFILE", "BINARY", "VARBINARY":
		f.goType = database.GoTypeBytes
	case "DATE", "TIME", "DATETIME", "TIMESTAMP":
		f.goType = database.GoTypeTime
	default:
		f.goType = database.GoTypeString
	}
	return f
}

// GoType returns the basic type of Go.
func (f *FieldType) GoType() database.GoType {
	return f.goType
}

// Scanner represents a scanner for reading data from a DM database.
type Scanner struct {
	f *Field
	database.BaseScanner
}

// NewScanner creates a scanner based on a field.
func NewScanner(f *Field) *Scanner {
	return &Scanner{
		f: f,
	}
}

// Scan converts data from the database to the corresponding value based on the field type.
func (s *Scanner) Scan(src any) (err error) {
	var cv element.ColumnValue
	byteSize := element.ByteSize(src)
	switch s.f.Type().DatabaseTypeName() {
	case "BIT", "BOOLEAN":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBoolColumnValue()
		case bool:
			cv = element.NewBoolColumnValue(data)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeBool)
		}

	case "TINYINT", "SMALLINT", "INT", "BIGINT":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBigIntColumnValue()
		case int64:
			cv = element.NewBigIntColumnValueFromInt64(data)
		case string:
			if cv, err = element.NewBigIntColumnValueFromString(data); err != nil {
				return
			}
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeBigInt)
		}
	case "FLOAT", "DOUBLE", "DECIMAL", "NUMERIC", "NUMBER":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilDecimalColumnValue()
		case float32:
			cv = element.NewDecimalColumnValueFromFloat32(data)
		case float64:
			cv = element.NewDecimalColumnValueFromFloat(data)
		case []byte:
			if cv, err = element.NewDecimalColumnValueFromString(string(data)); err != nil {
				return
			}
		case string:
			if cv, err = element.NewDecimalColumnValueFromString(data); err != nil {
				return
			}
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeDecimal)
		}
	case "CHAR", "VARCHAR", "VARCHAR2", "TEXT", "CLOB":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilStringColumnValue()
		case string:
			switch s.f.Type().DatabaseTypeName() {
			case "CHAR", "NCHAR":
				data = s.f.TrimStringChar(data)
			}
			cv = element.NewStringColumnValue(data)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeString)
		}
	case "BLOB", "BFILE", "BINARY", "VARBINARY":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBytesColumnValue()
		case []byte:
			cv = element.NewBytesColumnValueNoCopy(data)
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
	case "TIME", "DATETIME", "TIMESTAMP":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilTimeColumnValue()
		case time.Time:
			cv = element.NewTimeColumnValueWithDecoder(data, element.NewStringTimeDecoder(datetimeLayout))
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeTime)
		}
	default:
		return fmt.Errorf("src is %v(%T), but db type is %v", src, src, s.f.Type().DatabaseTypeName())
	}
	s.SetColumn(element.NewDefaultColumn(cv, s.f.Name(), byteSize))
	return
}
