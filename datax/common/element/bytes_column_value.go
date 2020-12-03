package element

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type NilBytesColumnValue struct {
}

func (n *NilBytesColumnValue) Type() ColumnType {
	return TypeBytes
}

type BytesColumnValue struct {
	notNilColumnValue
	val []byte
}

func (b *BytesColumnValue) Type() ColumnType {
	return TypeBytes
}

func (b *BytesColumnValue) AsBool() (bool, error) {
	v, err := strconv.ParseBool(b.String())
	if err != nil {
		return false, NewTransformError(b.Type(), TypeBool, fmt.Errorf("err: %v val: %v", err, b.String()))
	}
	return v, nil
}

func (b *BytesColumnValue) AsBigInt() (*big.Int, error) {

	if v, ok := new(big.Int).SetString(b.String(), 10); ok {
		return v, nil
	}
	return nil, NewTransformError(b.Type(), TypeBigInt, fmt.Errorf("val: %v ", b.String()))
}

func (b *BytesColumnValue) AsDecimal() (decimal.Decimal, error) {
	v, err := NewDecimalColumnValueFromString(b.String())
	if err != nil {
		return decimal.Decimal{}, NewTransformError(b.Type(), TypeDecimal, fmt.Errorf("err: %v, val: %v ", err, b.String()))
	}
	return v.AsDecimal()
}

func (b *BytesColumnValue) AsString() (string, error) {
	return b.String(), nil
}

func (b *BytesColumnValue) AsBytes() ([]byte, error) {
	v := make([]byte, len(b.val))
	copy(v, b.val)
	return v, nil
}

func (b *BytesColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, NewTransformError(b.Type(), TypeTime, fmt.Errorf(" val: %v", b.String()))
}

func (b *BytesColumnValue) String() string {
	return string(b.val)
}

func (b *BytesColumnValue) clone() ColumnValue {
	v := make([]byte, len(b.val))
	copy(v, b.val)
	return &BytesColumnValue{
		val: v,
	}
}
