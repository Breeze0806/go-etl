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
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"
)

// inspired by https://github.com/shopspring/decimal
type testDecimalStr struct {
	float  float64
	short  string
	exact  string
	intLen int
}

type testEnt struct {
	float   float64
	short   string
	exact   string
	inexact string
}

var testTable = []*testEnt{
	{3.141592653589793, "3.141592653589793", "", "3.14159265358979300000000000000000000000000000000000004"},
	{3, "3", "", "3.0000000000000000000000002"},
	{1234567890123456, "1234567890123456", "", "1234567890123456.00000000000000002"},
	{1234567890123456000, "1234567890123456000", "", "1234567890123456000.0000000000000008"},
	{1234.567890123456, "1234.567890123456", "", "1234.5678901234560000000000000009"},
	{.1234567890123456, "0.1234567890123456", "", "0.12345678901234560000000000006"},
	{0, "0", "", "0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"},
	{.1111111111111110, "0.111111111111111", "", "0.111111111111111000000000000000009"},
	{.1111111111111111, "0.1111111111111111", "", "0.111111111111111100000000000000000000023423545644534234"},
	{.1111111111111119, "0.1111111111111119", "", "0.111111111111111900000000000000000000000000000000000134123984192834"},
	{.000000000000000001, "0.000000000000000001", "", "0.00000000000000000100000000000000000000000000000000012341234"},
	{.000000000000000002, "0.000000000000000002", "", "0.0000000000000000020000000000000000000012341234123"},
	{.000000000000000003, "0.000000000000000003", "", "0.00000000000000000299999999999999999999999900000000000123412341234"},
	{.000000000000000005, "0.000000000000000005", "", "0.00000000000000000500000000000000000023412341234"},
	{.000000000000000008, "0.000000000000000008", "", "0.0000000000000000080000000000000000001241234432"},
	{.1000000000000001, "0.1000000000000001", "", "0.10000000000000010000000000000012341234"},
	{.1000000000000002, "0.1000000000000002", "", "0.10000000000000020000000000001234123412"},
	{.1000000000000003, "0.1000000000000003", "", "0.1000000000000003000000000000001234123412"},
	{.1000000000000005, "0.1000000000000005", "", "0.1000000000000005000000000000000006441234"},
	{.1000000000000008, "0.1000000000000008", "", "0.100000000000000800000000000000000009999999999999999999999999999"},
	{1e25, "10000000000000000000000000", "", ""},
	{1.5e14, "150000000000000", "", ""},
	{1.5e15, "1500000000000000", "", ""},
	{1.5e16, "15000000000000000", "", ""},
	{1.0001e25, "10001000000000000000000000", "", ""},
	{1.0001000000000000033e25, "10001000000000000000000000", "", ""},
	{2e25, "20000000000000000000000000", "", ""},
	{4e25, "40000000000000000000000000", "", ""},
	{8e25, "80000000000000000000000000", "", ""},
	{1e250, "10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", "", ""},
	{2e250, "20000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", "", ""},
	{math.MaxInt64, strconv.FormatFloat(float64(math.MaxInt64), 'f', -1, 64), "", FormatInt64(math.MaxInt64)},
	{1.29067116156722e-309, "0.00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000129067116156722", "", "0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001290671161567218558822290567835270536800098852722416870074139002112543896676308448335063375297788379444685193974290737962187240854947838776604607190387984577130572928111657710645015086812756013489109884753559084166516937690932698276436869274093950997935137476803610007959500457935217950764794724766740819156974617155861568214427828145972181876775307023388139991104942469299524961281641158436752347582767153796914843896176260096039358494077706152272661453132497761307744086665088096215425146090058519888494342944692629602847826300550628670375451325582843627504604013541465361435761965354140678551369499812124085312128659002910905639984075064968459581691226705666561364681985266583563078466180095375402399087817404368974165082030458595596655868575908243656158447265625000000000000000000000000000000000000004440000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
	// go Issue 29491.
	{498484681984085570, "498484681984085570", "", ""},
	{5.8339553793802237e+23, "583395537938022370000000", "", ""},
}

var testTableDecimalStr = []*testDecimalStr{
	{0, "+0000000103456789123456789012.123", "103456789123456789012.123", 21},
	{0, "0000000103456789123456789012.123", "103456789123456789012.123", 21},
	{0, "-0000000103456789123456789012.123", "-103456789123456789012.123", 22},
	{0, "-0000000103456789123456789012.12300000", "-103456789123456789012.123", 22},
	{0, "-0000000103456789123456789012.00000", "-103456789123456789012", 22},
	{0, "+0000000103456789123456789012.00000", "103456789123456789012", 21},
	{0, "-103456789123456789012.00000", "-103456789123456789012", 22},
	{0, "+103456789123456789012.00000", "103456789123456789012", 21},

	{0, "+.123", "0.123", 1},
	{0, "-.123", "-0.123", 2},
	{0, ".123", "0.123", 1},
	{0, "+0000000000.123", "0.123", 1},
	{0, "0000000000.123", "0.123", 1},
	{0, "-0000000000.123", "-0.123", 2},
	{0, "+0000000000.0000123", "0.0000123", 1},
	{0, "0000000000.0000123", "0.0000123", 1},
	{0, "-0000000000.0000123", "-0.0000123", 2},
	{0, "+0000000000.0000123", "0.0000123", 1},
	{0, "0000000000.000012301000", "0.000012301", 1},
	{0, "-0000000000.0000123", "-0.0000123", 2},
}

var testTableInt64 = map[string]string{
	"00000000000000000000000000000000000000000000000000":  "0",
	"+00000000000000000000000000000000000000000000000000": "0",
	"-00000000000000000000000000000000000000000000000000": "0",
	"-0":                              "0",
	"+0123456789012345678":            "123456789012345678",
	"+012345678901234567890":          "12345678901234567890",
	"+0000000103456789123456789012":   "103456789123456789012",
	"0000000103456789123456789012":    "103456789123456789012",
	"+103456789123456789012":          "103456789123456789012",
	"-103456789123456789012":          "-103456789123456789012",
	"103456789123456789012":           "103456789123456789012",
	"-0000000103456789123456789012":   "-103456789123456789012",
	FormatUInt64(math.MaxInt64 + 1):   FormatUInt64(math.MaxInt64 + 1),
	"-" + FormatUInt64(math.MaxInt64): "-" + FormatUInt64(math.MaxInt64),
	"18446744073709551616":            "18446744073709551616",
}

var testTableScientificNotation = map[string]string{

	FormatUInt64(math.MaxUint64):            FormatUInt64(math.MaxUint64),
	FormatUInt64(math.MaxInt64 + 1):         FormatUInt64(math.MaxInt64 + 1),
	"1e9":                                   "1000000000",
	"2.41E-3":                               "0.00241",
	"24.2E-4":                               "0.00242",
	"243E-5":                                "0.00243",
	"1e-5":                                  "0.00001",
	"245E3":                                 "245000",
	"1.2345E-1":                             "0.12345",
	"0e5":                                   "0",
	"0e-5":                                  "0",
	"0.e0":                                  "0",
	".0e0":                                  "0",
	"-0":                                    "0",
	"123.456e0":                             "123.456",
	"123.456e2":                             "12345.6",
	"123.456e10":                            "1234560000000",
	"123456789123456789123456789.123456e-2": "1234567891234567891234567.89123456",
	"123456789123456789123456789123456e-2":  "1234567891234567891234567891234.56",
}

var testErrors = []string{
	"",
	"qwert",
	"-",
	".",
	"-.",
	".-",
	"234-.56",
	"234-56",
	"2-",
	"..",
	"2..",
	"..2",
	".5.2",
	"8..2",
	"8.1.",
	"1e",
	"1-e",
	"1e9e",
	"1ee9",
	"1ee",
	"1eE",
	"1e-",
	"1e-.",
	"1e1.2",
	"123.456e1.3",
	"1e-1.2",
	"123.456e-1.3",
	"123.456Easdf",
	"123.456e" + FormatInt64(math.MinInt64),
	"123.456e" + FormatInt64(math.MinInt32),
	"512.99 USD",
	"$99.99",
	"51,850.00",
	"20_000_000.00",
	"$20_000_000.00",
}

func init() {
	for _, s := range testTable {
		s.exact = strconv.FormatFloat(s.float, 'f', 1500, 64)
		if strings.ContainsRune(s.exact, '.') {
			s.exact = strings.TrimRight(s.exact, "0")
			s.exact = strings.TrimRight(s.exact, ".")
		}
	}

	// add negatives
	withNeg := testTable[:]
	for _, s := range testTable {
		if s.float > 0 && s.short != "0" && s.exact != "0" {
			withNeg = append(withNeg, &testEnt{-s.float, "-" + s.short, "-" + s.exact, "-" + s.inexact})
		}
	}
	testTable = withNeg

	for e, s := range testTableScientificNotation {
		if string(e[0]) == "-" || string(e[0]) == "+" || s == "0" {
			continue
		}
		testTableScientificNotation["-"+e] = "-" + s
	}
}

var testNumConverter = &Converter{}
var testOldNumConverter = &OldConverter{}

func TestConverter_ConvertDecimal(t *testing.T) {
	for _, x := range testTable {
		s := x.short
		d, err := testNumConverter.ConvertDecimal(s)
		if err != nil {
			t.Errorf("error while parsing %s", s)
		} else if d.String() != s {
			t.Errorf("%s expected %s, got %s",
				s, s, d.String())
		}

		if f, _ := d.(DecimalNumber).Float64(); f != x.float {
			t.Errorf("%s expected %v, got %v", s, x.float, f)
		}
	}

	for _, x := range testTableDecimalStr {
		s := x.short
		d, err := testNumConverter.ConvertDecimal(s)
		if err != nil {
			t.Errorf("error while parsing %s", s)
		} else if d.(*DecimalStr).value != x.exact {
			t.Errorf("%s expected value %s(%v), got %s(%v)",
				s, x.exact, len(x.exact), d.(*DecimalStr).value, len(x.exact))
		} else if d.(*DecimalStr).intLen != x.intLen {
			t.Errorf("%s expected intLen %v, got %v",
				s, x.intLen, d.(*DecimalStr).intLen)
		}
	}

	for e, s := range testTableInt64 {
		d, err := testNumConverter.ConvertDecimal(e)
		if err != nil {
			t.Errorf("error while parsing %s", e)
		} else if d.String() != s {
			t.Errorf("%s expected %s, got %s",
				e, s, d.String())
		}
	}
	for e, s := range testTableScientificNotation {
		d, err := testNumConverter.ConvertDecimal(e)
		if err != nil {
			t.Errorf("error while parsing %s", e)
		} else if d.String() != s {
			t.Errorf("%s expected %s, got %s",
				e, s, d.String())
		}
	}
}

func TestOldConverter_ConvertDecimal(t *testing.T) {
	for _, x := range testTable {
		s := x.short
		d, err := testOldNumConverter.ConvertDecimal(s)
		if err != nil {
			t.Errorf("error while parsing %s", s)
		} else if d.String() != s {
			t.Errorf("%s expected %s, got %s",
				s, s, d.String())
		}

		if f, _ := d.(DecimalNumber).Float64(); f != x.float {
			t.Errorf("%s expected %v, got %v", s, x.float, f)
		}
	}

	for _, x := range testTableDecimalStr {
		s := x.short
		d, err := testOldNumConverter.ConvertDecimal(s)
		if err != nil {
			t.Errorf("error while parsing %s", s)
		} else if d.String() != x.exact {
			t.Errorf("%s expected value %s(%v), got %+v(%v)",
				s, x.exact, len(x.exact), *(d.(*Decimal).value), len(x.exact))
		}
	}

	for e, s := range testTableInt64 {
		d, err := testOldNumConverter.ConvertDecimal(e)
		if err != nil {
			t.Errorf("error while parsing %s", e)
		} else if d.String() != s {
			t.Errorf("%s expected %s, got %s",
				e, s, d.String())
		}
	}
	for e, s := range testTableScientificNotation {
		d, err := testOldNumConverter.ConvertDecimal(e)
		if err != nil {
			t.Errorf("error while parsing %s", e)
		} else if d.String() != s {
			t.Errorf("%s expected %s, got %s",
				e, s, d.String())
		}
	}
}
func TestConverter_ConvertNumberErrs(t *testing.T) {
	for _, s := range testErrors {
		_, err := testNumConverter.ConvertDecimal(s)

		if err == nil {
			t.Errorf("error expected when parsing %s", s)
		}
	}
}

func TestOldConverter_ConvertNumberErrs(t *testing.T) {
	for _, s := range testErrors {
		_, err := testOldNumConverter.ConvertDecimal(s)

		if err == nil {
			t.Errorf("error expected when parsing %s", s)
		}
	}
}

func TestOldConverter_ConvertBigInt(t *testing.T) {
	for e, s := range testTableInt64 {
		d, err := testOldNumConverter.ConvertBigInt(e)
		if err != nil {
			t.Errorf("error while parsing %s", e)
		} else if d.String() != s {
			t.Errorf("%s expected %s, got %s",
				e, s, d.String())
		}
	}
}

func TestConverter_ConvertBigInt(t *testing.T) {
	for e, s := range testTableInt64 {
		d, err := testNumConverter.ConvertBigInt(e)
		if err != nil {
			t.Errorf("error while parsing %s", e)
		} else if d.String() != s {
			t.Errorf("%s expected %s, got %s",
				e, s, d.String())
		}
	}
}

func TestConverter_ConvertBigInt_Errs(t *testing.T) {
	for _, s := range testErrors {
		_, err := testNumConverter.ConvertBigInt(s)

		if err == nil {
			t.Errorf("error expected when parsing %s", s)
		}
	}
}

func TestOldConverter_ConvertBigInt_Errs(t *testing.T) {
	for _, s := range testErrors {
		_, err := testOldNumConverter.ConvertBigInt(s)

		if err == nil {
			t.Errorf("error expected when parsing %s", s)
		}
	}
}

func TestOldConverter_ConvertBigIntFromInt(t *testing.T) {
	type args struct {
		i int64
	}
	tests := []struct {
		name    string
		args    args
		wantNum BigIntNumber
	}{
		{
			name: "1",
			args: args{
				i: math.MaxInt64,
			},
			wantNum: &BigInt{
				value: apd.NewBigInt(math.MaxInt64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNum := testOldNumConverter.ConvertBigIntFromInt(tt.args.i); !reflect.DeepEqual(gotNum, tt.wantNum) {
				t.Errorf("OldConverter.ConvertBigIntFromInt() = %v, want %v", gotNum, tt.wantNum)
			}
		})
	}
}

func TestOldConverter_ConvertDecimalFromFloat(t *testing.T) {
	type args struct {
		f float64
	}
	tests := []struct {
		name    string
		args    args
		wantNum DecimalNumber
	}{
		{
			name: "1",
			args: args{
				f: math.MaxFloat64,
			},
			wantNum: &Decimal{
				value: NewApdDecimalFromFloat(math.MaxFloat64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNum := testOldNumConverter.ConvertDecimalFromFloat(tt.args.f); !reflect.DeepEqual(gotNum, tt.wantNum) {
				t.Errorf("OldConverter.ConvertDecimalFromFloat() = %v, want %v", gotNum, tt.wantNum)
			}
		})
	}
}

func TestConverter_ConvertBigIntFromInt(t *testing.T) {
	type args struct {
		i int64
	}
	tests := []struct {
		name    string
		args    args
		wantNum BigIntNumber
	}{
		{
			name: "1",
			args: args{
				i: math.MaxInt64,
			},
			wantNum: &Int64{
				value: math.MaxInt64,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNum := testNumConverter.ConvertBigIntFromInt(tt.args.i); !reflect.DeepEqual(gotNum, tt.wantNum) {
				t.Errorf("Converter.ConvertBigIntFromInt() = %v, want %v", gotNum, tt.wantNum)
			}
		})
	}
}

func TestConverter_ConvertDecimalFromFloat(t *testing.T) {
	type args struct {
		f float64
	}
	tests := []struct {
		name    string
		args    args
		wantNum DecimalNumber
	}{
		{
			name: "1",
			args: args{
				f: math.MaxFloat64,
			},
			wantNum: &Float64{
				value: math.MaxFloat64,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNum := testNumConverter.ConvertDecimalFromFloat(tt.args.f); !reflect.DeepEqual(gotNum, tt.wantNum) {
				t.Errorf("Converter.ConvertDecimalFromFloat() = %v, want %v", gotNum, tt.wantNum)
			}
		})
	}
}

func TestOldConverter_ConvertBigIntFromUint(t *testing.T) {
	type args struct {
		i uint64
	}
	tests := []struct {
		name    string
		c       *OldConverter
		args    args
		wantNum BigIntNumber
	}{
		{
			name: "1",
			c:    &OldConverter{},
			args: args{
				i: math.MaxUint64,
			},
			wantNum: &BigInt{
				value: new(apd.BigInt).SetUint64(math.MaxUint64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNum := tt.c.ConvertBigIntFromUint(tt.args.i); !reflect.DeepEqual(gotNum, tt.wantNum) {
				t.Errorf("OldConverter.ConvertBigIntFromUint() = %v, want %v", gotNum, tt.wantNum)
			}
		})
	}
}

func TestOldConverter_ConvertDecimalFromFloat32(t *testing.T) {
	type args struct {
		f float32
	}
	tests := []struct {
		name    string
		c       *OldConverter
		args    args
		wantNum DecimalNumber
	}{
		{
			name: "1",
			c:    &OldConverter{},
			args: args{
				f: math.MaxFloat32,
			},
			wantNum: &Decimal{
				value: NewApdDecimalFromFloat32(math.MaxFloat32),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNum := tt.c.ConvertDecimalFromFloat32(tt.args.f); !reflect.DeepEqual(gotNum, tt.wantNum) {
				t.Errorf("OldConverter.ConvertDecimalFromFloat32() = %v, want %v", gotNum, tt.wantNum)
			}
		})
	}
}

func TestConverter_ConvertBigIntFromUint(t *testing.T) {
	type args struct {
		i uint64
	}
	tests := []struct {
		name    string
		c       *Converter
		args    args
		wantNum BigIntNumber
	}{
		{
			name: "1",
			c:    &Converter{},
			args: args{
				i: math.MaxUint64,
			},
			wantNum: &Uint64{
				value: math.MaxUint64,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNum := tt.c.ConvertBigIntFromUint(tt.args.i); !reflect.DeepEqual(gotNum, tt.wantNum) {
				t.Errorf("Converter.ConvertBigIntFromUint() = %v, want %v", gotNum, tt.wantNum)
			}
		})
	}
}

func TestConverter_ConvertDecimalFromFloat32(t *testing.T) {
	type args struct {
		f float32
	}
	tests := []struct {
		name    string
		c       *Converter
		args    args
		wantNum DecimalNumber
	}{
		{
			name: "1",
			c:    &Converter{},
			args: args{
				f: math.MaxFloat32,
			},
			wantNum: &Decimal{
				value: NewApdDecimalFromFloat32(math.MaxFloat32),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNum := tt.c.ConvertDecimalFromFloat32(tt.args.f); !reflect.DeepEqual(gotNum, tt.wantNum) {
				t.Errorf("Converter.ConvertDecimalFromFloat32() = %v, want %v", gotNum, tt.wantNum)
			}
		})
	}
}
