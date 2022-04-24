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

package postgres

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/lib/pq/oid"
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
func (f *Field) BindVar(i int) string {
	return "$" + strconv.Itoa(i)
}

//Select 查询时字段，用于SQL查询语句
func (f *Field) Select() string {
	return f.Quoted()
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
	case oid.TypeName[oid.T_bool]:
		f.goType = database.GoTypeBool
	case oid.TypeName[oid.T_int2], oid.TypeName[oid.T_int4],
		oid.TypeName[oid.T_int8]:
		f.goType = database.GoTypeInt64
	case oid.TypeName[oid.T_float4], oid.TypeName[oid.T_float8]:
		f.goType = database.GoTypeFloat64
	case oid.TypeName[oid.T_varchar], oid.TypeName[oid.T_text],
		oid.TypeName[oid.T_numeric]:
		f.goType = database.GoTypeString
	case oid.TypeName[oid.T_date], oid.TypeName[oid.T_time],
		oid.TypeName[oid.T_timetz], oid.TypeName[oid.T_timestamp],
		oid.TypeName[oid.T_timestamptz]:
		f.goType = database.GoTypeTime
	case oid.TypeName[oid.T_bpchar]:
		f.goType = database.GoTypeString
	}
	return f
}

//IsSupportted 是否支持解析
func (f *FieldType) IsSupportted() bool {
	return f.GoType() != database.GoTypeUnknown
}

//GoType 返回处理数值时的Golang类型
func (f *FieldType) GoType() database.GoType {
	return f.goType
}

//Scanner 扫描器
type Scanner struct {
	database.BaseScanner

	f *Field
}

//NewScanner 根据列类型生成扫描器
func NewScanner(f *Field) *Scanner {
	return &Scanner{
		f: f,
	}
}

//Scan 根据列类型读取数据
func (s *Scanner) Scan(src interface{}) (err error) {
	var cv element.ColumnValue
	//todo: byteSize is 0, fix it
	var byteSize int
	switch s.f.Type().DatabaseTypeName() {
	case oid.TypeName[oid.T_bool]:
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBoolColumnValue()
		case bool:
			cv = element.NewBoolColumnValue(data)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeBool)
		}
	case oid.TypeName[oid.T_int2], oid.TypeName[oid.T_int4],
		oid.TypeName[oid.T_int8]:
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBigIntColumnValue()
		case int64:
			cv = element.NewBigIntColumnValueFromInt64(data)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeBigInt)
		}
	case oid.TypeName[oid.T_bpchar]:
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBytesColumnValue()
		case []byte:
			cv = element.NewBytesColumnValue(data)
		default:
			return fmt.Errorf("src is %v(%T),but not %v", src, src, element.TypeBytes)
		}
	case oid.TypeName[oid.T_date]:
		switch data := src.(type) {
		case nil:
			cv = element.NewNilTimeColumnValue()
		case time.Time:
			cv = element.NewTimeColumnValueWithDecoder(data, element.NewStringTimeDecoder("2006-01-02"))
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeTime)
		}

	case oid.TypeName[oid.T_time], oid.TypeName[oid.T_timetz],
		oid.TypeName[oid.T_timestamp], oid.TypeName[oid.T_timestamptz]:
		switch data := src.(type) {
		case nil:
			cv = element.NewNilTimeColumnValue()
		case time.Time:
			cv = element.NewTimeColumnValueWithDecoder(data, element.NewStringTimeDecoder("2006-01-02 15:04:05"))
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeTime)
		}
	case oid.TypeName[oid.T_varchar], oid.TypeName[oid.T_text]:
		switch data := src.(type) {
		case nil:
			cv = element.NewNilStringColumnValue()
		case string:
			cv = element.NewStringColumnValue(data)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeString)
		}
	case oid.TypeName[oid.T_float4],
		oid.TypeName[oid.T_float8], oid.TypeName[oid.T_numeric]:
		switch data := src.(type) {
		case nil:
			cv = element.NewNilDecimalColumnValue()
		case float64:
			cv = element.NewDecimalColumnValueFromFloat(data)
		case []byte:
			if cv, err = element.NewDecimalColumnValueFromString(string(data)); err != nil {
				return
			}
		default:
			return fmt.Errorf("src is %v(%T), but type is %v", src, src, element.TypeDecimal)
		}
	default:
		return fmt.Errorf("src is %v(%T), but db type is %v", src, src, s.f.Type().DatabaseTypeName())
	}
	s.SetColumn(element.NewDefaultColumn(cv, s.f.Name(), byteSize))
	return
}
