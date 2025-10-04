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

var _IntZero = apd.NewBigInt(0)
var _IntTen = apd.NewBigInt(10)
var _IntFive = apd.NewBigInt(5)

// BigInt   Big integer
type BigInt struct {
	value *apd.BigInt
}

// Bool   Convert to boolean
func (b *BigInt) Bool() (bool, error) {
	return b.value.Cmp(_IntZero) != 0, nil
}

// Int64   Convert to 64-bit integer
func (b *BigInt) Int64() (int64, error) {
	if b.value.IsInt64() {
		return b.value.Int64(), nil
	}
	return 0, errors.New("element: BigInt to int64 fail for out of range")
}

// Float64   Convert to 64-bit floating-point number
func (b *BigInt) Float64() (v float64, err error) {
	f := new(big.Float).SetInt(b.value.MathBigInt())
	if v, _ = f.Float64(); math.Abs(v) > math.MaxFloat64 {
		v = 0.0
		err = errors.New("element: BigIntStr to float64 fail for out of range")
		return
	}
	return
}

// BigInt   Convert to high-precision integer
func (b *BigInt) BigInt() BigIntNumber {
	return b
}

// Decimal   Convert to high-precision decimal
func (b *BigInt) Decimal() DecimalNumber {
	return b
}

// Decimal   Convert to string
func (b *BigInt) String() string {
	return b.value.String()
}

// CloneBigInt   Clone high-precision integer
func (b *BigInt) CloneBigInt() BigIntNumber {
	return &BigInt{
		value: new(apd.BigInt).Set(b.value),
	}
}

// CloneDecimal   Clone high-precision decimal
func (b *BigInt) CloneDecimal() DecimalNumber {
	return &BigInt{
		value: new(apd.BigInt).Set(b.value),
	}
}

// AsBigInt   Convert to high-precision integer
func (b *BigInt) AsBigInt() *apd.BigInt {
	return b.value
}

// AsDecimal   Convert to high-precision decimal
func (b *BigInt) AsDecimal() *apd.Decimal {
	return apd.NewWithBigInt(b.AsBigInt(), 0)
}

// BigIntStr   High-precision integer string
type BigIntStr struct {
	value string
}

// Bool   Convert to boolean
func (b *BigIntStr) Bool() (bool, error) {
	return false, errors.New("element: BigIntStr to Bool fail for out of range")
}

// Int64   Convert to 64-bit integer
func (b *BigIntStr) Int64() (int64, error) {
	return 0, errors.New("element: BigIntStr to int64 fail for out of range")
}

// Float64   Convert to 64-bit floating-point number
func (b *BigIntStr) Float64() (v float64, err error) {
	f, _ := new(big.Float).SetString(b.value)
	if v, _ = f.Float64(); math.Abs(v) > math.MaxFloat64 {
		v = 0.0
		err = errors.New("element: BigIntStr to float64 fail for out of range")
		return
	}
	return
}

// BigInt   Convert to high-precision integer
func (b *BigIntStr) BigInt() BigIntNumber {
	return b
}

// Decimal   Convert to high-precision decimal
func (b *BigIntStr) Decimal() DecimalNumber {
	return b
}

// Decimal   Convert to string
func (b *BigIntStr) String() string {
	return b.value
}

// CloneBigInt   Clone high-precision integer
func (b *BigIntStr) CloneBigInt() BigIntNumber {
	return &BigIntStr{
		value: b.value,
	}
}

// CloneDecimal   Clone high-precision decimal
func (b *BigIntStr) CloneDecimal() DecimalNumber {
	return &BigIntStr{
		value: b.value,
	}
}

// AsBigInt   Convert to high-precision integer
func (b *BigIntStr) AsBigInt() *apd.BigInt {
	v, _ := new(apd.BigInt).SetString(b.value, 10)
	return v
}

// AsDecimal   Convert to high-precision decimal
func (b *BigIntStr) AsDecimal() *apd.Decimal {
	return apd.NewWithBigInt(b.AsBigInt(), 0)
}
