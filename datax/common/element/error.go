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
)

type TransformError struct {
	err error
	msg string
}

func NewTransformTypeError(one, other ColumnType) *TransformError {
	return &TransformError{
		msg: fmt.Sprintf("%s can not Transform to %s", one, other),
	}
}

func NewTransformError(one, other ColumnType, err error) *TransformError {
	return &TransformError{
		msg: fmt.Sprintf("%s Transform to %s", one, other),
		err: err,
	}
}

func (e *TransformError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s error: %v", e.msg, e.err)
	}
	return fmt.Sprintf("%s", e.msg)
}

type SetError struct {
	err error
	msg string
}

func NewSetError(i interface{}, other ColumnType, err error) *SetError {
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
