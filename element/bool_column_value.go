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

	"github.com/shopspring/decimal"
)

// NilBoolColumnValue 空值布尔列值
type NilBoolColumnValue struct {
	*nilColumnValue
}

// NewNilBoolColumnValue 生成空值布尔列值
func NewNilBoolColumnValue() ColumnValue {
	return &NilBoolColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

// Type 返回列类型
func (n *NilBoolColumnValue) Type() ColumnType {
	return TypeBool
}

// Clone 克隆空值布尔列值
func (n *NilBoolColumnValue) Clone() ColumnValue {
	return NewNilBoolColumnValue()
}

// BoolColumnValue 布尔列值
type BoolColumnValue struct {
	notNilColumnValue

	val bool //布尔值
}

// NewBoolColumnValue 从布尔值v生成布尔列值
func NewBoolColumnValue(v bool) ColumnValue {
	return &BoolColumnValue{
		val: v,
	}
}

// Type 返回列类型
func (b *BoolColumnValue) Type() ColumnType {
	return TypeBool
}

// AsBool 转化成布尔值
func (b *BoolColumnValue) AsBool() (bool, error) {
	return b.val, nil
}

// AsBigInt 转化成整数，true转化为1，false转化为0
func (b *BoolColumnValue) AsBigInt() (BigIntNumber, error) {
	if b.val {
		return NewBigIntColumnValue(big.NewInt(1)).AsBigInt()
	}
	return NewBigIntColumnValue(big.NewInt(0)).AsBigInt()
}

// AsDecimal 转化成高精度实数，true转化为1.0，false转化为0.0
func (b *BoolColumnValue) AsDecimal() (DecimalNumber, error) {
	if b.val {
		return NewDecimalColumnValue(decimal.New(1, 0)).AsDecimal()
	}
	return NewDecimalColumnValue(decimal.New(0, 1)).AsDecimal()
}

// AsString 转化成字符串，true转化为"true"，false转化为"false"
func (b *BoolColumnValue) AsString() (string, error) {
	if b.val {
		return b.String(), nil
	}
	return b.String(), nil
}

// AsBytes 转化成字节流，true转化为"true"，false转化为"false"
func (b *BoolColumnValue) AsBytes() ([]byte, error) {
	if b.val {
		return []byte(b.String()), nil
	}
	return []byte(b.String()), nil
}

// AsTime 目前布尔无法转化成时间
func (b *BoolColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, NewTransformErrorFormColumnTypes(b.Type(), TypeTime, fmt.Errorf(" val: %v", b.String()))
}

func (b *BoolColumnValue) String() string {
	if b.val {
		return "true"
	}
	return "false"
}

// Clone 克隆布尔列值
func (b *BoolColumnValue) Clone() ColumnValue {
	return NewBoolColumnValue(b.val)
}

// Cmp  返回1代表大于， 0代表相等， -1代表小于
func (b *BoolColumnValue) Cmp(right ColumnValue) (int, error) {
	rightValue, err := right.AsBool()
	if err != nil {
		return 0, err
	}

	if b.val == rightValue {
		return 0, nil
	}

	if b.val && !rightValue {
		return 1, nil
	}
	return -1, nil
}
