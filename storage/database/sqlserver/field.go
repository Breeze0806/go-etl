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
	"github.com/shopspring/decimal"
)

var (
	dateLayout     = element.DefaultTimeFormat[:10]
	datetimeLayout = element.DefaultTimeFormat[:26]
)

// Field 字段
type Field struct {
	*database.BaseField
}

// NewField 通过基本列属性生成字段
func NewField(bf *database.BaseField) *Field {
	return &Field{
		BaseField: bf,
	}
}

// Quoted 引用，用于SQL语句
func (f *Field) Quoted() string {
	return Quoted(f.Name())
}

// BindVar SQL占位符，用于SQL语句
func (f *Field) BindVar(i int) string {
	return fmt.Sprintf("@p%d", i)
}

// Select 查询时字段，用于SQL查询语句
func (f *Field) Select() string {
	return f.Quoted()
}

// Type 字段类型
func (f *Field) Type() database.FieldType {
	return NewFieldType(f.FieldType())
}

// Scanner 扫描器，用于读取数据
func (f *Field) Scanner() database.Scanner {
	return NewScanner(f)
}

// Valuer 赋值器，采用GoValuer处理数据
func (f *Field) Valuer(c element.Column) database.Valuer {
	return NewValuer(f, c)
}

// FieldType 字段类型
type FieldType struct {
	*database.BaseFieldType

	goType database.GoType
}

// NewFieldType 创建新的字段类型
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

// IsSupportted 是否支持解析
func (f *FieldType) IsSupportted() bool {
	return f.GoType() != database.GoTypeUnknown
}

// GoType 返回处理数值时的Golang类型
func (f *FieldType) GoType() database.GoType {
	return f.goType
}

// Scanner 扫描器
type Scanner struct {
	f *Field
	database.BaseScanner
}

// NewScanner 根据列类型生成扫描器
func NewScanner(f *Field) *Scanner {
	return &Scanner{
		f: f,
	}
}

// Scan 根据列类型读取数据
func (s *Scanner) Scan(src interface{}) (err error) {
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
			cv = element.NewDecimalColumnValue(decimal.NewFromFloat32(data))
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

// Valuer 赋值器
type Valuer struct {
	f *Field
	c element.Column
}

// NewValuer 创建新赋值器
func NewValuer(f *Field, c element.Column) *Valuer {
	return &Valuer{
		f: f,
		c: c,
	}
}

// Value 赋值
func (v *Valuer) Value() (driver.Value, error) {
	//不能直接nil，golang的[]byte(nil)的类型是[]byte，但是值是nil，会导致以下错误:
	//mssql: Implicit conversion from data type nvarchar to binary is not allowed.
	//Use the CONVERT function to run this query.
	//原因是传入nil，在mssql.go的makeParam时TypeId是typeNull，导致makeDecl返回"nvarchar(1)"
	if v.c.IsNil() {
		switch v.f.Type().(*FieldType).GoType() {
		case database.GoTypeBytes:
			return []byte(nil), nil
		}
	}

	return database.NewGoValuer(v.f, v.c).Value()
}
