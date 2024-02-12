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

// NilBoolColumnValue represents an empty Boolean column value
type NilBoolColumnValue struct {
	*nilColumnValue
}

// NewNilBoolColumnValue creates an empty Boolean column value
func NewNilBoolColumnValue() ColumnValue {
	return &NilBoolColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

// Type returns the type of the column
func (n *NilBoolColumnValue) Type() ColumnType {
	return TypeBool
}

// Clone clones an empty Boolean column value
func (n *NilBoolColumnValue) Clone() ColumnValue {
	return NewNilBoolColumnValue()
}

// BoolColumnValue represents a Boolean column value
type BoolColumnValue struct {
	notNilColumnValue

	val bool // Boolean Value
}

// NewBoolColumnValue creates a Boolean column value from the boolean value v
func NewBoolColumnValue(v bool) ColumnValue {
	return &BoolColumnValue{
		val: v,
	}
}

// Type returns the type of the column
func (b *BoolColumnValue) Type() ColumnType {
	return TypeBool
}

// AsBool converts it to a boolean value
func (b *BoolColumnValue) AsBool() (bool, error) {
	return b.val, nil
}

// AsBigInt converts it to a big integer, where true becomes 1 and false becomes 0
func (b *BoolColumnValue) AsBigInt() (BigIntNumber, error) {
	if b.val {
		return NewBigIntColumnValue(big.NewInt(1)).AsBigInt()
	}
	return NewBigIntColumnValue(big.NewInt(0)).AsBigInt()
}

// AsDecimal converts it to a high-precision decimal number, where true becomes 1.0 and false becomes 0.0
func (b *BoolColumnValue) AsDecimal() (DecimalNumber, error) {
	if b.val {
		return NewDecimalColumnValue(decimal.New(1, 0)).AsDecimal()
	}
	return NewDecimalColumnValue(decimal.New(0, 1)).AsDecimal()
}

// AsString converts it to a string, where true becomes true and false becomes false
func (b *BoolColumnValue) AsString() (string, error) {
	if b.val {
		return b.String(), nil
	}
	return b.String(), nil
}

// AsBytes converts it to a byte stream, where true becomes true and false becomes false
func (b *BoolColumnValue) AsBytes() ([]byte, error) {
	if b.val {
		return []byte(b.String()), nil
	}
	return []byte(b.String()), nil
}

// AsTime: Currently, a Boolean cannot be converted to a time value
func (b *BoolColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, NewTransformErrorFormColumnTypes(b.Type(), TypeTime, fmt.Errorf(" val: %v", b.String()))
}

func (b *BoolColumnValue) String() string {
	if b.val {
		return "true"
	}
	return "false"
}

// Clone clones a Boolean column value
func (b *BoolColumnValue) Clone() ColumnValue {
	return NewBoolColumnValue(b.val)
}

// Cmp: Returns 1 for greater than, 0 for equal, and -1 for less than
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
