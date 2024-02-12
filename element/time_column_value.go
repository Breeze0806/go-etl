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
)

// NilTimeColumnValue represents an empty time column value
type NilTimeColumnValue struct {
	*nilColumnValue
}

// NewNilTimeColumnValue creates an empty time column value
func NewNilTimeColumnValue() ColumnValue {
	return &NilTimeColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

// Type returns the type of the column
func (n *NilTimeColumnValue) Type() ColumnType {
	return TypeTime
}

// Clone clones an empty time column value
func (n *NilTimeColumnValue) Clone() ColumnValue {
	return NewNilTimeColumnValue()
}

// TimeColumnValue represents a time column value
type TimeColumnValue struct {
	notNilColumnValue
	TimeDecoder // Time decoder

	val time.Time
}

// NewTimeColumnValue creates a time column value from the time t
func NewTimeColumnValue(t time.Time) ColumnValue {
	return NewTimeColumnValueWithDecoder(t, NewStringTimeDecoder(DefaultTimeFormat))
}

// NewTimeColumnValueWithDecoder creates a time column value from the time t and the time decoder
func NewTimeColumnValueWithDecoder(t time.Time, d TimeDecoder) ColumnValue {
	return &TimeColumnValue{
		TimeDecoder: d,
		val:         t,
	}
}

// Type returns the type of the column
func (t *TimeColumnValue) Type() ColumnType {
	return TypeTime
}

// AsBool: Cannot convert to a boolean value
func (t *TimeColumnValue) AsBool() (bool, error) {
	return false, NewTransformErrorFormColumnTypes(t.Type(), TypeBool, fmt.Errorf("val: %v", t.String()))
}

// AsBigInt: Cannot convert to a big integer
func (t *TimeColumnValue) AsBigInt() (BigIntNumber, error) {
	return nil, NewTransformErrorFormColumnTypes(t.Type(), TypeBigInt, fmt.Errorf("val: %v", t.String()))
}

// AsDecimal: Cannot convert to a high-precision decimal number
func (t *TimeColumnValue) AsDecimal() (DecimalNumber, error) {
	return nil, NewTransformErrorFormColumnTypes(t.Type(), TypeDecimal, fmt.Errorf("val: %v", t.String()))
}

// AsString: Converts to a string
func (t *TimeColumnValue) AsString() (s string, err error) {
	var i interface{}
	i, err = t.TimeDecode(t.val)
	if err != nil {
		return "", NewTransformErrorFormColumnTypes(t.Type(), TypeString, fmt.Errorf("val: %v", t.String()))
	}
	return i.(string), nil
}

// AsBytes: Converts to a byte stream
func (t *TimeColumnValue) AsBytes() (b []byte, err error) {
	var i interface{}
	i, err = t.TimeDecode(t.val)
	if err != nil {
		return nil, NewTransformErrorFormColumnTypes(t.Type(), TypeString, fmt.Errorf("val: %v", t.String()))
	}
	return []byte(i.(string)), nil
}

// AsTime: Converts to a time value
func (t *TimeColumnValue) AsTime() (time.Time, error) {
	return t.val, nil
}

func (t *TimeColumnValue) String() string {
	return t.val.Format(DefaultTimeFormat)
}

// Clone clones a time column value
func (t *TimeColumnValue) Clone() ColumnValue {
	return &TimeColumnValue{
		val: t.val,
	}
}

// Cmp: Returns 1 for greater than, 0 for equal, and -1 for less than
func (t *TimeColumnValue) Cmp(right ColumnValue) (int, error) {
	rightValue, err := right.AsTime()
	if err != nil {
		return 0, err
	}

	if t.val.After(rightValue) {
		return 1, nil
	}
	if t.val.Before(rightValue) {
		return -1, nil
	}
	return 0, nil
}
