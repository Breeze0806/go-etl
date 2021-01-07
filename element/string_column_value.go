package element

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type NilStringColumnValue struct {
	*nilColumnValue
}

func NewNilStringColumnValue() ColumnValue {
	return &NilStringColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

func (n *NilStringColumnValue) Type() ColumnType {
	return TypeString
}

func (n *NilStringColumnValue) clone() ColumnValue {
	return NewNilStringColumnValue()
}

//StringColumnValue 注意：Decimal 123.0（val:1230,exp:-1）和123（val:123,exp:0）不一致
type StringColumnValue struct {
	*notNilColumnValue
	TimeEncoder
	val string
}

func NewStringColumnValue(s string) ColumnValue {
	return NewStringColumnValueWithEncoder(s, NewStringTimeEncoder(time.RFC3339Nano))
}

func NewStringColumnValueWithEncoder(s string, e TimeEncoder) ColumnValue {
	return &StringColumnValue{
		notNilColumnValue: &notNilColumnValue{},
		TimeEncoder:       e,
		val:               s,
	}
}

func (s *StringColumnValue) Type() ColumnType {
	return TypeString
}

func (s *StringColumnValue) AsBool() (v bool, err error) {
	v, err = strconv.ParseBool(s.val)
	if err != nil {
		return false, NewTransformError(s.Type(), TypeBool, fmt.Errorf("err: %v, val: %v ", err, s.val))
	}
	return
}

func (s *StringColumnValue) AsBigInt() (*big.Int, error) {
	v, err := NewDecimalColumnValueFromString(s.val)
	if err != nil {
		return nil, NewTransformError(s.Type(), TypeBigInt, fmt.Errorf("err: %v, val: %v ", err, s.val))
	}
	return v.AsBigInt()
}

func (s *StringColumnValue) AsDecimal() (decimal.Decimal, error) {
	v, err := NewDecimalColumnValueFromString(s.val)
	if err != nil {
		return decimal.Decimal{}, NewTransformError(s.Type(), TypeDecimal, fmt.Errorf("err: %v, val: %v ", err, s.val))
	}
	return v.AsDecimal()
}

func (s *StringColumnValue) AsString() (string, error) {
	return s.val, nil
}

func (s *StringColumnValue) AsBytes() ([]byte, error) {
	return []byte(s.val), nil
}

func (s *StringColumnValue) AsTime() (t time.Time, err error) {
	t, err = s.TimeEncode(s.val)
	if err != nil {
		return time.Time{}, NewTransformError(s.Type(), TypeTime, fmt.Errorf(" val: %v", s.val))
	}
	return
}

func (s *StringColumnValue) String() string {
	return s.val
}

func (s *StringColumnValue) clone() ColumnValue {
	return NewStringColumnValue(s.val)
}
