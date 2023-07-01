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
	"golang.org/x/text/encoding/simplifiedchinese"
)

var (
	dateLayout      = element.DefaultTimeFormat[:10]
	timestampLayout = element.DefaultTimeFormat[:26]
	timeLayout      = timestampLayout
)

//Field 字段
type Field struct {
	*database.BaseField
}

//NewField 通过基本列属性生成字段
func NewField(bf *database.BaseField) *Field {
	return &Field{
		BaseField: bf,
	}
}

//Quoted 引用，用于SQL语句
func (f *Field) Quoted() string {
	return Quoted(f.Name())
}

//BindVar SQL占位符，用于SQL语句
func (f *Field) BindVar(_ int) string {
	return "?"
}

//Select 查询时字段，用于SQL查询语句
func (f *Field) Select() string {
	return Quoted(f.Name())
}

//Type 字段类型
func (f *Field) Type() database.FieldType {
	return NewFieldType(f.FieldType())
}

//Scanner 扫描器，用于读取数据
func (f *Field) Scanner() database.Scanner {
	return NewScanner(f)
}

//Valuer 赋值器，采用GoValuer处理数据
func (f *Field) Valuer(c element.Column) database.Valuer {
	return database.NewGoValuer(f, c)
}

//FieldType 字段类型
type FieldType struct {
	*database.BaseFieldType

	goType database.GoType
}

//NewFieldType 创建新的字段类型
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

//IsSupportted 是否支持解析
func (f *FieldType) IsSupportted() bool {
	return f.goType != database.GoTypeUnknown
}

//GoType 返回处理数值时的Golang类型
func (f *FieldType) GoType() database.GoType {
	return f.goType
}

//Scanner 扫描器
type Scanner struct {
	f *Field
	database.BaseScanner
}

//NewScanner 根据列类型生成扫描器
func NewScanner(f *Field) *Scanner {
	return &Scanner{
		f: f,
	}
}

//Scan 根据列类型读取数据
//"INTEGER", "BIGINT", "SMALLINT"作为整形处理
//"DOUBLE", "REAL", "DECIMAL"作为高精度实数处理
//"DATE", "TIME", "TIMESTAMP" 作为时间处理
//"CHAR", "VARCHAR"作为字符串处理
//"BLOB" 作为字节流处理
//"BOOLEAN" 作为布尔值处理
func (s *Scanner) Scan(src interface{}) (err error) {
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
			var v []byte
			v, err = simplifiedchinese.GBK.NewDecoder().Bytes(data)
			if err != nil {
				return err
			}
			cv = element.NewStringColumnValue(string(v))
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
