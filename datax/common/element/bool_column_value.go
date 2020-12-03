package element

import (
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

type NilBoolColumnValue struct {
	nilColumnValue
}

func (n *NilBoolColumnValue) Type() ColumnType {
	return TypeBool
}

type BoolColumnValue struct {
	notNilColumnValue
	val bool
}

func (b *BoolColumnValue) Type() ColumnType {
	return TypeBool
}

func (b *BoolColumnValue) AsBool() (bool, error) {
	return b.val, nil
}

func (b *BoolColumnValue) AsBigInt() (*big.Int, error) {
	if b.val {
		return _IntOne, nil
	}
	return _IntZero, nil
}

func (b *BoolColumnValue) AsDecimal() (decimal.Decimal, error) {
	if b.val {
		return decimal.New(1, 0), nil
	}
	return decimal.New(0, 0), nil
}

func (b *BoolColumnValue) AsString() (string, error) {
	if b.val {
		return b.String(), nil
	}
	return b.String(), nil
}

func (b *BoolColumnValue) AsBytes() ([]byte, error) {
	if b.val {
		return []byte(b.String()), nil
	}
	return []byte(b.String()), nil
}

func (b *BoolColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, NewTransformError(b.Type(), TypeTime, fmt.Errorf(" val: %v", b.String()))
}

func (b *BoolColumnValue) String() string {
	if b.val {
		return "true"
	}
	return "false"
}

func (b *BoolColumnValue) clone() ColumnValue {
	return &BoolColumnValue{
		val: b.val,
	}
}
