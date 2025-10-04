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
	"fmt"
	"math"
	"strconv"

	"github.com/cockroachdb/apd/v3"
)

var _DefaultNumberConverter NumberConverter = &Converter{}

// NumberConverter   Number converter
type NumberConverter interface {
	ConvertBigIntFromInt(i int64) (num BigIntNumber)
	ConvertBigIntFromUint(i uint64) (num BigIntNumber)
	ConvertDecimalFromFloat32(f float32) (num DecimalNumber)
	ConvertDecimalFromFloat(f float64) (num DecimalNumber)
	ConvertBigInt(s string) (num BigIntNumber, err error)
	ConvertDecimal(s string) (num DecimalNumber, err error)
}

// OldConverter   Unchecked conversion
type OldConverter struct{}

// ConvertBigIntFromInt   Convert to decimal from integer
func (c *OldConverter) ConvertBigIntFromInt(i int64) (num BigIntNumber) {
	return &BigInt{
		value: apd.NewBigInt(i),
	}
}

// ConvertBigIntFromUint   Convert to decimal from unsigned integer
func (c *OldConverter) ConvertBigIntFromUint(i uint64) (num BigIntNumber) {
	return &BigInt{
		value: new(apd.BigInt).SetUint64(i),
	}
}

// ConvertDecimalFromFloat32   Convert to decimal from 32-bit loating-point number
func (c *OldConverter) ConvertDecimalFromFloat32(f float32) (num DecimalNumber) {
	return &Decimal{
		value: NewApdDecimalFromFloat32(f),
	}
}

// ConvertDecimalFromFloat   Convert to decimal from floating-point number
func (c *OldConverter) ConvertDecimalFromFloat(f float64) (num DecimalNumber) {
	return &Decimal{
		value: NewApdDecimalFromFloat(f),
	}
}

// ConvertDecimal   Convert string to decimal
func (c *OldConverter) ConvertDecimal(s string) (num DecimalNumber, err error) {
	var d *apd.Decimal
	if d, _, err = apd.NewFromString(s); err != nil {
		return
	}
	num = &Decimal{
		value: d,
	}
	return
}

// ConvertBigInt   Convert string to integer
func (c *OldConverter) ConvertBigInt(s string) (num BigIntNumber, err error) {
	b, ok := new(apd.BigInt).SetString(s, 10)
	if !ok {
		err = errors.New("number is not int")
		return
	}
	num = &BigInt{
		value: b,
	}
	return
}

// Converter   Number converter
type Converter struct{}

// ConvertBigIntFromInt   Convert to decimal from integer
func (c *Converter) ConvertBigIntFromInt(i int64) (num BigIntNumber) {
	return &Int64{
		value: i,
	}
}

// ConvertBigIntFromUint   Convert to decimal from unsigned integer
func (c *Converter) ConvertBigIntFromUint(i uint64) (num BigIntNumber) {
	return &Uint64{
		value: i,
	}
}

// ConvertDecimalFromFloat32   Convert to decimal from 32-bit loating-point number
func (c *Converter) ConvertDecimalFromFloat32(f float32) (num DecimalNumber) {
	return &Decimal{
		value: NewApdDecimalFromFloat32(f),
	}
}

// ConvertDecimalFromFloat   Convert to decimal from floating-point number
func (c *Converter) ConvertDecimalFromFloat(f float64) (num DecimalNumber) {
	return &Float64{
		value: f,
	}
}

