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
	"strconv"
	"time"
)

// NilBytesColumnValue 空值字节流列值
type NilBytesColumnValue struct {
	*nilColumnValue
}

// NewNilBytesColumnValue 创建空值字节流列值
func NewNilBytesColumnValue() ColumnValue {
	return &NilBytesColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

// Type 返回列类型
func (n *NilBytesColumnValue) Type() ColumnType {
	return TypeBytes
}

// Clone 克隆空值字节流列值
func (n *NilBytesColumnValue) Clone() ColumnValue {
	return NewNilBytesColumnValue()
}

// BytesColumnValue 字节流列值
type BytesColumnValue struct {
	notNilColumnValue
	TimeEncoder //时间编码器

	val []byte //字节流值
}

// NewBytesColumnValue 从字节流v 生成字节流列值,做拷贝
func NewBytesColumnValue(v []byte) ColumnValue {
	new := make([]byte, len(v))
	copy(new, v)
	return NewBytesColumnValueNoCopy(new)
}

// NewBytesColumnValueNoCopy 从字节流v 生成字节流列值,不做拷贝
func NewBytesColumnValueNoCopy(v []byte) ColumnValue {
	return NewBytesColumnValueWithEncoderNoCopy(v, NewStringTimeEncoder(DefaultTimeFormat))
}

// NewBytesColumnValueWithEncoder 从字节流v 和时间编码器e 生成字节流列值,做拷贝
func NewBytesColumnValueWithEncoder(v []byte, e TimeEncoder) ColumnValue {
	new := make([]byte, len(v))
	copy(new, v)
	return NewBytesColumnValueWithEncoderNoCopy(new, NewStringTimeEncoder(DefaultTimeFormat))
}

// NewBytesColumnValueWithEncoderNoCopy 从字节流v 和时间编码器e,不做拷贝
func NewBytesColumnValueWithEncoderNoCopy(v []byte, e TimeEncoder) ColumnValue {
	return &BytesColumnValue{
		val:         v,
		TimeEncoder: e,
	}
}

// Type 返回列类型
func (b *BytesColumnValue) Type() ColumnType {
	return TypeBytes
}

// AsBool 1, t, T, TRUE, true, True转化为true
// 0, f, F, FALSE, false, False转化为false，如果不是上述情况会报错
func (b *BytesColumnValue) AsBool() (bool, error) {
	v, err := strconv.ParseBool(b.String())
	if err != nil {
		return false, NewTransformErrorFormColumnTypes(b.Type(), TypeBool, fmt.Errorf("err: %v val: %v", err, b.String()))
	}
	return v, nil
}

// AsBigInt 转化为整数，实数型以及科学性计数法字符串会被取整，不是数值型的会报错
// 如123.67转化为123 123.12转化为123
func (b *BytesColumnValue) AsBigInt() (BigIntNumber, error) {
	v, err := NewDecimalColumnValueFromString(b.String())
	if err != nil {
		return nil, NewTransformErrorFormColumnTypes(b.Type(), TypeBigInt, fmt.Errorf("err: %v, val: %v ", err, b.String()))
	}
	return v.AsBigInt()
}

// AsDecimal 转化为高精度师叔，实数型以及科学性计数法字符串能够转化，不是数值型的会报错
func (b *BytesColumnValue) AsDecimal() (DecimalNumber, error) {
	v, err := NewDecimalColumnValueFromString(b.String())
	if err != nil {
		return nil, NewTransformErrorFormColumnTypes(b.Type(), TypeDecimal, fmt.Errorf("err: %v, val: %v ", err, b.String()))
	}
	return v.AsDecimal()
}

// AsString 转化为字符串
func (b *BytesColumnValue) AsString() (string, error) {
	return b.String(), nil
}

// AsBytes 转化成字节流
func (b *BytesColumnValue) AsBytes() ([]byte, error) {
	v := make([]byte, len(b.val))
	copy(v, b.val)
	return v, nil
}

// AsTime 根据时间编码器转化成时间，不符合时间编码器格式会报错
func (b *BytesColumnValue) AsTime() (t time.Time, err error) {
	t, err = b.TimeEncode(b.String())
	if err != nil {
		return time.Time{}, NewTransformErrorFormColumnTypes(b.Type(), TypeTime, fmt.Errorf(" val: %v", b.String()))
	}
	return
}

func (b *BytesColumnValue) String() string {
	return string(b.val)
}

// Clone 克隆字节流列值
func (b *BytesColumnValue) Clone() ColumnValue {
	v := make([]byte, len(b.val))
	copy(v, b.val)
	return NewBytesColumnValue(v)
}

// Cmp  返回1代表大于， 0代表相等， -1代表小于
func (b *BytesColumnValue) Cmp(right ColumnValue) (int, error) {
	rightValue, err := right.AsBytes()
	if err != nil {
		return 0, err
	}

	if string(b.val) > string(rightValue) {
		return 1, nil
	}
	if string(b.val) == string(rightValue) {
		return 0, nil
	}
	return -1, nil
}
