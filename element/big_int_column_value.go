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
	"math/big"
	"time"
)

// NilBigIntColumnValue 空值整数列值
type NilBigIntColumnValue struct {
	*nilColumnValue
}

// NewNilBigIntColumnValue 创建空值整数列值
func NewNilBigIntColumnValue() ColumnValue {
	return &NilBigIntColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

// Type 返回列类型
func (n *NilBigIntColumnValue) Type() ColumnType {
	return TypeBigInt
}

// Clone 克隆空值整数列值
func (n *NilBigIntColumnValue) Clone() ColumnValue {
	return NewNilBigIntColumnValue()
}

// BigIntColumnValue 整数列值
type BigIntColumnValue struct {
	notNilColumnValue

	val BigIntNumber
}

// NewBigIntColumnValueFromInt64 从int64 v中获取整数列值
func NewBigIntColumnValueFromInt64(v int64) ColumnValue {
	return &BigIntColumnValue{
		val: _DefaultNumberConverter.ConvertBigIntFromInt(v),
	}
}

// NewBigIntColumnValue 从big.Int v中获取整数列值
func NewBigIntColumnValue(v *big.Int) ColumnValue {
	return &BigIntColumnValue{
		val: &BigInt{
			value: new(big.Int).Set(v),
		},
	}
}

// NewBigIntColumnValueFromString 从string v中获取整数列值
// 当string v不是整数时,返回错误
func NewBigIntColumnValueFromString(v string) (ColumnValue, error) {
	num, err := _DefaultNumberConverter.ConvertBigInt(v)
	if err != nil {
		return nil, NewSetError(v, TypeBigInt, fmt.Errorf("string %v is not valid int", v))
	}
	return &BigIntColumnValue{
		val: num,
	}, nil
}

// Type 返回列类型
func (b *BigIntColumnValue) Type() ColumnType {
	return TypeBigInt
}

// AsBool 转化成布尔值，不是0的转化为true,是0的转化成false
func (b *BigIntColumnValue) AsBool() (bool, error) {
	return b.val.Bool()
}

// AsBigInt 转化成整数
func (b *BigIntColumnValue) AsBigInt() (BigIntNumber, error) {
	return b.val, nil
}

// AsDecimal 转化成高精度实数
func (b *BigIntColumnValue) AsDecimal() (DecimalNumber, error) {
	return b.val.Decimal(), nil
}

// AsString 转化成字符串，如1234556790转化为1234556790
func (b *BigIntColumnValue) AsString() (string, error) {
	return b.val.String(), nil
}

// AsBytes 转化成字节流，如1234556790转化为1234556790
func (b *BigIntColumnValue) AsBytes() ([]byte, error) {
	return []byte(b.val.String()), nil
}

// AsTime 目前整数无法转化成时间
func (b *BigIntColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, NewTransformErrorFormColumnTypes(b.Type(), TypeTime, fmt.Errorf(" val: %v", b.String()))
}

func (b *BigIntColumnValue) String() string {
	return b.val.String()
}

// Clone 克隆整数列属性
func (b *BigIntColumnValue) Clone() ColumnValue {
	return &BigIntColumnValue{
		val: b.val.CloneBigInt(),
	}
}

// Cmp  返回1代表大于， 0代表相等， -1代表小于
func (b *BigIntColumnValue) Cmp(right ColumnValue) (int, error) {
	rightValue, err := right.AsBigInt()
	if err != nil {
		return 0, err
	}
	return b.val.AsBigInt().Cmp(rightValue.AsBigInt()), nil
}
