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

// NilStringColumnValue - Null value for a string column
type NilStringColumnValue struct {
	*nilColumnValue
}

// NewNilStringColumnValue - Create a null value for a string column
func NewNilStringColumnValue() ColumnValue {
	return &NilStringColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

// Type - Column type
func (n *NilStringColumnValue) Type() ColumnType {
	return TypeString
}

// Clone - Clone the null string value
func (n *NilStringColumnValue) Clone() ColumnValue {
	return NewNilStringColumnValue()
}

// StringColumnValue - Value for a string column. Note: Decimal 123.0 (val:1230, exp:-1) is not equivalent to 123 (val:123, exp:0)
type StringColumnValue struct {
	notNilColumnValue
	TimeEncoder
	val string
}

// NewStringColumnValue - Create a string column value based on the string s
func NewStringColumnValue(s string) ColumnValue {
	return NewStringColumnValueWithEncoder(s, NewStringTimeEncoder(DefaultTimeFormat))
}

// NewStringColumnValueWithEncoder - Create a string column value based on the string s and time encoder e
func NewStringColumnValueWithEncoder(s string, e TimeEncoder) ColumnValue {
	return &StringColumnValue{
		TimeEncoder: e,
		val:         s,
	}
}

// Type - Column type
func (s *StringColumnValue) Type() ColumnType {
	return TypeString
}

// AsBool - Convert 1, t, T, TRUE, true, True to true
// Convert 0, f, F, FALSE, false, False to false. If none of the above, an error is thrown.
func (s *StringColumnValue) AsBool() (v bool, err error) {
	v, err = strconv.ParseBool(s.val)
	if err != nil {
		return false, NewTransformErrorFormColumnTypes(s.Type(), TypeBool, fmt.Errorf("err: %v, val: %v ", err, s.val))
	}
	return
}

// AsBigInt - Convert to a big integer. Floating-point numbers and scientific notation strings will be rounded. Non-numeric values will throw an error.
// E.g., 123.67 is converted to 123, and 123.12 is converted to 123.
func (s *StringColumnValue) AsBigInt() (BigIntNumber, error) {
	v, err := NewDecimalColumnValueFromString(s.val)
	if err != nil {
		return nil, NewTransformErrorFormColumnTypes(s.Type(), TypeBigInt, fmt.Errorf("err: %v, val: %v ", err, s.val))
	}
	return v.AsBigInt()
}

// AsDecimal - Convert to a decimal. Floating-point numbers and scientific notation strings can be converted. Non-numeric values will throw an error.
func (s *StringColumnValue) AsDecimal() (DecimalNumber, error) {
	v, err := NewDecimalColumnValueFromString(s.val)
	if err != nil {
		return nil, NewTransformErrorFormColumnTypes(s.Type(), TypeDecimal,
			fmt.Errorf("err: %v, val: %v ", err, s.val))
	}
	return v.AsDecimal()
}

// AsString - Convert to a string
func (s *StringColumnValue) AsString() (string, error) {
	return s.val, nil
}

// AsBytes - Convert to a byte stream
func (s *StringColumnValue) AsBytes() ([]byte, error) {
	return []byte(s.val), nil
}

// AsTime - Convert to a time based on the time encoder. If the format does not match the time encoder, an error is thrown.
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

// Clone - Clone the string column value
func (s *StringColumnValue) Clone() ColumnValue {
	return NewStringColumnValue(s.val)
}

// Cmp - Return 1 for greater than, 0 for equal, -1 for less than
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
