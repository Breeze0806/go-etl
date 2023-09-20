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
func (f *Field) BindVar(_ int) string {
	return "?"
}

// Select 查询时字段，用于SQL查询语句
func (f *Field) Select() string {
	return Quoted(f.Name())
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
	return database.NewGoValuer(f, c)
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
	//由于存在非负整数，如果直接变为对应的int类型，则会导致转化错误
	//TIME存在负数无法正常转化，YEAR就是TINYINT
	//todo: test YEAR
	case "MEDIUMINT", "INT", "BIGINT", "SMALLINT", "TINYINT",
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
// "MEDIUMINT", "INT", "BIGINT", "SMALLINT", "TINYINT", "YEAR"作为整形处理
// "DOUBLE", "FLOAT", "DECIMAL"作为高精度实数处理
// "DATE", "DATETIME", "TIMESTAMP" 作为时间处理
// "TEXT", "LONGTEXT", "MEDIUMTEXT", "TINYTEXT", "CHAR", "VARCHAR", "TIME"作为字符串处理
// "BLOB", "LONGBLOB", "MEDIUMBLOB", "BINARY", "TINYBLOB", "VARBINARY"作为字节流处理
func (s *Scanner) Scan(src interface{}) (err error) {
	var cv element.ColumnValue
	byteSize := element.ByteSize(src)

	switch s.f.Type().DatabaseTypeName() {
	//todo: test year
	case "MEDIUMINT", "INT", "BIGINT", "SMALLINT", "TINYINT", "YEAR":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBigIntColumnValue()
		case []byte:
			if cv, err = element.NewBigIntColumnValueFromString(string(data)); err != nil {
				return
			}
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
			cv = element.NewDecimalColumnValue(decimal.NewFromFloat32(data))
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