// ConvertDecimal  Convert string to decimal
// inspired by https://github.com/shopspring/decimal
func (c *Converter) ConvertDecimal(s string) (num DecimalNumber, err error) {
	eIndex := -1
	for i, r := range s {
		if r == 'E' || r == 'e' {
			if eIndex > -1 {
				return nil, fmt.Errorf("can't convert %s to decimal: multiple 'E' characters found", s)
			}
			eIndex = i
			continue
		}
	}

	if eIndex == -1 {
		return convertDecimal(s)
	}
	num, err = convertDecimal(s[:eIndex])
	if err != nil {
		return nil, err
	}

	var expInt int64
	if expInt, err = strconv.ParseInt(s[eIndex+1:], 10, 32); err != nil {
		return nil, fmt.Errorf("can't convert %s to decimal: fractional part too long", s)
	}
	exp := expInt
	var dValue *apd.BigInt
	switch data := num.(type) {
	case *Int64:
		dValue = apd.NewBigInt(data.value)
	case *BigIntStr:
		dValue, _ = new(apd.BigInt).SetString(data.value, 10)
	case *DecimalStr:
		intString := data.value
		if data.intLen+1 < len(data.value) {
			intString = data.value[:data.intLen] + data.value[data.intLen+1:]
			expInt := -len(data.value[data.intLen+1:])
			exp += int64(expInt)
		}
		if len(intString) <= 18 {
			parsed64, _ := strconv.ParseInt(intString, 10, 64)
			dValue = apd.NewBigInt(parsed64)
		} else {
			dValue, _ = new(apd.BigInt).SetString(intString, 10)
		}
	}

	if exp < math.MinInt32 || exp > math.MaxInt32 {
		// NOTE(vadim): I doubt a string could realistically be this long
		return nil, fmt.Errorf("can't convert %s to decimal: fractional part too long", s)
	}
	return &Decimal{
		value: apd.NewWithBigInt(dValue, int32(exp)),
	}, nil
}

// ConvertBigInt converts a string to an integer
// inspired by https://github.com/shopspring/decimal
func (c *Converter) ConvertBigInt(s string) (num BigIntNumber, err error) {
	if len(s) == 0 {
		err = errors.New("element: convertDecimal empty string")
		return
	}
	sign := ""
	if s[0] == '+' || s[0] == '-' {
		if s[0] == '-' {
			sign = "-"
		}
		s = s[1:]
	}

	if len(s) == 0 {
		err = errors.New("element: convertDecimal empty string")
		return
	}

	start := len(s)
	for i := 0; i < len(s); i++ {
		if s[i] != '0' {
			start = i
			break
		}
	}
	s = s[start:]

	if err = checkInt(s); err != nil {
		return
	}
	first := s
	if len(first) == 0 {
		first = "0"
	}
	if first == "0" {
		sign = ""
	}
	return convertBigInt(sign + first).(BigIntNumber), nil
}

func convertBigInt(s string) (n Number) {
	v, err := parseInt64(s)
	if err != nil {
		return &BigIntStr{
			value: s,
		}
	}
	return &Int64{
		value: v,
	}
}

func convertDecimal(s string) (num DecimalNumber, err error) {
	if len(s) == 0 {
		err = errors.New("element: convertDecimal empty string")
		return
	}
	sign := ""
	if s[0] == '+' || s[0] == '-' {
		if s[0] == '-' {
			sign = "-"
		}
		s = s[1:]
	}

	if len(s) == 0 {
		err = errors.New("element: convertDecimal empty string")
		return
	}

	if s[0] == '.' && len(s[1:]) == 0 {
		err = errors.New("element: convertDecimal only dot")
		return
	}

	start := len(s)
	for i := 0; i < len(s); i++ {
		if s[i] != '0' {
			start = i
			break
		}
	}
	s = s[start:]
	pIndex := -1

	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			if pIndex > -1 {
				return nil, errors.New("element: convertDecimal too many dots")
			}
			pIndex = i
		}
	}

	if pIndex != -1 {
		end := -1
		for i := len(s) - 1; i >= 0; i-- {
			if s[i] != '0' {
				end = i
				break
			}
		}
		s = s[:end+1]

		if err = checkInt(s[:pIndex]); err != nil {
			return
		}

		if err = checkInt(s[pIndex+1:]); err != nil {
			return
		}
		first := s[:pIndex]
		if len(first) == 0 {
			first = "0"
		}
		if len(s[pIndex+1:]) == 0 {
			if first == "0" {
				sign = ""
			}
			return &DecimalStr{
				value:  sign + first,
				intLen: len(sign) + len(first),
			}, nil
		}

		return &DecimalStr{
			value:  sign + first + s[pIndex:],
			intLen: len(sign) + len(first),
		}, nil
	}

	if err = checkInt(s); err != nil {
		return
	}

	first := s
	if len(first) == 0 {
		first = "0"
	}

	return convertBigInt(sign + first).(DecimalNumber), nil
}

func checkInt(s string) (err error) {
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			err = errors.New("element: convertDecimal invalid syntax")
			return
		}
	}
	return
}
