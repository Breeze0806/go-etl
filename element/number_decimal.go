// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package element

import (
	"errors"
	"math"
	"math/big"

	"github.com/cockroachdb/apd/v3"
)

var _DecimalZero = apd.New(0, 0)
var _StrZero = "0"

// DecimalStr   High-precision decimal string
type DecimalStr struct {
	value  string
	intLen int
}

// Bool   Convert to boolean
func (d *DecimalStr) Bool() (bool, error) {
	return d.value != _StrZero, nil
}

// Float64   Convert to 64-bit floating-point number
func (d *DecimalStr) Float64() (v float64, err error) {
	f, _ := new(big.Float).SetString(d.value)
	if v, _ = f.Float64(); math.Abs(v) > math.MaxFloat64 {
		v = 0.0
		err = errors.New("element: DecimalStr to float64 fail for out of range")
		return
	}
	return
}

// BigInt   Convert to high-precision integer
func (d *DecimalStr) BigInt() BigIntNumber {
	return convertBigInt(d.value[:d.intLen]).(BigIntNumber)
}

// Decimal   Convert to high-precision decimal
func (d *DecimalStr) Decimal() DecimalNumber {
	return d
}

// Decimal   Convert to string
func (d *DecimalStr) String() string {
	return d.value
}

// CloneDecimal   Clone high-precision decimal
func (d *DecimalStr) CloneDecimal() DecimalNumber {
	return &DecimalStr{
		value:  d.value,
		intLen: d.intLen,
	}
}

// AsDecimal   Convert to high-precision decimal
func (d *DecimalStr) AsDecimal() *apd.Decimal {
	intString := d.value
	if d.intLen+1 < len(d.value) {
		intString = d.value[:d.intLen] + d.value[d.intLen+1:]
	}
	v, _ := new(apd.BigInt).SetString(intString, 10)
	return apd.NewWithBigInt(v, int32(-len(d.value)+d.intLen+1))
}

// Decimal   High-precision decimal
type Decimal struct {
	value *apd.Decimal
}

// Bool   Convert to boolean
func (d *Decimal) Bool() (bool, error) {
	return d.value.Cmp(_DecimalZero) != 0, nil
}

// Float64   Convert to 64-bit floating-point number
func (d *Decimal) Float64() (v float64, err error) {
	f, _ := new(big.Float).SetString(d.String())
	if v, _ = f.Float64(); math.Abs(v) > math.MaxFloat64 {
		return 0, errors.New("element: Decimal to float64 fail out of range")
	}
	return v, nil
}

// BigInt   Convert to high-precision integer
func (d *Decimal) BigInt() BigIntNumber {
	exp := d.value.Exponent
	value := &d.value.Coeff
	if d.value.Negative {
		value = value.Neg(value)
	}
	diff := math.Abs(-float64(exp))
	expScale := new(apd.BigInt).Exp(_IntTen, apd.NewBigInt(int64(diff)), nil)
	if 0 > exp {
		value = value.Quo(value, expScale)
	} else if 0 < exp {
		value = value.Mul(value, expScale)
	}

	return &BigInt{
		value: value,
	}
}

// Decimal   Convert to high-precision decimal
func (d *Decimal) Decimal() DecimalNumber {
	return d
}

// Decimal   Convert to string
func (d *Decimal) String() string {
	s, _ := convertDecimal(d.value.Text('f'))
	return s.String()
}

// CloneDecimal   Clone high-precision decimal
func (d *Decimal) CloneDecimal() DecimalNumber {
	return &Decimal{
		value: d.value.Set(d.value),
	}
}

// AsDecimal   Convert to high-precision decimal
func (d *Decimal) AsDecimal() *apd.Decimal {
	return d.value
}
