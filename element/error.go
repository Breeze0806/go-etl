package element

import (
	"errors"
	"fmt"
)

//错误
var (
	ErrPrecisionNotEnough     = errors.New("precision is not enough")      //精度不足错误
	ErrColumnExist            = errors.New("column exist")                 //列存在错误
	ErrColumnNotExist         = errors.New("column does not exist")        //列不存在错误
	ErrNilValue               = errors.New("column value is nil")          //空值错误
	ErrIndexOutOfRange        = errors.New("column index is out of range") //索引值超出范围
	ErrValueNotInt64          = errors.New("value is not int64")           //不是64错误
	ErrValueInfinity          = errors.New("value is infinity")            //无穷大实数错误
	ErrNotColumnValueClonable = errors.New("columnValue is not clonable")  //不是可克隆列值
)

//TransformError 转化错误
type TransformError struct {
	err error
	msg string
}

//NewTransformError 根据消息msg和错误err生成转化错误
func NewTransformError(msg string, err error) *TransformError {
	for uerr := err; uerr != nil; uerr = errors.Unwrap(err) {
		err = uerr
	}
	return &TransformError{
		msg: msg,
		err: err,
	}
}

//NewTransformErrorFormColumnTypes 从one类型到other类型转化错误err生成转化错误
func NewTransformErrorFormColumnTypes(one, other ColumnType, err error) *TransformError {
	return NewTransformError(fmt.Sprintf("%s transform to %s", one, other), err)
}

//NewTransformErrorFormString 从one到other转化错误err生成转化错误
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

//SetError 设置错误
type SetError struct {
	err error
	msg string
}

//NewSetError 通过值i设置成累心other类型的错误err生成设置错误
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
