package element

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type NilBytesColumnValue struct {
	*nilColumnValue
}

func NewNilBytesColumnValue() ColumnValue {
	return &NilBytesColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

func (n *NilBytesColumnValue) Type() ColumnType {
	return TypeBytes
}

func (n *NilBytesColumnValue) clone() ColumnValue {
	return NewNilBytesColumnValue()
}

type BytesColumnValue struct {
	*notNilColumnValue
	TimeEncoder
	val []byte
}

func NewBytesColumnValue(v []byte) ColumnValue {
	return NewBytesColumnValueWithEncoder(v, NewStringTimeEncoder(time.RFC3339Nano))
}

func NewBytesColumnValueWithEncoder(v []byte, e TimeEncoder) ColumnValue {
	return &BytesColumnValue{
		notNilColumnValue: &notNilColumnValue{},
		val:               v,
		TimeEncoder:       e,
	}
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
	v, err := NewDecimalColumnValueFromString(b.String())
	if err != nil {
		return nil, NewTransformError(b.Type(), TypeBigInt, fmt.Errorf("err: %v, val: %v ", err, b.String()))
	}
	return v.AsBigInt()
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

func (b *BytesColumnValue) AsTime() (t time.Time, err error) {
	t, err = b.TimeEncode(b.String())
	if err != nil {
		return time.Time{}, NewTransformError(b.Type(), TypeTime, fmt.Errorf(" val: %v", b.String()))
	}
	return
}

func (b *BytesColumnValue) String() string {
	return string(b.val)
}

func (b *BytesColumnValue) clone() ColumnValue {
	v := make([]byte, len(b.val))
	copy(v, b.val)
	return NewBytesColumnValue(v)
}
