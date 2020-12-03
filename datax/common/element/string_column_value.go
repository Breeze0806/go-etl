package element

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type NilStringColumnValue struct {
}

func (n *NilStringColumnValue) Type() ColumnType {
	return TypeString
}

type StringColumnValue struct {
	notNilColumnValue
	val string
}

func NewStringColumnValue(s string) (ColumnValue, error) {
	return &StringColumnValue{
		val: s,
	}, nil
}

func (s *StringColumnValue) Type() ColumnType {
	return TypeString
}

func (s *StringColumnValue) AsBool() (bool, error) {
	return strconv.ParseBool(s.val)
}

func (s *StringColumnValue) AsBigInt() (*big.Int, error) {
	if v, ok := new(big.Int).SetString(s.val, 10); ok {
		return v, nil
	}
	return nil, NewTransformError(s.Type(), TypeBigInt, fmt.Errorf("val: %v ", s.val))
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

func (s *StringColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, NewTransformError(s.Type(), TypeTime, fmt.Errorf(" val: %v", s.val))
}

func (s *StringColumnValue) String() string {
	return s.val
}

func (s *StringColumnValue) clone() ColumnValue {
	return &StringColumnValue{
		val: s.val,
	}
}
