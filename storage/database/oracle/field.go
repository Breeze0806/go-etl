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
	return ":" + strconv.Itoa(i)
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
	return NewValuer(f, c)
}

//FieldType 字段类型
type FieldType struct {
	*database.BaseFieldType

	supportted bool
}

//NewFieldType 创建新的字段类型
func NewFieldType(typ database.ColumnType) *FieldType {
	f := &FieldType{
		BaseFieldType: database.NewBaseFieldType(typ),
	}
	switch f.DatabaseTypeName() {
	//由于oracle特殊的转化机制导致所有的数据需要转化为string类型进行插入
	case "BOOLEAN",
		"BINARY_INTEGER",
		"NUMBER", "FLOAT", "DOUBLE",
		"TIMESTAMP", "TIMESTAMP WITH TIME ZONE", "TIMESTAMP WITH LOCAL TIME ZONE", "DATE",
		"VARCHAR2", "NVARCHAR2", "CHAR", "NCHAR", "LONG",
		"CLOB", "NCLOB", "BFILE", "BLOB", "RAW", "LONG RAW":
		f.supportted = true
	}
	return f
}

//IsSupportted 是否支持解析
func (f *FieldType) IsSupportted() bool {
	return f.supportted
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
// "BOOLEAN" 做为bool类型处理
// "BINARY_INTEGER" 做为bigint类型处理
// "NUMBER", "FLOAT", "DOUBLE" 做为decimal类型处理
// "TIMESTAMP", "TIMESTAMP WITH TIME ZONE", "TIMESTAMP WITH LOCAL TIME ZONE", "DATE"做为time类型处理
// "CLOB", "NCLOB", "BFILE", "BLOB", "VARCHAR2", "NVARCHAR2", "CHAR", "NCHAR", "LONG"做为string类型处理
// "RAW", "LONG RAW"做为bytes类型处理
func (s *Scanner) Scan(src interface{}) (err error) {
	var cv element.ColumnValue
	//todo: byteSize is 0, fix it
	var byteSize int
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
	case "RAW", "LONG RAW":
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
	case "CLOB", "NCLOB", "BFILE", "BLOB", "VARCHAR2", "NVARCHAR2", "CHAR", "NCHAR", "LONG":
		switch data := src.(type) {
		case string:
			if data == "" {
				cv = element.NewNilStringColumnValue()
			} else {
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

//Valuer 赋值器
type Valuer struct {
	f *Field
	c element.Column
}

//NewValuer 创建新赋值器
func NewValuer(f *Field, c element.Column) *Valuer {
	return &Valuer{
		f: f,
		c: c,
	}
}

//Value 赋值
func (v *Valuer) Value() (value driver.Value, err error) {
	if v.c.IsNil() {
		return "", nil
	}
	if v.f.Type().DatabaseTypeName() == "BOOLEAN" {
		var b bool
		if b, err = v.c.AsBool(); err != nil {
			return nil, err
		}
		if b {
			return "1", nil
		}
		return "0", nil
	}

	return v.c.AsString()
}
