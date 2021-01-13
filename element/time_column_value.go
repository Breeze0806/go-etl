package element

import (
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

type NilTimeColumnValue struct {
	nilColumnValue
}

func NewNilTimeColumnValue() ColumnValue {
	return &NilTimeColumnValue{}
}

func (n *NilTimeColumnValue) Type() ColumnType {
	return TypeTime
}

func (n *NilTimeColumnValue) Clone() ColumnValue {
	return NewNilTimeColumnValue()
}

type TimeColumnValue struct {
	notNilColumnValue
	TimeDecoder
	val time.Time
}

func NewTimeColumnValue(t time.Time) ColumnValue {
	return NewTimeColumnValueWithDecoder(t, NewStringTimeDecoder(time.RFC3339Nano))
}

func NewTimeColumnValueWithDecoder(t time.Time, d TimeDecoder) ColumnValue {
	return &TimeColumnValue{
		TimeDecoder:       d,
		val:               t,
	}
}

func (t *TimeColumnValue) Type() ColumnType {
	return TypeTime
}

func (t *TimeColumnValue) AsBool() (bool, error) {
	return false, NewTransformErrorFormColumnTypes(t.Type(), TypeBool, fmt.Errorf("val: %v", t.String()))
}

func (t *TimeColumnValue) AsBigInt() (*big.Int, error) {
	return nil, NewTransformErrorFormColumnTypes(t.Type(), TypeBigInt, fmt.Errorf("val: %v", t.String()))
}

func (t *TimeColumnValue) AsDecimal() (decimal.Decimal, error) {
	return decimal.Decimal{}, NewTransformErrorFormColumnTypes(t.Type(), TypeDecimal, fmt.Errorf("val: %v", t.String()))
}

func (t *TimeColumnValue) AsString() (s string, err error) {
	var i interface{}
	i, err = t.TimeDecode(t.val)
	if err != nil {
		return "", NewTransformErrorFormColumnTypes(t.Type(), TypeString, fmt.Errorf("val: %v", t.String()))
	}
	return i.(string), nil
}

func (t *TimeColumnValue) AsBytes() (b []byte, err error) {
	var i interface{}
	i, err = t.TimeDecode(t.val)
	if err != nil {
		return nil, NewTransformErrorFormColumnTypes(t.Type(), TypeString, fmt.Errorf("val: %v", t.String()))
	}
	return []byte(i.(string)), nil
}

func (t *TimeColumnValue) AsTime() (time.Time, error) {
	return t.val, nil
}

func (t *TimeColumnValue) String() string {
	return t.val.Format(defaultTimeFormat)
}

func (t *TimeColumnValue) Clone() ColumnValue {
	return &TimeColumnValue{
		val: t.val,
	}
}
