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

// NilBigIntColumnValue - Null value for a big integer column
type NilBigIntColumnValue struct {
	*nilColumnValue
}

// NewNilBigIntColumnValue - Create a null value for a big integer column
func NewNilBigIntColumnValue() ColumnValue {
	return &NilBigIntColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

// Type - Return the type of the column
func (n *NilBigIntColumnValue) Type() ColumnType {
	return TypeBigInt
}

// Clone - Clone the null value for a big integer column
func (n *NilBigIntColumnValue) Clone() ColumnValue {
	return NewNilBigIntColumnValue()
}

// BigIntColumnValue - Value for a big integer column
type BigIntColumnValue struct {
	notNilColumnValue

	val BigIntNumber
}

// NewBigIntColumnValueFromInt64 - Create a big integer column value from an int64 v
func NewBigIntColumnValueFromInt64(v int64) ColumnValue {
	return &BigIntColumnValue{
		val: _DefaultNumberConverter.ConvertBigIntFromInt(v),
	}
}

// NewBigIntColumnValue - Create a big integer column value from a big.Int v
func NewBigIntColumnValue(v *big.Int) ColumnValue {
	return &BigIntColumnValue{
		val: &BigInt{
			value: new(big.Int).Set(v),
		},
	}
}

// NewBigIntColumnValueFromString - Create a big integer column value from a string v
// If the string v is not an integer, an error is returned
func NewBigIntColumnValueFromString(v string) (ColumnValue, error) {
	num, err := _DefaultNumberConverter.ConvertBigInt(v)
	if err != nil {
		return nil, NewSetError(v, TypeBigInt, fmt.Errorf("string %v is not valid int", v))
	}
	return &BigIntColumnValue{
		val: num,
	}, nil
}

// Type - Return the type of the column
func (b *BigIntColumnValue) Type() ColumnType {
	return TypeBigInt
}

// AsBool - Convert to a boolean value, where non-zero becomes true and zero becomes false
func (b *BigIntColumnValue) AsBool() (bool, error) {
	return b.val.Bool()
}

// AsBigInt - Convert to a big integer
func (b *BigIntColumnValue) AsBigInt() (BigIntNumber, error) {
	return b.val, nil
}

// AsDecimal - Convert to a high-precision decimal number
func (b *BigIntColumnValue) AsDecimal() (DecimalNumber, error) {
	return b.val.Decimal(), nil
}

// AsString - Convert to a string, e.g., 1234556790 becomes 1234556790
func (b *BigIntColumnValue) AsString() (string, error) {
	return b.val.String(), nil
}

// AsBytes - Convert to a byte stream, e.g., 1234556790 becomes [49, 50, 51, 52, 53, 54, 55, 56, 57, 48]
func (b *BigIntColumnValue) AsBytes() ([]byte, error) {
	return []byte(b.val.String()), nil
}

// AsTime - Currently, integers cannot be converted to time
func (b *BigIntColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, NewTransformErrorFormColumnTypes(b.Type(), TypeTime, fmt.Errorf(" val: %v", b.String()))
}

func (b *BigIntColumnValue) String() string {
	return b.val.String()
}

// Clone - Clone the big integer column value
func (b *BigIntColumnValue) Clone() ColumnValue {
	return &BigIntColumnValue{
		val: b.val.CloneBigInt(),
	}
}

// Cmp - Return 1 for greater than, 0 for equal, and -1 for less than
func (b *BigIntColumnValue) Cmp(right ColumnValue) (int, error) {
	rightValue, err := right.AsBigInt()
	if err != nil {
		return 0, err
	}
	return b.val.AsBigInt().Cmp(rightValue.AsBigInt()), nil
}
