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
	"errors"
	"fmt"
)

// Error
var (
	ErrColumnExist              = errors.New("column exist")                  // Column exists error
	ErrColumnNotExist           = errors.New("column does not exist")         // Column does not exist error
	ErrNilValue                 = errors.New("column value is nil")           // Null value error
	ErrIndexOutOfRange          = errors.New("column index is out of range")  // Index value out of range
	ErrValueNotInt64            = errors.New("value is not int64")            // Not an int64 error
	ErrValueInfinity            = errors.New("value is infinity")             // Infinite real number error
	ErrNotColumnValueClonable   = errors.New("columnValue is not clonable")   // Not a clonable column value
	ErrNotColumnValueComparable = errors.New("columnValue is not comparable") // Not a comparable column value
	ErrColumnNameNotEqual       = errors.New("column name is not equal")      // Column names differ
)

// TransformError: Conversion error
type TransformError struct {
	err error
	msg string
}

// NewTransformError: Creates a conversion error based on the message msg and error err
func NewTransformError(msg string, err error) *TransformError {
	for uerr := err; uerr != nil; uerr = errors.Unwrap(err) {
		err = uerr
	}
	return &TransformError{
		msg: msg,
		err: err,
	}
}

// NewTransformErrorFromColumnTypes: Generates a conversion error from the error err when converting from type one to type other
func NewTransformErrorFormColumnTypes(one, other ColumnType, err error) *TransformError {
	return NewTransformError(fmt.Sprintf("%s transform to %s", one, other), err)
}

// NewTransformErrorFromString: Generates a conversion error from the error err when converting from one to other
func NewTransformErrorFormString(one, other string, err error) *TransformError {
	return NewTransformError(fmt.Sprintf("%s transform to %s", one, other), err)
}

func (e *TransformError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s error: %v", e.msg, e.err)
	}
	return fmt.Sprintf("%s", e.msg)
}

func (e *TransformError) Unwrap() error {
	return e.err
}

// SetError: Sets an error
type SetError struct {
	err error
	msg string
}

// NewSetError: Generates a setting error by setting the value i to the specified other type with the error err
func NewSetError(i interface{}, other ColumnType, err error) *SetError {
	for uerr := err; uerr != nil; uerr = errors.Unwrap(err) {
		err = uerr
	}
	return &SetError{
		msg: fmt.Sprintf("%T set to %s", i, other),
		err: err,
	}
}

func (e *SetError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s error: %v", e.msg, e.err)
	}
	return fmt.Sprintf("%s", e.msg)
}

func (e *SetError) Unwrap() error {
	return e.err
}
