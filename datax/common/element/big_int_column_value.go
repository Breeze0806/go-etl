package element

import (
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

var _IntZero = big.NewInt(0)
var _IntOne = big.NewInt(1)
var _IntTen = big.NewInt(10)

type NilBigIntColumnValue struct {
	nilColumnValue
}

func (n *NilBigIntColumnValue) Type() ColumnType {
	return TypeBigInt
}

type BigIntColumnValue struct {
	notNilColumnValue
	val *big.Int
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
	return decimal.NewFromBigInt(b.val, 0), nil
}

func (b *BigIntColumnValue) AsString() (string, error) {
	return b.val.String(), nil
}

func (b *BigIntColumnValue) AsBytes() ([]byte, error) {
	return b.val.Bytes(), nil
}

func (b *BigIntColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, NewTransformError(b.Type(), TypeTime, fmt.Errorf(" val: %v", b.String()))
}

func (b *BigIntColumnValue) String() string {
	return b.val.String()
}

func (b *BigIntColumnValue) clone() ColumnValue {
	return &BigIntColumnValue{
		val: new(big.Int).Set(b.val),
	}
}
