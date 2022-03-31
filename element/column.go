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

package element

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

//ColumnType 列类型
type ColumnType string

//列类型枚举
const (
	TypeUnknown ColumnType = "unknown" //未知类型
	TypeBool    ColumnType = "bool"    //布尔类型
	TypeBigInt  ColumnType = "bigInt"  //整数类型
	TypeDecimal ColumnType = "decimal" //高精度实数类型
	TypeString  ColumnType = "string"  //字符串类型
	TypeBytes   ColumnType = "bytes"   //字节流类型
	TypeTime    ColumnType = "time"    //时间类型
)

//String 打印显示
func (c ColumnType) String() string {
	return string(c)
}

//ColumnValue 列值
type ColumnValue interface {
	fmt.Stringer

	Type() ColumnType                    //列类型
	IsNil() bool                         //是否为空
	AsBool() (bool, error)               //转化为布尔值
	AsBigInt() (*big.Int, error)         //转化为整数
	AsDecimal() (decimal.Decimal, error) //转化为高精度实数
	AsString() (string, error)           //转化为字符串
	AsBytes() ([]byte, error)            //转化为字节流
	AsTime() (time.Time, error)          // 转化为时间
}

//ColumnValueClonable 可克隆列值
type ColumnValueClonable interface {
	Clone() ColumnValue //克隆
}

//ColumnValueComparabale 可比较列值
type ColumnValueComparabale interface {
	//比较 1代表大于， 0代表相等， -1代表小于
	Cmp(ColumnValue) (int, error)
}

//Column 列
type Column interface {
	ColumnValue
	AsInt8() (int8, error)       //转化为8位整数
	AsInt16() (int16, error)     //转化为16位整数
	AsInt32() (int32, error)     //转化为32位整数
	AsInt64() (int64, error)     //转化为64位整数
	AsFloat32() (float32, error) //转化为32位实数
	AsFloat64() (float64, error) //转化为64位实数
	Clone() (Column, error)      //克隆
	Cmp(Column) (int, error)     //比较, 1代表大于， 0代表相等， -1代表小于
	Name() string                //列名
	ByteSize() int64             //字节流大小
	MemorySize() int64           //内存大小
}

type notNilColumnValue struct{}

//IsNil  是否为空
func (n *notNilColumnValue) IsNil() bool {
	return false
}

type nilColumnValue struct{}

//Type  列类型
func (n *nilColumnValue) Type() ColumnType {
	return TypeUnknown
}

//IsNil  是否为空
func (n *nilColumnValue) IsNil() bool {
	return true
}

//AsBool 无法转化布尔值
func (n *nilColumnValue) AsBool() (bool, error) {
	return false, ErrNilValue
}

//AsBigInt 无法转化整数
func (n *nilColumnValue) AsBigInt() (*big.Int, error) {
	return nil, ErrNilValue
}

//AsDecimal 无法转化高精度实数
func (n *nilColumnValue) AsDecimal() (decimal.Decimal, error) {
	return decimal.Decimal{}, ErrNilValue
}

//AsString 无法转化字符串
func (n *nilColumnValue) AsString() (string, error) {
	return "", ErrNilValue
}

//AsBytes 无法转化字节流
func (n *nilColumnValue) AsBytes() ([]byte, error) {
	return nil, ErrNilValue
}

//AsTime 无法转化时间
func (n *nilColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, ErrNilValue
}

//String 打印显示
func (n *nilColumnValue) String() string {
	return "<nil>"
}

//DefaultColumn 默认值
type DefaultColumn struct {
	ColumnValue // 列值

	name     string
	byteSize int
}

//NewDefaultColumn 根据列值v,列名name,字节流大小byteSiz，生成默认列
func NewDefaultColumn(v ColumnValue, name string, byteSize int) Column {
	return &DefaultColumn{
		ColumnValue: v,
		name:        name,
		byteSize:    byteSize,
	}
}

//Name 列名
func (d *DefaultColumn) Name() string {
	return d.name
}

//Cmp 比较列，如果不是可比较列值，就会报错
func (d *DefaultColumn) Cmp(c Column) (int, error) {
	if d.Name() != c.Name() {
		return 0, ErrColumnNameNotEqual
	}
	comparabale, ok := d.ColumnValue.(ColumnValueComparabale)
	if !ok {
		return 0, ErrNotColumnValueComparable
	}
	return comparabale.Cmp(c)
}

