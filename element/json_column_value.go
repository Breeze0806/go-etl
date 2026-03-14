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

// NilJsonColumnValue represents a nil JSON column value
type NilJsonColumnValue struct {
	*nilColumnValue
}

// NewNilJsonColumnValue creates a nil JSON column value
func NewNilJsonColumnValue() ColumnValue {
	return &NilJsonColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

// Type returns the type of the column
func (n *NilJsonColumnValue) Type() ColumnType {
	return TypeJSON
}

// Clone clones a nil JSON column value
func (n *NilJsonColumnValue) Clone() ColumnValue {
	return NewNilJsonColumnValue()
}

// JsonColumnValue represents a JSON column value
type JsonColumnValue struct {
	notNilColumnValue

	val JSON // JSON Value
}

// NewJsonColumnValueFromString creates a JSON column value from a string
func NewJsonColumnValueFromString(s string) (ColumnValue, error) {
	json, err := _DefaultJSONConverter.ConvertFromString(s)
	if err != nil {
		return nil, err
	}
	return &JsonColumnValue{val: json}, nil
}

// NewJsonColumnValueFromBytes creates a JSON column value from bytes
func NewJsonColumnValueFromBytes(b []byte) (ColumnValue, error) {
	json, err := _DefaultJSONConverter.ConvertFromBytes(b)
	if err != nil {
		return nil, err
	}
	return &JsonColumnValue{val: json}, nil
}

// Type returns the type of the column
func (j *JsonColumnValue) Type() ColumnType {
	return TypeJSON
}

// AsBool cannot convert to a boolean
func (j *JsonColumnValue) AsBool() (bool, error) {
	return false, NewTransformErrorFormColumnTypes(j.Type(), TypeBool, fmt.Errorf(" val: %v", j.String()))
}

// AsBigInt cannot convert to a big integer
func (j *JsonColumnValue) AsBigInt() (BigIntNumber, error) {
	return nil, NewTransformErrorFormColumnTypes(j.Type(), TypeBigInt, fmt.Errorf(" val: %v", j.String()))
}

// AsDecimal cannot convert to a high-precision decimal number
func (j *JsonColumnValue) AsDecimal() (DecimalNumber, error) {
	return nil, NewTransformErrorFormColumnTypes(j.Type(), TypeDecimal, fmt.Errorf(" val: %v", j.String()))
}

// AsString converts it to a string
func (j *JsonColumnValue) AsString() (string, error) {
	return j.val.ToString(), nil
}

// AsBytes converts it to a byte stream
func (j *JsonColumnValue) AsBytes() ([]byte, error) {
	return j.val.ToBytes(), nil
}

// AsTime cannot convert to a time value
func (j *JsonColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, NewTransformErrorFormColumnTypes(j.Type(), TypeTime, fmt.Errorf(" val: %v", j.String()))
}

// AsJSON converts to a JSON value
func (j *JsonColumnValue) AsJSON() (JSON, error) {
	return j.val, nil
}

func (j *JsonColumnValue) String() string {
	return j.val.ToString()
}

// Clone clones a JSON column value
func (j *JsonColumnValue) Clone() ColumnValue {
	return &JsonColumnValue{val: j.val.Clone()}
}

// Cmp Returns 1 for greater than, 0 for equal, and -1 for less than
func (j *JsonColumnValue) Cmp(right ColumnValue) (int, error) {
	rightValue, err := right.AsString()
	if err != nil {
		return 0, err
	}

	if j.val.ToString() > rightValue {
		return 1, nil
	}
	if j.val.ToString() == rightValue {
		return 0, nil
	}
	return -1, nil
}
