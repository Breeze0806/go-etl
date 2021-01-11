package element

import (
	"errors"
	"fmt"
)

var (
	ErrPrecisionNotEnough = errors.New("precision is not enough")
	ErrColumnExist        = errors.New("column exist")
	ErrColumnNotExist     = errors.New("column does not exist")
	ErrNilValue           = errors.New("column value is nil")
	ErrIndexOutOfRange    = errors.New("column index is out of range")
	ErrValueNotInt64      = errors.New("value is not int64")
	ErrValueInfinity      = errors.New("Value is infinity")
)

type TransformError struct {
	err error
	msg string
}

func NewTransformError(msg string, err error) *TransformError {
	for uerr := err; uerr != nil; uerr = errors.Unwrap(err) {
		err = uerr
	}
	return &TransformError{
		msg: msg,
		err: err,
	}
}

func NewTransformErrorFormColumnTypes(one, other ColumnType, err error) *TransformError {
	return NewTransformError(fmt.Sprintf("%s transform to %s", one, other), err)
}

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

type SetError struct {
	err error
	msg string
}

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