//Clone 克隆列，如果不是可克隆列值，就会报错
func (d *DefaultColumn) Clone() (Column, error) {
	colnable, ok := d.ColumnValue.(ColumnValueClonable)
	if !ok {
		return nil, ErrNotColumnValueClonable
	}

	return &DefaultColumn{
		ColumnValue: colnable.Clone(),
		name:        d.name,
		byteSize:    d.byteSize,
	}, nil
}

//ByteSize 字节流大小
func (d *DefaultColumn) ByteSize() int64 {
	return int64(d.byteSize)
}

//MemorySize 内存大小
func (d *DefaultColumn) MemorySize() int64 {
	return int64(d.byteSize + len(d.name) + 4)
}

//AsInt8 转化为8位整数
func (d *DefaultColumn) AsInt8() (int8, error) {
	bi, err := d.AsBigInt()
	if err != nil {
		return 0, NewTransformErrorFormString(d.Type().String(), "int8", err)
	}
	if bi.IsInt64() {
		v := bi.Int64()
		if v > math.MaxInt8 || v < math.MinInt8 {
			return 0, NewTransformErrorFormString(d.Type().String(), "int8", fmt.Errorf("%v %v", v, strconv.ErrRange))
		}
		return int8(bi.Int64()), nil
	}

	return 0, NewTransformErrorFormString(d.Type().String(), "int8", fmt.Errorf("%v %v", d.String(), ErrValueNotInt64))
}

//AsInt16 转化为16位整数
func (d *DefaultColumn) AsInt16() (int16, error) {
	bi, err := d.AsBigInt()
	if err != nil {
		return 0, NewTransformErrorFormString(d.Type().String(), "int16", err)
	}
	if bi.IsInt64() {
		v := bi.Int64()
		if v > math.MaxInt16 || v < math.MinInt16 {
			return 0, NewTransformErrorFormString(d.Type().String(), "int16", fmt.Errorf("%v %v", v, strconv.ErrRange))
		}
		return int16(bi.Int64()), nil
	}
	return 0, NewTransformErrorFormString(d.Type().String(), "int16", fmt.Errorf("%v %v", d.String(), ErrValueNotInt64))
}

//AsInt32 转化为32位整数
func (d *DefaultColumn) AsInt32() (int32, error) {
	bi, err := d.AsBigInt()
	if err != nil {
		return 0, NewTransformErrorFormString(d.Type().String(), "int32", err)
	}
	if bi.IsInt64() {
		v := bi.Int64()
		if v > math.MaxInt32 || v < math.MinInt32 {
			return 0, NewTransformErrorFormString(d.Type().String(), "int32", fmt.Errorf("%v %v", v, strconv.ErrRange))
		}
		return int32(bi.Int64()), nil
	}
	return 0, NewTransformErrorFormString(d.Type().String(), "int32", fmt.Errorf("%v %v", d.String(), ErrValueNotInt64))
}

//AsInt64 转化为64位整数
func (d *DefaultColumn) AsInt64() (int64, error) {
	bi, err := d.AsBigInt()
	if err != nil {
		return 0, NewTransformErrorFormString(d.Type().String(), "int64", err)
	}
	if bi.IsInt64() {
		return int64(bi.Int64()), nil
	}
	return 0, NewTransformErrorFormString(d.Type().String(), "int64", fmt.Errorf("%v %v", d.String(), ErrValueNotInt64))
}

//AsFloat32 转化为32位实数
func (d *DefaultColumn) AsFloat32() (float32, error) {
	dec, err := d.AsDecimal()
	if err != nil {
		return 0, NewTransformErrorFormString(d.Type().String(), "float32", err)
	}
	v, _ := dec.Rat().Float32()
	if math.Abs(float64(v)) > math.MaxFloat32 {
		return 0, NewTransformErrorFormString(d.Type().String(), "float32",
			fmt.Errorf("%v %v", d.String(), strconv.ErrRange))
	}
	return v, nil
}

//AsFloat64 转化为64位实数
func (d *DefaultColumn) AsFloat64() (float64, error) {
	dec, err := d.AsDecimal()
	if err != nil {
		return 0, NewTransformErrorFormString(d.Type().String(), "float64", err)
	}
	v, _ := dec.Float64()
	if math.Abs(v) > math.MaxFloat64 {
		return 0, NewTransformErrorFormString(d.Type().String(), "float64",
			fmt.Errorf("%v %v", d.String(), strconv.ErrRange))
	}
	return v, nil
}
