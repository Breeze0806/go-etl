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
	"time"

	"github.com/shopspring/decimal"
)

// NilDecimalColumnValue represents a null value for a high-precision decimal column.
type NilDecimalColumnValue struct {
	*nilColumnValue
}

// NewNilDecimalColumnValue creates a null value for a high-precision decimal column.
func NewNilDecimalColumnValue() ColumnValue {
	return &NilDecimalColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

// Type represents the type of the column.
func (n *NilDecimalColumnValue) Type() ColumnType {
	return TypeDecimal
}

// Clone creates a copy (clone) of a high-precision decimal column value.
func (n *NilDecimalColumnValue) Clone() ColumnValue {
	return NewNilDecimalColumnValue()
}

// DecimalColumnValue represents a high-precision decimal column value.
type DecimalColumnValue struct {
	notNilColumnValue

	val DecimalNumber // High-precision decimal
}

// NewDecimalColumnValueFromFloat creates a high-precision decimal column value from a float64 value.
func NewDecimalColumnValueFromFloat(f float64) ColumnValue {
	return &DecimalColumnValue{
		val: _DefaultNumberConverter.ConvertDecimalFromFloat(f),
	}
}

// NewDecimalColumnValue creates a high-precision decimal column value from a high-precision decimal value.
func NewDecimalColumnValue(d decimal.Decimal) ColumnValue {
	return &DecimalColumnValue{
		val: &Decimal{
			value: d,
		},
	}
}

// NewDecimalColumnValueFromString creates a high-precision decimal column value from a string.
// If the string is not a numeric value or in scientific notation, an error will be reported.
func NewDecimalColumnValueFromString(s string) (ColumnValue, error) {
	num, err := _DefaultNumberConverter.ConvertDecimal(s)
	if err != nil {
		return nil, NewSetError(s, TypeDecimal, fmt.Errorf("string %v is not valid decimal", s))
	}
	return &DecimalColumnValue{
		val: num,
	}, nil
}

// Type represents the type of the column.
func (d *DecimalColumnValue) Type() ColumnType {
	return TypeDecimal
}

// AsBool converts non-zero values to true and zero values to false.
func (d *DecimalColumnValue) AsBool() (bool, error) {
	return d.val.Bool()
}

// AsBigInt rounds down the high-precision decimal value, e.g., 123.67 becomes 123, and 123.12 becomes 123.
func (d *DecimalColumnValue) AsBigInt() (BigIntNumber, error) {
	return d.val.BigInt(), nil
}

// AsDecimal converts to a high-precision decimal.
func (d *DecimalColumnValue) AsDecimal() (DecimalNumber, error) {
	return d.val, nil
}

// AsString converts to a string, e.g., 10.123 becomes 10.123.
func (d *DecimalColumnValue) AsString() (string, error) {
	return d.val.String(), nil
}

// AsBytes converts to a byte stream, e.g., 10.123 becomes the byte representation of 10.123.
func (d *DecimalColumnValue) AsBytes() ([]byte, error) {
	return []byte(d.val.String()), nil
}

// AsTime currently cannot convert to a time value.
func (d *DecimalColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, NewTransformErrorFormColumnTypes(d.Type(), TypeTime, fmt.Errorf(" val: %v", d.String()))
}

func (d *DecimalColumnValue) String() string {
	return d.val.String()
}

// Clone creates a copy (clone) of a high-precision decimal column value.
func (d *DecimalColumnValue) Clone() ColumnValue {
	return &DecimalColumnValue{
		val: d.val,
	}
}

// Cmp returns 1 for greater than, 0 for equal, and -1 for less than.
func (d *DecimalColumnValue) Cmp(right ColumnValue) (int, error) {
	rightValue, err := right.AsDecimal()
	if err != nil {
		return 0, err
	}

	return d.val.AsDecimal().Cmp(rightValue.AsDecimal()), nil
}
