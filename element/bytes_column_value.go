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

// NilBytesColumnValue - Null byte stream column value
type NilBytesColumnValue struct {
	*nilColumnValue
}

// NewNilBytesColumnValue - Create a null byte stream column value
func NewNilBytesColumnValue() ColumnValue {
	return &NilBytesColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

// Type - Return the column type
func (n *NilBytesColumnValue) Type() ColumnType {
	return TypeBytes
}

// Clone - Clone the null byte stream column value
func (n *NilBytesColumnValue) Clone() ColumnValue {
	return NewNilBytesColumnValue()
}

// BytesColumnValue - Byte stream column value
type BytesColumnValue struct {
	notNilColumnValue
	TimeEncoder // Time Encoder

	val []byte // Byte Stream Value
}

// NewBytesColumnValue - Generate a byte stream column value from byte stream v, making a copy
func NewBytesColumnValue(v []byte) ColumnValue {
	new := make([]byte, len(v))
	copy(new, v)
	return NewBytesColumnValueNoCopy(new)
}

// NewBytesColumnValueNoCopy - Generate a byte stream column value from byte stream v, without making a copy
func NewBytesColumnValueNoCopy(v []byte) ColumnValue {
	return NewBytesColumnValueWithEncoderNoCopy(v, NewStringTimeEncoder(DefaultTimeFormat))
}

// NewBytesColumnValueWithEncoder - Generate a byte stream column value from byte stream v and time encoder e, making a copy
func NewBytesColumnValueWithEncoder(v []byte, e TimeEncoder) ColumnValue {
	new := make([]byte, len(v))
	copy(new, v)
	return NewBytesColumnValueWithEncoderNoCopy(new, NewStringTimeEncoder(DefaultTimeFormat))
}

// NewBytesColumnValueWithEncoderNoCopy - Generate a byte stream column value from byte stream v and time encoder e, without making a copy
func NewBytesColumnValueWithEncoderNoCopy(v []byte, e TimeEncoder) ColumnValue {
	return &BytesColumnValue{
		val:         v,
		TimeEncoder: e,
	}
}

// Type - Return the column type
func (b *BytesColumnValue) Type() ColumnType {
	return TypeBytes
}

// AsBool - Convert 1, t, T, TRUE, true, True to true
// Convert 0, f, F, FALSE, false, False to false. If none of the above, an error is reported.
func (b *BytesColumnValue) AsBool() (bool, error) {
	v, err := strconv.ParseBool(b.String())
	if err != nil {
		return false, NewTransformErrorFormColumnTypes(b.Type(), TypeBool, fmt.Errorf("err: %v val: %v", err, b.String()))
	}
	return v, nil
}

// AsBigInt - Convert to an integer. Real numbers and scientific notation strings will be rounded. Non-numeric values will report an error.
// E.g., 123.67 is converted to 123, and 123.12 is converted to 123.
func (b *BytesColumnValue) AsBigInt() (BigIntNumber, error) {
	v, err := NewDecimalColumnValueFromString(b.String())
	if err != nil {
		return nil, NewTransformErrorFormColumnTypes(b.Type(), TypeBigInt, fmt.Errorf("err: %v, val: %v ", err, b.String()))
	}
	return v.AsBigInt()
}

// AsDecimal - Convert to a high-precision decimal. Real numbers and scientific notation strings can be converted. Non-numeric values will report an error.
func (b *BytesColumnValue) AsDecimal() (DecimalNumber, error) {
	v, err := NewDecimalColumnValueFromString(b.String())
	if err != nil {
		return nil, NewTransformErrorFormColumnTypes(b.Type(), TypeDecimal, fmt.Errorf("err: %v, val: %v ", err, b.String()))
	}
	return v.AsDecimal()
}

// AsString - Convert to a string
func (b *BytesColumnValue) AsString() (string, error) {
	return b.String(), nil
}

// AsBytes - Convert to a byte stream
func (b *BytesColumnValue) AsBytes() ([]byte, error) {
	v := make([]byte, len(b.val))
	copy(v, b.val)
	return v, nil
}

// AsTime - Convert to time based on the time encoder. If it does not match the time encoder format, an error is reported.
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

// Clone - Clone the byte stream column value
func (b *BytesColumnValue) Clone() ColumnValue {
	v := make([]byte, len(b.val))
	copy(v, b.val)
	return NewBytesColumnValue(v)
}

// Cmp - Return 1 for greater than, 0 for equal, and -1 for less than
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
