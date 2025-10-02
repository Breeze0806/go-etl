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
	"math/rand"
	"strconv"
	"testing"
	"testing/quick"

	"github.com/cockroachdb/apd/v3"
)

type Inp struct {
	float float64
	exp   int32
}

func TestNewFromFloat(t *testing.T) {
	for _, x := range testTable {
		s := x.short
		d := NewFromFloat(x.float)
		dn, _ := convertDecimal(d.Text('f'))
		if dn.String() != s {
			t.Errorf("expected %s, got %s (float: %v) (%s, %d)",
				s, d.String(), x.float,
				d.Coeff.String(), d.Exponent)
		}
	}

	shouldPanicOn := []float64{
		math.NaN(),
		math.Inf(1),
		math.Inf(-1),
	}

	for _, n := range shouldPanicOn {
		var d *apd.Decimal
		if !didPanic(func() { d = NewFromFloat(n) }) {
			t.Fatalf("Expected panic when creating a Decimal from %v, got %v instead", n, d.String())
		}
	}
}

func TestNewFromFloatRandom(t *testing.T) {
	n := 0
	rng := rand.New(rand.NewSource(0xdead1337))
	for {
		n++
		if n == 10 {
			break
		}
		in := (rng.Float64() - 0.5) * math.MaxFloat64 * 2
		want, _, err := apd.NewFromString(strconv.FormatFloat(in, 'f', -1, 64))
		if err != nil {
			t.Error(err)
			continue
		}
		got := NewFromFloat(in)
		if want.Cmp(got) != 0 {
			t.Errorf("in: %v, expected %s (%s, %d), got %s (%s, %d) ",
				in, want.String(), want.Coeff.String(), want.Exponent,
				got.String(), got.Coeff.String(), got.Exponent)
		}
	}
}

