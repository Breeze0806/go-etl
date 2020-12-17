package element

import (
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

type mockTimeDecoder struct{}

func (m *mockTimeDecoder) TimeDecode(t time.Time) (interface{}, error) {
	return time.Time{}, fmt.Errorf("mockTimeDecoder error")
}

func testBigIntFromString(v string) *big.Int {
	bi, _ := new(big.Int).SetString(v, 10)
	return bi
}

func testBigIntColumnValueFromString(v string) *BigIntColumnValue {
	c, _ := NewBigIntColumnValueFromString(v)
	return c.(*BigIntColumnValue)
}

func testDecimalFormString(v string) decimal.Decimal {
	d, _ := decimal.NewFromString(v)
	return d
}

func testDecimalColumnValueFormString(v string) ColumnValue {
	d, _ := NewDecimalColumnValueFromString(v)
	return d
}
