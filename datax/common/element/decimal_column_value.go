package element

import (
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

type NilDecimalColumnValue struct {
	nilColumnValue
}

func (n *NilDecimalColumnValue) Type() ColumnType {
	return TypeDecimal
}

type DecimalColumnValue struct {
	notNilColumnValue
	val decimal.Decimal
}

func NewDecimalColumnValueFromFloat(f float64) (ColumnValue, error) {
	return &DecimalColumnValue{
		val: decimal.NewFromFloat(f),
	}, nil
}

func NewDecimalColumnValueFromString(s string) (ColumnValue, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}
	return &DecimalColumnValue{
		val: d,
	}, nil
}

func (d *DecimalColumnValue) Type() ColumnType {
	return TypeDecimal
}

func (d *DecimalColumnValue) AsBool() (bool, error) {
	return d.val.Cmp(decimal.Zero) != 0, nil
}

func (d *DecimalColumnValue) AsBigInt() (*big.Int, error) {
	exp := d.val.Exponent()
	value := d.val.Coefficient()
	if exp == 0 {
		return value, nil
	}
	diff := math.Abs(-float64(exp))

	expScale := new(big.Int).Exp(_IntTen, big.NewInt(int64(diff)), nil)
	if 0 > exp {
		value = value.Quo(value, expScale)
	} else if 0 < exp {
		value = value.Mul(value, expScale)
	}

	return value, nil
}

func (d *DecimalColumnValue) AsDecimal() (decimal.Decimal, error) {
	return d.val, nil
}

func (d *DecimalColumnValue) AsString() (string, error) {
	return d.val.String(), nil
}

func (d *DecimalColumnValue) AsBytes() ([]byte, error) {
	return []byte(d.val.String()), nil
}

func (d *DecimalColumnValue) AsTime() (time.Time, error) {
	return time.Time{}, NewTransformError(d.Type(), TypeTime, fmt.Errorf(" val: %v", d.String()))
}

func (d *DecimalColumnValue) String() string {
	return d.val.String()
}

func (d *DecimalColumnValue) clone() ColumnValue {
	return &DecimalColumnValue{
		val: d.val,
	}
}
