package element

import (
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

var _IntZero = big.NewInt(0)
var _IntTen = big.NewInt(10)

type NilBigIntColumnValue struct {
	*nilColumnValue
}

func NewNilBigIntColumnValue() ColumnValue {
	return &NilBigIntColumnValue{
		nilColumnValue: &nilColumnValue{},
	}
}

func (n *NilBigIntColumnValue) Type() ColumnType {
	return TypeBigInt
}

func (n *NilBigIntColumnValue) clone() ColumnValue {
	return NewNilBigIntColumnValue()
}

type BigIntColumnValue struct {
	notNilColumnValue
	val *big.Int
}

func NewBigIntColumnValueFromInt64(v int64) ColumnValue {
	return &BigIntColumnValue{
		val: big.NewInt(v),
	}
}

func NewBigIntColumnValue(v *big.Int) ColumnValue {
	return &BigIntColumnValue{
		val: new(big.Int).Set(v),
	}
}

func NewBigIntColumnValueFromString(v string) (ColumnValue, error) {
	bi, ok := new(big.Int).SetString(v, 10)
	if !ok {
		return nil, NewSetError(v, TypeBigInt, fmt.Errorf("string %v is not valid int", v))
	}
	return &BigIntColumnValue{
		val: bi,
	}, nil
}

func (b *BigIntColumnValue) Type() ColumnType {
	return TypeBigInt
}

func (b *BigIntColumnValue) AsBool() (bool, error) {
	return b.val.Cmp(_IntZero) != 0, nil
}

func (b *BigIntColumnValue) AsBigInt() (*big.Int, error) {
	return new(big.Int).Set(b.val), nil
}

func (b *BigIntColumnValue) AsDecimal() (decimal.Decimal, error) {
	if b.val.Cmp(_IntZero) != 0 {
		return decimal.NewFromBigInt(b.val, 0), nil
	}
	return decimal.New(0, 1), nil
}

func (b *BigIntColumnValue) AsString() (string, error) {
	return b.val.String(), nil
}

func (b *BigIntColumnValue) AsBytes() ([]byte, error) {
	return []byte(b.val.String()), nil
}

func (b *BigIntColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, NewTransformErrorFormColumnTypes(b.Type(), TypeTime, fmt.Errorf(" val: %v", b.String()))
}

func (b *BigIntColumnValue) String() string {
	return b.val.String()
}

func (b *BigIntColumnValue) clone() ColumnValue {
	return NewBigIntColumnValue(b.val)
}
