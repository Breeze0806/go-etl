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
	"strconv"

	"github.com/cockroachdb/apd/v3"
)

// Int64   64-bit integer
type Int64 struct {
	value int64
}

// Bool   Convert to boolean
func (i *Int64) Bool() (bool, error) {
	return i.value != 0, nil
}

// Int64   Convert to 64-bit integer
func (i *Int64) Int64() (int64, error) {
	return i.value, nil
}

// Float64   Convert to 64-bit floating-point number. But it will result in loss of precision.
func (i *Int64) Float64() (float64, error) {
	return float64(i.value), nil
}

// BigInt   Convert to high-precision integer
func (i *Int64) BigInt() BigIntNumber {
	return i
}

// Decimal   Convert to high-precision decimal
func (i *Int64) Decimal() DecimalNumber {
	return i
}

// Decimal   Convert to string
func (i *Int64) String() string {
	return FormatInt64(i.value)
}

// CloneBigInt   Clone high-precision integer
func (i *Int64) CloneBigInt() BigIntNumber {
	return &Int64{
		value: i.value,
	}
}

// CloneDecimal   Clone high-precision decimal
func (i *Int64) CloneDecimal() DecimalNumber {
	return &Int64{
		value: i.value,
	}
}

// AsBigInt   Convert to high-precision integer
func (i *Int64) AsBigInt() *apd.BigInt {
	return apd.NewBigInt(i.value)
}

// AsDecimal   Convert to high-precision decimal
func (i *Int64) AsDecimal() *apd.Decimal {
	return apd.New(i.value, 0)
}

// Uint64  ungined 64-bit integer
type Uint64 struct {
	value uint64
}

// Bool   Convert to boolean
func (i *Uint64) Bool() (bool, error) {
	return i.value != 0, nil
}

// Int64   Convert to 64-bit integer
func (i *Uint64) Int64() (int64, error) {
	if i.value > uint64(math.MaxInt64) {
		return 0, errors.New("element: uint64 to int64 fail for out of range")
	}
	return int64(i.value), nil
}

// Float64   Convert to 64-bit floating-point number. But it will result in loss of precision.
func (i *Uint64) Float64() (float64, error) {
	return float64(i.value), nil
}

// BigInt   Convert to high-precision integer
func (i *Uint64) BigInt() BigIntNumber {
	return i
}

// Decimal   Convert to high-precision decimal
func (i *Uint64) Decimal() DecimalNumber {
	return i
}

// Decimal   Convert to string
func (i *Uint64) String() string {
	return FormatUInt64(i.value)
}

// CloneBigInt   Clone high-precision integer
func (i *Uint64) CloneBigInt() BigIntNumber {
	return &Uint64{
		value: i.value,
	}
}

// CloneDecimal   Clone high-precision decimal
func (i *Uint64) CloneDecimal() DecimalNumber {
	return &Uint64{
		value: i.value,
	}
}

// AsBigInt   Convert to high-precision integer
func (i *Uint64) AsBigInt() (bi *apd.BigInt) {
	bi = new(apd.BigInt).SetUint64(i.value)
	return
}

// AsDecimal   Convert to high-precision decimal
func (i *Uint64) AsDecimal() *apd.Decimal {
	return apd.NewWithBigInt(i.AsBigInt(), 0)
}

// base on stdlib strconv , which has the following license:
// """
// BSD 3-Clause "New" or "Revised" License
// Copyright 2009 The Go Authors.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google LLC nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// """

func FormatInt64(i int64) (s string) {
	if 0 <= i && i < nSmalls {
		return small(int(i))
	}
	return formatBits(uint64(i), i < 0)
}

func FormatUInt64(i uint64) (s string) {
	if i < nSmalls {
		return small(int(i))
	}
	return formatBits(uint64(i), false)
}

const maxUInt63 = 1 << 63

const smallsString = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"
const digits = "0123456789"

const host32bit = ^uint(0)>>32 == 0

const nSmalls = 100

// small returns the string for an i with 0 <= i < nSmalls.
func small(i int) string {
	if i < 10 {
		return digits[i : i+1]
	}
	return smallsString[i*2 : i*2+2]
}

func formatBits(u uint64, neg bool) (s string) {
	var a [21 + 1]byte // +1 for sign of 64bit value in base 2
	i := len(a)

	if neg {
		u = -u
	}

	// common case: use constants for / because
	// the compiler can optimize it into a multiply+shift
	if host32bit {
		// convert the lower digits using 32bit operations
		for u >= 1e9 {
			// Avoid using r = a%b in addition to q = a/b
			// since 64bit division and modulo operations
			// are calculated by runtime functions on 32bit machines.
			q := u / 1e9
			us := uint(u - q*1e9) // u % 1e9 fits into a uint
			for j := 4; j > 0; j-- {
				is := us % 100 * 2
				us /= 100
				i -= 2
				a[i+1] = smallsString[is+1]
				a[i+0] = smallsString[is+0]
			}

			// us < 10, since it contains the last digit
			// from the initial 9-digit us.
			i--
			a[i] = smallsString[us*2+1]

			u = q
		}
		// u < 1e9
	}

	// u guaranteed to fit into a uint
	us := uint(u)
	for us >= 100 {
		is := us % 100 * 2
		us /= 100
		i -= 2
		a[i+1] = smallsString[is+1]
		a[i+0] = smallsString[is+0]
	}

	// us < 100
	is := us * 2
	i--
	a[i] = smallsString[is+1]
	if us >= 10 {
		i--
		a[i] = smallsString[is]
	}

	// add sign, if any
	if neg {
		i--
		a[i] = '-'
	}

	s = string(a[i:])
	return
}

func parseInt64(s string) (i int64, err error) {
	const fnParseInt = "ParseInt"
	// Pick off leading sign.
	s0 := s
	neg := false
	if s[0] == '-' {
		neg = true
		s = s[1:]
	}

	// Convert unsigned and check range.
	var un uint64
	un = parseUint64(s)

	cutoff := uint64(maxUInt63)
	if !neg && un >= cutoff {
		return int64(cutoff - 1), rangeError(fnParseInt, s0)
	}
	if neg && un > cutoff {
		return -int64(cutoff), rangeError(fnParseInt, s0)
	}
	n := int64(un)
	if neg {
		n = -n
	}
	return n, nil
}

func parseUint64(s string) (n uint64) {
	const fnParseUint = "ParseUint"
	cutoff := uint64(math.MaxUint64/10 + 1)

	maxVal := uint64(math.MaxUint64)

	for _, c := range s {
		d := c - '0'

		if n >= cutoff {
			// n*base overflows
			return maxVal
		}
		n *= uint64(10)

		n1 := n + uint64(d)
		if n1 < n || n1 > maxVal {
			// n+d overflows
			return maxVal
		}
		n = n1
	}

	return n
}

func rangeError(fn, str string) *strconv.NumError {
	return &strconv.NumError{Func: fn, Num: cloneString(str), Err: strconv.ErrRange}
}

func cloneString(x string) string { return string([]byte(x)) }