func TestNewFromFloatQuick(t *testing.T) {
	err := quick.Check(func(f float64) bool {
		want, _, werr := apd.NewFromString(strconv.FormatFloat(f, 'f', -1, 64))
		if werr != nil {
			return true
		}
		got := NewFromFloat(f)
		return got.Cmp(want) == 0
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestNewFromFloat32Random(t *testing.T) {
	n := 0
	rng := rand.New(rand.NewSource(0xdead1337))
	for {
		n++
		if n == 10 {
			break
		}
		in := float32((rng.Float64() - 0.5) * math.MaxFloat32 * 2)
		want, _, err := apd.NewFromString(strconv.FormatFloat(float64(in), 'f', -1, 32))
		if err != nil {
			t.Error(err)
			continue
		}
		got := NewFromFloat32(in)
		if want.Cmp(got) != 0 {
			t.Errorf("in: %v, expected %s (%s, %d), got %s (%s, %d) ",
				in, want.String(), want.Coeff.String(), want.Exponent,
				got.String(), got.Coeff.String(), got.Exponent)
		}
	}
}

func TestNewFromFloat32Quick(t *testing.T) {
	err := quick.Check(func(f float32) bool {
		want, _, werr := apd.NewFromString(strconv.FormatFloat(float64(f), 'f', -1, 32))
		if werr != nil {
			return true
		}
		got := NewFromFloat32(f)
		return got.Cmp(want) == 0
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func didPanic(f func()) bool {
	ret := false
	func() {

		defer func() {
			if message := recover(); message != nil {
				ret = true
			}
		}()

		// call the target function
		f()

	}()

	return ret

}

func TestNewFromFloatWithExponent(t *testing.T) {

	// some tests are taken from here https://www.cockroachlabs.com/blog/rounding-implementations-in-go/
	tests := map[Inp]string{
		{123.4, -3}:                 "123.4",
		{123.4, -1}:                 "123.4",
		{123.412345, 1}:             "120",
		{123.412345, 0}:             "123",
		{123.412345, -5}:            "123.41235",
		{123.412345, -6}:            "123.412345",
		{123.412345, -7}:            "123.412345",
		{123.412345, -28}:           "123.4123450000000019599610823207",
		{1230000000, 3}:             "1230000000",
		{123.9999999999999999, -7}:  "124",
		{123.8989898999999999, -7}:  "123.8989899",
		{0.49999999999999994, 0}:    "0",
		{0.5, 0}:                    "1",
		{0., -1000}:                 "0",
		{0.5000000000000001, 0}:     "1",
		{1.390671161567e-309, 0}:    "0",
		{4.503599627370497e+15, 0}:  "4503599627370497",
		{4.503599627370497e+60, 0}:  "4503599627370497110902645731364739935039854989106233267453952",
		{4.503599627370497e+60, 1}:  "4503599627370497110902645731364739935039854989106233267453950",
		{4.503599627370497e+60, -1}: "4503599627370497110902645731364739935039854989106233267453952",
		{50, 2}:                     "100",
		{49, 2}:                     "0",
		{50, 3}:                     "0",
		// subnormals
		{1.390671161567e-309, -2000}: "0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001390671161567000864431395448332752540137009987788957394095829635554502771758698872408926974382819387852542087331897381878220271350970912568035007740861074263206736245957501456549756342151614772544950978154339064833880234531754156635411349342950306987480369774780312897442981323940546749863054846093718407237782253156822124910364044261653195961209878120072488178603782495270845071470243842997312255994555557251870400944414666445871039673491570643357351279578519863428540219295076767898526278029257129758694673164251056158277568765100904638511604478844087596428177947970563689475826736810456067108202083804368114484417399279328807983736233036662284338182105684628835292230438999173947056675615385756827890872955322265625",
		{1.390671161567e-309, -862}:  "0.0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013906711615670008644313954483327525401370099877889573940958296355545027717586988724089269743828193878525420873318973818782202713509709125680350077408610742632067362459575014565497563421516147725449509781543390648338802345317541566354113493429503069874803697747803128974429813239405467498630548460937184072377822531568221249103640442616531959612098781200724881786037824952708450714702438429973122559945555572518704009444146664458710396734915706433573512795785198634285402192950767678985262780292571297586946731642510561582775687651009046385116044788440876",
		{1.390671161567e-309, -863}:  "0.0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013906711615670008644313954483327525401370099877889573940958296355545027717586988724089269743828193878525420873318973818782202713509709125680350077408610742632067362459575014565497563421516147725449509781543390648338802345317541566354113493429503069874803697747803128974429813239405467498630548460937184072377822531568221249103640442616531959612098781200724881786037824952708450714702438429973122559945555572518704009444146664458710396734915706433573512795785198634285402192950767678985262780292571297586946731642510561582775687651009046385116044788440876",
	}

	// add negatives
	for p, s := range tests {
		if p.float > 0 {
			if s != "0" {
				tests[Inp{-p.float, p.exp}] = "-" + s
			} else {
				tests[Inp{-p.float, p.exp}] = "0"
			}
		}
	}

	for input, s := range tests {
		d := newFromFloatWithExponent(input.float, input.exp)
		dn, _ := convertDecimal(d.Text('f'))
		if dn.String() != s {
			t.Errorf("expected %s, got %s (%s, %d)",
				s, d.String(),
				d.Coeff.String(), d.Exponent)
		}
	}

	shouldPanicOn := []float64{
		math.NaN(),
		math.Inf(1),
		math.Inf(-1),
	}

	for _, n := range shouldPanicOn {
		var d *apd.Decimal
		if !didPanic(func() { d = newFromFloatWithExponent(n, 0) }) {
			t.Fatalf("Expected panic when creating a Decimal from %v, got %v instead", n, d.String())
		}
	}
}

func TestNewFromFloat32(t *testing.T) {
	type args struct {
		value float32
	}
	tests := []struct {
		name string
		args args
		want *apd.Decimal
	}{
		{
			name: "Zero",
			args: args{
				value: 0,
			},
			want: _DecimalZero,
		},
		{
			name: "BigFloat",
			args: args{
				value: math.MaxFloat32,
			},
			want: apd.New(34028235, 31),
		},
		{
			name: "-BigFloat",
			args: args{
				value: -math.MaxFloat32,
			},
			want: apd.New(-34028235, 31),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFromFloat32(tt.args.value); got.Cmp(tt.want) != 0 {
				t.Errorf("NewFromFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_NNewFromFloat(t *testing.T) {
	type args struct {
		val float64
	}
	tests := []struct {
		name string
		args args
		want *apd.Decimal
	}{
		{
			name: "Zero",
			args: args{
				val: 0,
			},
			want: _DecimalZero,
		},
		{
			name: "BigFloat",
			args: args{
				val: math.MaxFloat64,
			},
			want: apd.New(17976931348623157, 292),
		},
		{
			name: "-BigFloat",
			args: args{
				val: -math.MaxFloat64,
			},
			want: apd.New(-17976931348623157, 292),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFromFloat(tt.args.val); got.Cmp(tt.want) != 0 {
				t.Errorf("newFromFloat() = %+v, want %+v", *got, *(tt.want))
			}
		})
	}
}
