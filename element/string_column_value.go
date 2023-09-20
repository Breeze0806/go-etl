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

// NilStringColumnValue 空值字符串列值
type NilStringColumnValue struct {
	*nilColumnValue
}

// NewNilStringColumnValue 创建空值字符串列值
func NewNilStringColumnValue() ColumnValue {
	return &NilStringColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

// Type 列类型
func (n *NilStringColumnValue) Type() ColumnType {
	return TypeString
}

// Clone 克隆空值字符串
func (n *NilStringColumnValue) Clone() ColumnValue {
	return NewNilStringColumnValue()
}

// StringColumnValue 字符串列名 注意：Decimal 123.0（val:1230,exp:-1）和123（val:123,exp:0）不一致
type StringColumnValue struct {
	notNilColumnValue
	TimeEncoder
	val string
}

// NewStringColumnValue 根据字符串s 生成字符串列值
func NewStringColumnValue(s string) ColumnValue {
	return NewStringColumnValueWithEncoder(s, NewStringTimeEncoder(DefaultTimeFormat))
}

// NewStringColumnValueWithEncoder 根据字符串s 时间编码器e生成字符串列值
func NewStringColumnValueWithEncoder(s string, e TimeEncoder) ColumnValue {
	return &StringColumnValue{
		TimeEncoder: e,
		val:         s,
	}
}

// Type 列类型
func (s *StringColumnValue) Type() ColumnType {
	return TypeString
}

// AsBool 1, t, T, TRUE, true, True转化为true
// 0, f, F, FALSE, false, False转化为false，如果不是上述情况会报错
func (s *StringColumnValue) AsBool() (v bool, err error) {
	v, err = strconv.ParseBool(s.val)
	if err != nil {
		return false, NewTransformErrorFormColumnTypes(s.Type(), TypeBool, fmt.Errorf("err: %v, val: %v ", err, s.val))
	}
	return
}

// AsBigInt 转化为整数，实数型以及科学性计数法字符串会被取整，不是数值型的会报错
// 如123.67转化为123 123.12转化为123
func (s *StringColumnValue) AsBigInt() (BigIntNumber, error) {
	v, err := NewDecimalColumnValueFromString(s.val)
	if err != nil {
		return nil, NewTransformErrorFormColumnTypes(s.Type(), TypeBigInt, fmt.Errorf("err: %v, val: %v ", err, s.val))
	}
	return v.AsBigInt()
}

// AsDecimal 转化为整数，实数型以及科学性计数法字符串能够转化，不是数值型的会报错
func (s *StringColumnValue) AsDecimal() (DecimalNumber, error) {
	v, err := NewDecimalColumnValueFromString(s.val)
	if err != nil {
		return nil, NewTransformErrorFormColumnTypes(s.Type(), TypeDecimal,
			fmt.Errorf("err: %v, val: %v ", err, s.val))
	}
	return v.AsDecimal()
}

// AsString 转化为字符串
func (s *StringColumnValue) AsString() (string, error) {
	return s.val, nil
}

// AsBytes 转化成字节流
func (s *StringColumnValue) AsBytes() ([]byte, error) {
	return []byte(s.val), nil
}

// AsTime 根据时间编码器转化成时间，不符合时间编码器格式会报错
func (s *StringColumnValue) AsTime() (t time.Time, err error) {
	t, err = s.TimeEncode(s.val)
	if err != nil {
		return time.Time{}, NewTransformErrorFormColumnTypes(s.Type(), TypeTime, fmt.Errorf("err: %v val: %v", err, s.val))
	}
	return
}

func (s *StringColumnValue) String() string {
	return s.val
}

// Clone 克隆字符串列值
func (s *StringColumnValue) Clone() ColumnValue {
	return NewStringColumnValue(s.val)
}

// Cmp  返回1代表大于， 0代表相等， -1代表小于
func (s *StringColumnValue) Cmp(right ColumnValue) (int, error) {
	rightValue, err := right.AsString()
	if err != nil {
		return 0, err
	}

	if s.val > rightValue {
		return 1, nil
	}

	if s.val == rightValue {
		return 0, nil
	}

	return -1, nil
}
