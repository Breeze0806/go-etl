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
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

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
	{math.MaxInt64, strconv.FormatFloat(float64(math.MaxInt64), 'f', -1, 64), "", strconv.FormatInt(math.MaxInt64, 10)},
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
	"+012345678901234567890":                              "12345678901234567890",
	"+0000000103456789123456789012":                       "103456789123456789012",
	"0000000103456789123456789012":                        "103456789123456789012",
	"+103456789123456789012":                              "103456789123456789012",
	"-103456789123456789012":                              "-103456789123456789012",
	"103456789123456789012":                               "103456789123456789012",
	"-0000000103456789123456789012":                       "-103456789123456789012",
}

var testTableScientificNotation = map[string]string{

	strconv.FormatUint(math.MaxUint64, 10):  strconv.FormatUint(math.MaxUint64, 10),
	strconv.FormatUint(math.MaxInt64+1, 10): strconv.FormatUint(math.MaxInt64+1, 10),
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
	"123.456e0":                             "123.456",
	"123.456e2":                             "12345.6",
	"123.456e10":                            "1234560000000",
	"123456789123456789123456789.123456e-2": "1234567891234567891234567.89123456",
	"123456789123456789123456789123456e-2":  "1234567891234567891234567891234.56",
}

var testNumConverter = &Converter{}
var testOldNumConverter = &OldConverter{}

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
			t.Errorf("%s expected value %s(%v), got %s(%v)",
				s, x.exact, len(x.exact), d.(*DecimalStr).value, len(x.exact))
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
	tests := []string{
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
		"123.456e" + strconv.FormatInt(math.MinInt64, 10),
		"123.456e" + strconv.FormatInt(math.MinInt32, 10),
		"512.99 USD",
		"$99.99",
		"51,850.00",
		"20_000_000.00",
		"$20_000_000.00",
	}

	for _, s := range tests {
		_, err := testNumConverter.ConvertDecimal(s)

		if err == nil {
			t.Errorf("error expected when parsing %s", s)
		}
	}
}

