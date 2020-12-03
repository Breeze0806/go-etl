package element

import (
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

var defaultTimeFormat = time.RFC3339Nano

type NilTimeColumnValue struct {
	nilColumnValue
}

func (n *NilTimeColumnValue) Type() ColumnType {
	return TypeTime
}

type TimeColumnValue struct {
	notNilColumnValue
	val time.Time
}

func (t *TimeColumnValue) Type() ColumnType {
	return TypeTime
}

func (t *TimeColumnValue) AsBool() (bool, error) {
	return false, NewTransformError(t.Type(), TypeBool, fmt.Errorf("val: %v", t.String()))
}

func (t *TimeColumnValue) AsBigInt() (*big.Int, error) {
	return nil, NewTransformError(t.Type(), TypeBigInt, fmt.Errorf("val: %v", t.String()))
}

func (t *TimeColumnValue) AsDecimal() (decimal.Decimal, error) {
	return decimal.Decimal{}, NewTransformError(t.Type(), TypeDecimal, fmt.Errorf("val: %v", t.String()))
}

func (t *TimeColumnValue) AsString() (string, error) {
	return t.val.String(), nil
}

func (t *TimeColumnValue) AsBytes() ([]byte, error) {
	return []byte(t.val.String()), nil
}

func (t *TimeColumnValue) AsTime() (time.Time, error) {
	return t.val, nil
}

func (t *TimeColumnValue) String() string {
	return t.val.Format(defaultTimeFormat)
}

func (t *TimeColumnValue) clone() ColumnValue {
	return &TimeColumnValue{
		val: t.val,
	}
}
