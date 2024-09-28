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

package sqlserver

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

var (
	dateLayout     = element.DefaultTimeFormat[:10]
	datetimeLayout = element.DefaultTimeFormat[:26]
)

// Field - Represents a field in a database table.
type Field struct {
	database.BaseConfigSetter

	*database.BaseField
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
	return fmt.Sprintf("@p%d", i)
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
	return NewValuer(f, c)
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
	case "BIT":
		f.goType = database.GoTypeBool
	case "TINYINT", "SMALLINT", "INT", "BIGINT":
		f.goType = database.GoTypeInt64
	case "REAL", "FLOAT":
		f.goType = database.GoTypeFloat64
	case "DECIMAL",
		"VARCHAR", "NVARCHAR", "CHAR", "NCHAR", "TEXT", "NTEXT":
		f.goType = database.GoTypeString
	case "SMALLDATETIME", "DATETIME", "DATETIME2", "DATE", "TIME", "DATETIMEOFFSET":
		f.goType = database.GoTypeTime
	case "VARBINARY", "BINARY":
		f.goType = database.GoTypeBytes
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
	f *Field
	database.BaseScanner
}

// NewScanner - Generates a scanner based on the column type.
func NewScanner(f *Field) *Scanner {
	return &Scanner{
		f: f,
	}
}

// Scan - Reads data from a column based on its type.
func (s *Scanner) Scan(src any) (err error) {
	var cv element.ColumnValue
	byteSize := element.ByteSize(src)
	switch s.f.Type().DatabaseTypeName() {
	case "BIT":
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
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeBigInt)
		}
	case "REAL", "FLOAT", "DECIMAL":
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
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeDecimal)
		}
	case "VARCHAR", "NVARCHAR", "CHAR", "NCHAR", "TEXT", "NTEXT":
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
	case "VARBINARY", "BINARY":
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
	case "SMALLDATETIME", "DATETIME", "DATETIME2", "TIME", "DATETIMEOFFSET":
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

// Valuer - Assigns values to a field.
type Valuer struct {
	f *Field
	c element.Column
}

// NewValuer - Creates a new valuer.
func NewValuer(f *Field, c element.Column) *Valuer {
	return &Valuer{
		f: f,
		c: c,
	}
}

// Value - Represents the value assigned to a field.
func (v *Valuer) Value() (driver.Value, error) {
	// Cannot directly use nil. In Golang, []byte(nil) has a type of []byte but a value of nil, which can cause the following error:
	// mssql: Implicit conversion from data type nvarchar to binary is not allowed.
	// Use the CONVERT function to run this query.
	// The reason is that passing nil results in a TypeId of typeNull in mssql.go's makeParam, which leads to makeDecl returning nvarchar(1).
	if v.c.IsNil() {
		switch v.f.Type().(*FieldType).GoType() {
		case database.GoTypeBytes:
			return []byte(nil), nil
		}
	}

	return database.NewGoValuer(v.f, v.c).Value()
}