func TestOldConverter_ConvertNumberErrs(t *testing.T) {
	tests := []string{
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
		"123.456e" + strconv.FormatInt(math.MinInt64, 10),
		"123.456e" + strconv.FormatInt(math.MinInt32, 10),
		"512.99 USD",
		"$99.99",
		"51,850.00",
		"20_000_000.00",
		"$20_000_000.00",
	}

	for _, s := range tests {
		_, err := testOldNumConverter.ConvertDecimal(s)

		if err == nil {
			t.Errorf("error expected when parsing %s", s)
		}
	}
}
func TestInt64_Bool(t *testing.T) {
	tests := []struct {
		name    string
		i       *Int64
		want    bool
		wantErr bool
	}{
		{
			name: "1",
			i: &Int64{
				value: 0,
			},
			want: false,
		},
		{
			name: "2",
			i: &Int64{
				value: math.MaxInt64,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Bool()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64.Bool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Int64.Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_Int64(t *testing.T) {
	tests := []struct {
		name    string
		i       *Int64
		want    int64
		wantErr bool
	}{
		{
			name: "1",
			i: &Int64{
				value: 0,
			},
			want: 0,
		},
		{
			name: "2",
			i: &Int64{
				value: math.MaxInt64,
			},
			want: math.MaxInt64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Int64()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64.Int64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Int64.Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_Float64(t *testing.T) {
	tests := []struct {
		name    string
		i       *Int64
		want    float64
		wantErr bool
	}{
		{
			name: "1",
			i: &Int64{
				value: 0,
			},
			want: 0.0,
		},
		{
			name: "2",
			i: &Int64{
				value: math.MaxInt64,
			},
			want: 9.223372036854776e+18,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Float64()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64.Float64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Int64.Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_BigInt(t *testing.T) {
	tests := []struct {
		name string
		i    *Int64
		want BigIntNumber
	}{
		{
			name: "1",
			i: &Int64{
				value: math.MaxInt64,
			},
			want: &Int64{
				value: math.MaxInt64,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.BigInt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64.BigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_Decimal(t *testing.T) {
	tests := []struct {
		name string
		i    *Int64
		want DecimalNumber
	}{
		{
			name: "1",
			i: &Int64{
				value: math.MaxInt64,
			},
			want: &Int64{
				value: math.MaxInt64,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Decimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64.Decimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_CloneBigInt(t *testing.T) {
	tests := []struct {
		name string
		i    *Int64
		want BigIntNumber
	}{
		{
			name: "1",
			i: &Int64{
				value: math.MaxInt64,
			},
			want: &Int64{
				value: math.MaxInt64,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.CloneBigInt(); !reflect.DeepEqual(got, tt.want) || got == tt.i {
				t.Errorf("Int64.CloneBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_CloneDecimal(t *testing.T) {
	tests := []struct {
		name string
		i    *Int64
		want DecimalNumber
	}{
		{
			name: "1",
			i: &Int64{
				value: math.MaxInt64,
			},
			want: &Int64{
				value: math.MaxInt64,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.CloneDecimal(); !reflect.DeepEqual(got, tt.want) || got == tt.i {
				t.Errorf("Int64.CloneDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_AsBigInt(t *testing.T) {
	tests := []struct {
		name string
		i    *Int64
		want *big.Int
	}{
		{
			name: "1",
			i: &Int64{
				value: math.MaxInt64,
			},
			want: big.NewInt(math.MaxInt64),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.AsBigInt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64.AsBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_AsDecimal(t *testing.T) {
	tests := []struct {
		name string
		i    *Int64
		want decimal.Decimal
	}{
		{
			name: "1",
			i: &Int64{
				value: math.MaxInt64,
			},
			want: testDecimalFormString(strconv.FormatInt(math.MaxInt64, 10)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.AsDecimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigInt_Bool(t *testing.T) {
	tests := []struct {
		name    string
		b       *BigInt
		want    bool
		wantErr bool
	}{
		{
			name: "1",
			b: &BigInt{
				value: big.NewInt(math.MaxInt64),
			},
			want: true,
		},
		{
			name: "2",
			b: &BigInt{
				value: big.NewInt(0),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Bool()
			if (err != nil) != tt.wantErr {
				t.Errorf("BigInt.Bool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BigInt.Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigInt_Int64(t *testing.T) {
	tests := []struct {
		name    string
		b       *BigInt
		want    int64
		wantErr bool
	}{
		{
			name: "1",
			b: &BigInt{
				value: big.NewInt(math.MaxInt64),
			},
			want: math.MaxInt64,
		},
		{
			name: "2",
			b: &BigInt{
				value: big.NewInt(0),
			},
			want: 0,
		},
		{
			name: "3",
			b: &BigInt{
				value: testBigIntFromString(testDecimalFormString("1e1000").String()),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Int64()
			if (err != nil) != tt.wantErr {
				t.Errorf("BigInt.Int64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BigInt.Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigInt_Float64(t *testing.T) {
	tests := []struct {
		name    string
		b       *BigInt
		wantV   float64
		wantErr bool
	}{
		{
			name: "1",
			b: &BigInt{
				value: big.NewInt(math.MaxInt64),
			},
			wantV: 9.223372036854776e+18,
		},
		{
			name: "2",
			b: &BigInt{
				value: big.NewInt(0),
			},
			wantV: 0.0,
		},
		{
			name: "3",
			b: &BigInt{
				value: testBigIntFromString(testDecimalFormString("1e1000").String()),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotV, err := tt.b.Float64()
			if (err != nil) != tt.wantErr {
				t.Errorf("BigInt.Float64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotV != tt.wantV {
				t.Errorf("BigInt.Float64() = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

func TestBigInt_BigInt(t *testing.T) {
	tests := []struct {
		name string
		b    *BigInt
		want BigIntNumber
	}{
		{
			name: "1",
			b: &BigInt{
				value: big.NewInt(math.MaxInt64),
			},
			want: &BigInt{
				value: big.NewInt(math.MaxInt64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.BigInt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BigInt.BigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigInt_Decimal(t *testing.T) {
	tests := []struct {
		name string
		b    *BigInt
		want DecimalNumber
	}{
		{
			name: "1",
			b: &BigInt{
				value: big.NewInt(math.MaxInt64),
			},
			want: &BigInt{
				value: big.NewInt(math.MaxInt64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Decimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BigInt.Decimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigInt_CloneDecimal(t *testing.T) {
	tests := []struct {
		name string
		b    *BigInt
		want DecimalNumber
	}{
		{
			name: "1",
			b: &BigInt{
				value: big.NewInt(math.MaxInt64),
			},
			want: &BigInt{
				value: big.NewInt(math.MaxInt64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.CloneDecimal(); !reflect.DeepEqual(got, tt.want) || got == tt.b {
				t.Errorf("BigInt.CloneDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigInt_String(t *testing.T) {
	tests := []struct {
		name string
		b    *BigInt
		want string
	}{
		{
			name: "1",
			b: &BigInt{
				value: big.NewInt(math.MaxInt64),
			},
			want: strconv.FormatInt(math.MaxInt64, 10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("BigInt.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigInt_CloneBigInt(t *testing.T) {
	tests := []struct {
		name string
		b    *BigInt
		want BigIntNumber
	}{
		{
			name: "1",
			b: &BigInt{
				value: big.NewInt(math.MaxInt64),
			},
			want: &BigInt{
				value: big.NewInt(math.MaxInt64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.CloneBigInt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BigInt.CloneBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntStr_Bool(t *testing.T) {
	tests := []struct {
		name    string
		b       *BigIntStr
		want    bool
		wantErr bool
	}{
		{
			name: "1",
			b: &BigIntStr{
				value: "1234567890123456789",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Bool()
			if (err != nil) != tt.wantErr {
				t.Errorf("BigIntStr.Bool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BigIntStr.Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntStr_Int64(t *testing.T) {
	tests := []struct {
		name    string
		b       *BigIntStr
		want    int64
		wantErr bool
	}{
		{
			name: "1",
			b: &BigIntStr{
				value: "1234567890123456789",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Int64()
			if (err != nil) != tt.wantErr {
				t.Errorf("BigIntStr.Int64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BigIntStr.Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntStr_Float64(t *testing.T) {
	tests := []struct {
		name    string
		b       *BigIntStr
		wantV   float64
		wantErr bool
	}{
		{
			name: "1",
			b: &BigIntStr{
				value: testDecimalFormString("1e100").String(),
			},
			wantV: 1e100,
		},

		{
			name: "3",
			b: &BigIntStr{
				value: testDecimalFormString("1e1000").String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotV, err := tt.b.Float64()
			if (err != nil) != tt.wantErr {
				t.Errorf("BigIntStr.Float64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotV != tt.wantV {
				t.Errorf("BigIntStr.Float64() = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

func TestBigIntStr_BigInt(t *testing.T) {
	tests := []struct {
		name string
		b    *BigIntStr
		want BigIntNumber
	}{
		{
			name: "1",
			b: &BigIntStr{
				value: new(big.Int).SetUint64(math.MaxUint16).String(),
			},
			want: &BigIntStr{
				value: new(big.Int).SetUint64(math.MaxUint16).String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.BigInt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BigIntStr.BigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntStr_Decimal(t *testing.T) {
	tests := []struct {
		name string
		b    *BigIntStr
		want DecimalNumber
	}{
		{
			name: "1",
			b: &BigIntStr{
				value: new(big.Int).SetUint64(math.MaxUint16).String(),
			},
			want: &BigIntStr{
				value: new(big.Int).SetUint64(math.MaxUint16).String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Decimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BigIntStr.Decimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntStr_CloneBigInt(t *testing.T) {
	tests := []struct {
		name string
		b    *BigIntStr
		want BigIntNumber
	}{
		{
			name: "1",
			b: &BigIntStr{
				value: new(big.Int).SetUint64(math.MaxUint16).String(),
			},
			want: &BigIntStr{
				value: new(big.Int).SetUint64(math.MaxUint16).String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.CloneBigInt(); !reflect.DeepEqual(got, tt.want) || got == tt.b {
				t.Errorf("BigIntStr.CloneBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntStr_CloneDecimal(t *testing.T) {
	tests := []struct {
		name string
		b    *BigIntStr
		want DecimalNumber
	}{
		{
			name: "1",
			b: &BigIntStr{
				value: new(big.Int).SetUint64(math.MaxUint16).String(),
			},
			want: &BigIntStr{
				value: new(big.Int).SetUint64(math.MaxUint16).String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.CloneDecimal(); !reflect.DeepEqual(got, tt.want) || got == tt.b {
				t.Errorf("BigIntStr.CloneDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntStr_AsBigInt(t *testing.T) {
	tests := []struct {
		name string
		b    *BigIntStr
		want *big.Int
	}{
		{
			name: "1",
			b: &BigIntStr{
				value: new(big.Int).SetUint64(math.MaxUint16).String(),
			},
			want: new(big.Int).SetUint64(math.MaxUint16),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.AsBigInt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BigIntStr.AsBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntStr_AsDecimal(t *testing.T) {
	tests := []struct {
		name string
		b    *BigIntStr
		want decimal.Decimal
	}{
		{
			name: "1",
			b: &BigIntStr{
				value: new(big.Int).SetUint64(math.MaxUint16).String(),
			},
			want: testDecimalFormString(new(big.Int).SetUint64(math.MaxUint16).String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.AsDecimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BigIntStr.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalStr_Bool(t *testing.T) {
	tests := []struct {
		name    string
		d       *DecimalStr
		want    bool
		wantErr bool
	}{
		{
			name: "1",
			d: &DecimalStr{
				value:  "123456.0789",
				intLen: 6,
			},
			want: true,
		},
		{
			name: "2",
			d: &DecimalStr{
				value:  "0",
				intLen: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Bool()
			if (err != nil) != tt.wantErr {
				t.Errorf("DecimalStr.Bool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecimalStr.Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalStr_Float64(t *testing.T) {
	tests := []struct {
		name    string
		d       *DecimalStr
		want    float64
		wantErr bool
	}{
		{
			name: "1",
			d: &DecimalStr{
				value:  testDecimalFormString("1e1000").String(),
				intLen: 1000,
			},
			wantErr: true,
		},
		{
			name: "2",
			d: &DecimalStr{
				value:  "123456.0789",
				intLen: 6,
			},
			want: 123456.0789,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Float64()
			if (err != nil) != tt.wantErr {
				t.Errorf("DecimalStr.Float64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecimalStr.Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalStr_BigInt(t *testing.T) {
	tests := []struct {
		name string
		d    *DecimalStr
		want BigIntNumber
	}{
		{
			name: "1",
			d: &DecimalStr{
				value:  testDecimalFormString("1.797693134862315708145274237317043567981e+308").String() + ".123",
				intLen: 309,
			},
			want: &BigIntStr{
				value: testDecimalFormString("1.797693134862315708145274237317043567981e+308").String(),
			},
		},
		{
			name: "2",
			d: &DecimalStr{
				value:  strconv.FormatInt(math.MaxInt64, 10),
				intLen: len(strconv.FormatInt(math.MaxInt64, 10)),
			},
			want: &Int64{
				value: math.MaxInt64,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.BigInt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecimalStr.BigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalStr_Decimal(t *testing.T) {
	tests := []struct {
		name string
		d    *DecimalStr
		want DecimalNumber
	}{
		{
			name: "1",
			d: &DecimalStr{
				value:  testDecimalFormString("1.797693134862315708145274237317043567981e+308").String() + ".123",
				intLen: 309,
			},
			want: &DecimalStr{
				value:  testDecimalFormString("1.797693134862315708145274237317043567981e+308").String() + ".123",
				intLen: 309,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Decimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecimalStr.Decimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalStr_CloneDecimal(t *testing.T) {
	tests := []struct {
		name string
		d    *DecimalStr
		want DecimalNumber
	}{
		{
			name: "1",
			d: &DecimalStr{
				value:  testDecimalFormString("1.797693134862315708145274237317043567981e+308").String() + ".123",
				intLen: 309,
			},
			want: &DecimalStr{
				value:  testDecimalFormString("1.797693134862315708145274237317043567981e+308").String() + ".123",
				intLen: 309,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.CloneDecimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecimalStr.CloneDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalStr_AsDecimal(t *testing.T) {
	tests := []struct {
		name string
		d    *DecimalStr
		want decimal.Decimal
	}{
		{
			name: "1",
			d: &DecimalStr{
				value:  testDecimalFormString("1.797693134862315708145274237317043567981e+308").String() + ".123",
				intLen: 309,
			},
			want: testDecimalFormString(testDecimalFormString("1.797693134862315708145274237317043567981e+308").String() + ".123"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.AsDecimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecimalStr.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_Bool(t *testing.T) {
	tests := []struct {
		name    string
		d       *Decimal
		want    bool
		wantErr bool
	}{
		{
			name: "1",
			d: &Decimal{
				value: decimal.Zero,
			},
			want: false,
		},
		{
			name: "2",
			d: &Decimal{
				value: decimal.NewFromFloat(1e32),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Bool()
			if (err != nil) != tt.wantErr {
				t.Errorf("Decimal.Bool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decimal.Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_Float64(t *testing.T) {
	tests := []struct {
		name    string
		d       *Decimal
		want    float64
		wantErr bool
	}{
		{
			name: "1",
			d: &Decimal{
				value: testDecimalFormString("1e1000"),
			},
			wantErr: true,
		},
		{
			name: "2",
			d: &Decimal{
				value: decimal.NewFromFloat(math.MaxFloat64),
			},
			want: math.MaxFloat64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Float64()
			if (err != nil) != tt.wantErr {
				t.Errorf("Decimal.Float64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decimal.Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_BigInt(t *testing.T) {
	tests := []struct {
		name string
		d    *Decimal
		want BigIntNumber
	}{
		{
			name: "1",
			d: &Decimal{
				value: testDecimalFormString("123456232542542525.525254252524"),
			},
			want: &BigInt{
				value: big.NewInt(123456232542542525),
			},
		},
		{
			name: "2",
			d: &Decimal{
				value: testDecimalFormString("-123456232542542525.525254252524"),
			},
			want: &BigInt{
				value: big.NewInt(-123456232542542525),
			},
		},
		{
			name: "3",
			d: &Decimal{
				value: testDecimalFormString("0.00122323123"),
			},
			want: &BigInt{
				value: big.NewInt(0),
			},
		},
		{
			name: "4",
			d: &Decimal{
				value: testDecimalFormString("123450000"),
			},
			want: &BigInt{
				value: big.NewInt(123450000),
			},
		},

		{
			name: "5",
			d: &Decimal{
				value: testDecimalFormString("1.00122323123"),
			},
			want: &BigInt{
				value: big.NewInt(1),
			},
		},
		{
			name: "6",
			d: &Decimal{
				value: testDecimalFormString("12345.67e-3"),
			},
			want: &BigInt{
				value: big.NewInt(12),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.BigInt(); got.String() != tt.want.String() {
				t.Errorf("Decimal.BigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_Decimal(t *testing.T) {
	tests := []struct {
		name string
		d    *Decimal
		want DecimalNumber
	}{
		{
			name: "1",
			d: &Decimal{
				value: testDecimalFormString("123456232542542525.525254252524"),
			},
			want: &Decimal{
				value: testDecimalFormString("123456232542542525.525254252524"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Decimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decimal.Decimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_CloneDecimal(t *testing.T) {
	tests := []struct {
		name string
		d    *Decimal
		want DecimalNumber
	}{
		{
			name: "1",
			d: &Decimal{
				value: testDecimalFormString("123456232542542525.525254252524"),
			},
			want: &Decimal{
				value: testDecimalFormString("123456232542542525.525254252524"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.CloneDecimal(); !reflect.DeepEqual(got, tt.want) || got == tt.want {
				t.Errorf("Decimal.CloneDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_AsDecimal(t *testing.T) {
	tests := []struct {
		name string
		d    *Decimal
		want decimal.Decimal
	}{
		{
			name: "1",
			d: &Decimal{
				value: testDecimalFormString("123456232542542525.525254252524"),
			},
			want: testDecimalFormString("123456232542542525.525254252524"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.AsDecimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decimal.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_String(t *testing.T) {
	tests := []struct {
		name string
		d    *Decimal
		want string
	}{
		{
			name: "1",
			d: &Decimal{
				value: testDecimalFormString("123456232542542525.525254252524"),
			},
			want: "123456232542542525.525254252524",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("Decimal.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigInt_AsDecimal(t *testing.T) {
	tests := []struct {
		name string
		b    *BigInt
		want decimal.Decimal
	}{
		{
			name: "1",
			b: &BigInt{
				value: big.NewInt(math.MaxInt64),
			},
			want: decimal.NewFromInt(math.MaxInt64),
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.AsDecimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BigInt.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
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
	tests := []string{
		"",
		"qwert",
		"-",
		".",
		"-.",
		".-",
		"234-.56",
		"234-56",
		"2-",
		"2.",
		".2",
		".5.2",
		"8.2",
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
		"123.456e" + strconv.FormatInt(math.MinInt64, 10),
		"123.456e" + strconv.FormatInt(math.MaxUint16, 10),
		"512.99 USD",
		"$99.99",
		"51,850.00",
		"20_000_000.00",
		"$20_000_000.00",
	}

	for _, s := range tests {
		_, err := testNumConverter.ConvertBigInt(s)

		if err == nil {
			t.Errorf("error expected when parsing %s", s)
		}
	}
}

func TestOldConverter_ConvertBigInt_Errs(t *testing.T) {
	tests := []string{
		"",
		"qwert",
		"-",
		".",
		"-.",
		".-",
		"234-.56",
		"234-56",
		"2-",
		"2.",
		".2",
		".5.2",
		"8.2",
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
		"123.456e" + strconv.FormatInt(math.MinInt64, 10),
		"123.456e" + strconv.FormatInt(math.MaxUint16, 10),
		"512.99 USD",
		"$99.99",
		"51,850.00",
		"20_000_000.00",
		"$20_000_000.00",
	}

	for _, s := range tests {
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
				value: big.NewInt(math.MaxInt64),
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
				value: decimal.NewFromFloat(math.MaxFloat64),
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
			wantNum: &Decimal{
				value: decimal.NewFromFloat(math.MaxFloat64),
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

func TestUint64_Bool(t *testing.T) {
	tests := []struct {
		name    string
		i       *Uint64
		want    bool
		wantErr bool
	}{
		{
			name:    "1",
			i:       &Uint64{value: 0},
			want:    false,
			wantErr: false,
		},
		{
			name:    "2",
			i:       &Uint64{value: math.MaxUint64},
			want:    true,
			wantErr: false,
		},
		{
			name:    "3",
			i:       &Uint64{value: 123456789},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Bool()
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint64.Bool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Uint64.Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_Int64(t *testing.T) {
	tests := []struct {
		name    string
		i       *Uint64
		want    int64
		wantErr bool
	}{
		{
			name:    "1",
			i:       &Uint64{value: 0},
			want:    0,
			wantErr: false,
		},
		{
			name:    "2",
			i:       &Uint64{value: uint64(math.MaxInt64)},
			want:    math.MaxInt64,
			wantErr: false,
		},
		{
			name:    "3",
			i:       &Uint64{value: uint64(math.MaxInt64 + 1)},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Int64()
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint64.Int64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Uint64.Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_Float64(t *testing.T) {
	tests := []struct {
		name    string
		i       *Uint64
		want    float64
		wantErr bool
	}{
		{
			name:    "1",
			i:       &Uint64{value: 0},
			want:    0.0,
			wantErr: false,
		},
		{
			name:    "2",
			i:       &Uint64{value: math.MaxInt64},
			want:    9.223372036854776e+18,
			wantErr: false,
		},
		{
			name:    "3",
			i:       &Uint64{value: math.MaxUint64},
			want:    1.8446744073709552e+19,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Float64()
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint64.Float64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Uint64.Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_BigInt(t *testing.T) {
	tests := []struct {
		name string
		i    *Uint64
		want BigIntNumber
	}{
		{
			name: "1",
			i:    &Uint64{value: math.MaxUint64},
			want: &Uint64{value: math.MaxUint64},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.BigInt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64.BigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_Decimal(t *testing.T) {
	tests := []struct {
		name string
		i    *Uint64
		want DecimalNumber
	}{
		{
			name: "1",
			i:    &Uint64{value: math.MaxUint64},
			want: &Uint64{value: math.MaxUint64},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Decimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64.Decimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_String(t *testing.T) {
	tests := []struct {
		name string
		i    *Uint64
		want string
	}{
		{
			name: "1",
			i:    &Uint64{value: 1234567890},
			want: "1234567890",
		},
		{
			name: "2",
			i:    &Uint64{value: math.MaxUint64},
			want: "18446744073709551615",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.String(); got != tt.want {
				t.Errorf("Uint64.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_CloneBigInt(t *testing.T) {
	tests := []struct {
		name string
		i    *Uint64
		want BigIntNumber
	}{
		{
			name: "1",
			i:    &Uint64{value: math.MaxUint64},
			want: &Uint64{value: math.MaxUint64},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.CloneBigInt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64.CloneBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_CloneDecimal(t *testing.T) {
	tests := []struct {
		name string
		i    *Uint64
		want DecimalNumber
	}{
		{
			name: "1",
			i:    &Uint64{value: math.MaxUint64},
			want: &Uint64{value: math.MaxUint64},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.CloneDecimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64.CloneDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64_AsBigInt(t *testing.T) {
	tests := []struct {
		name   string
		i      *Uint64
		wantBi *big.Int
	}{
		{
			name:   "1",
			i:      &Uint64{value: math.MaxInt64 + 1},
			wantBi: new(big.Int).SetUint64(math.MaxInt64 + 1),
		},
		{
			name:   "2",
			i:      &Uint64{value: math.MaxUint64},
			wantBi: new(big.Int).SetUint64(math.MaxUint64),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotBi := tt.i.AsBigInt(); !reflect.DeepEqual(gotBi, tt.wantBi) {
				t.Errorf("Uint64.AsBigInt() = %v, want %v", gotBi, tt.wantBi)
			}
		})
	}
}

func TestUint64_AsDecimal(t *testing.T) {
	tests := []struct {
		name string
		i    *Uint64
		want decimal.Decimal
	}{
		{
			name: "1",
			i:    &Uint64{value: math.MaxInt64 + 1},
			want: decimal.NewFromUint64(math.MaxInt64 + 1),
		},
		{
			name: "2",
			i:    &Uint64{value: math.MaxUint64},
			want: decimal.NewFromUint64(math.MaxUint64),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.AsDecimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64.AsDecimal() = %v, want %v", got, tt.want)
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
				value: new(big.Int).SetUint64(math.MaxUint64),
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
				value: decimal.NewFromFloat32(math.MaxFloat32),
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
				value: decimal.NewFromFloat32(math.MaxFloat32),
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
