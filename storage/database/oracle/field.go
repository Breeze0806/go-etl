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

package oracle

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/godror/godror"
	"github.com/shopspring/decimal"
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
func (f *Field) BindVar(i int) string {
	// Fix the time format error ORA-01861: literal does not match format string
	switch f.FieldType().DatabaseTypeName() {
	case "DATE":
		return "to_date(:" + strconv.Itoa(i) + ",'yyyy-mm-dd hh24:mi:ss')"
	case "TIMESTAMP", "TIMESTAMP WITH TIME ZONE", "TIMESTAMP WITH LOCAL TIME ZONE":
		return "to_timestamp(:" + strconv.Itoa(i) + ",'yyyy-mm-dd hh24:mi:ss.ff9')"
	}

	return ":" + strconv.Itoa(i)
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
	return NewValuer(f, c)
}

// FieldType Field type
type FieldType struct {
	*database.BaseFieldType

	supportted bool
}

// NewFieldType Create a new field type
func NewFieldType(typ database.ColumnType) *FieldType {
	f := &FieldType{
		BaseFieldType: database.NewBaseFieldType(typ),
	}
	switch f.DatabaseTypeName() {

	case "BOOLEAN",
		"BINARY_INTEGER",
		"NUMBER", "FLOAT", "DOUBLE",
		"TIMESTAMP", "TIMESTAMP WITH TIME ZONE", "TIMESTAMP WITH LOCAL TIME ZONE", "DATE",
		"VARCHAR2", "NVARCHAR2", "CHAR", "NCHAR", "LONG",
		"CLOB", "NCLOB", "BLOB", "RAW", "LONG RAW":
		f.supportted = true
	}
	return f
}

// IsSupported Whether it supports parsing
func (f *FieldType) IsSupported() bool {
	return f.supportted
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
// BOOLEAN is treated as a bool type
// BINARY_INTEGER is treated as a bigint type
// NUMBER, FLOAT, DOUBLE are treated as decimal types
// TIMESTAMP, TIMESTAMP WITH TIME ZONE, TIMESTAMP WITH LOCAL TIME ZONE, DATE are treated as time types
// CLOB, NCLOB, VARCHAR2, NVARCHAR2, CHAR, NCHAR are treated as string types
// BLOB, RAW, LONG RAW, LONG are treated as byte types
func (s *Scanner) Scan(src interface{}) (err error) {
	var cv element.ColumnValue
	byteSize := element.ByteSize(src)

	switch s.f.Type().DatabaseTypeName() {
	case "BOOLEAN":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBoolColumnValue()
		case bool:
			cv = element.NewBoolColumnValue(data)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeBigInt)
		}
	case "BINARY_INTEGER":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBigIntColumnValue()
		case int64:
			cv = element.NewBigIntColumnValueFromInt64(data)
		case uint64:
			cv = element.NewBigIntColumnValue(new(big.Int).SetUint64(data))
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeBigInt)
		}
	// todo test BFILE
	case // BFILE,
		"BLOB", "LONG", "RAW", "LONG RAW":
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
	case "TIMESTAMP", "TIMESTAMP WITH TIME ZONE", "TIMESTAMP WITH LOCAL TIME ZONE":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilTimeColumnValue()
		case time.Time:
			cv = element.NewTimeColumnValueWithDecoder(data, element.NewStringTimeDecoder(datetimeLayout))
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeTime)
		}
	case "CLOB", "NCLOB", "VARCHAR2", "NVARCHAR2", "CHAR", "NCHAR":
		switch data := src.(type) {
		case string:
			if data == "" {
				cv = element.NewNilStringColumnValue()
			} else {
				switch s.f.Type().DatabaseTypeName() {
				case "CHAR", "NCHAR":
					data = s.f.TrimStringChar(data)
				}
				cv = element.NewStringColumnValue(data)
			}
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeString)
		}
	case "NUMBER", "FLOAT", "DOUBLE":
		s := ""
		switch data := src.(type) {
		case nil:
			cv = element.NewNilDecimalColumnValue()
		case float32:
			cv = element.NewDecimalColumnValue(decimal.NewFromFloat32(data))
		case float64:
			cv = element.NewDecimalColumnValueFromFloat(data)
		case int64:
			s = strconv.FormatInt(data, 10)
		case uint64:
			s = strconv.FormatUint(data, 10)
		case bool:
			s = "0"
			if data {
				s = "1"
			}
		case godror.Number:
			s = string(data)
			byteSize = len(s)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeDecimal)
		}
		if s != "" {
			if cv, err = element.NewDecimalColumnValueFromString(s); err != nil {
				return
			}
		}
	default:
		return fmt.Errorf("src is %v(%T), but db type is %v", src, src, s.f.Type().DatabaseTypeName())
	}
	s.SetColumn(element.NewDefaultColumn(cv, s.f.Name(), byteSize))
	return
}

// Valuer Valuer
type Valuer struct {
	f *Field
	c element.Column
}

// NewValuer Create a new valuer
func NewValuer(f *Field, c element.Column) *Valuer {
	return &Valuer{
		f: f,
		c: c,
	}
}

// Value Assignment
func (v *Valuer) Value() (value driver.Value, err error) {
	switch v.f.Type().DatabaseTypeName() {
	case "BOOLEAN":
		// In Oracle, inserting an empty string is actually treated as nil, corresponding to NULL
		if v.c.IsNil() {
			return "", nil
		}
		var b bool
		if b, err = v.c.AsBool(); err != nil {
			return nil, err
		}
		if b {
			return "1", nil
		}
		return "0", nil
		// todo test BFILE
	case // BFILE,
		"BLOB", "LONG", "RAW", "LONG RAW":
		// For these types, inserting nil corresponds to NULL
		if v.c.IsNil() {
			return nil, nil
		}
		return v.c.AsBytes()
	}
	// In Oracle, inserting an empty string is actually treated as nil, corresponding to NULL
	if v.c.IsNil() {
		return "", nil
	}
	// Due to Oracle's special conversion mechanism, all data needs to be converted to string type for insertion
	return v.c.AsString()
}
