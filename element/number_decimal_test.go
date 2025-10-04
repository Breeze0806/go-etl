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
	"testing"

	"github.com/cockroachdb/apd/v3"
)

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
				value:  testDecimalFormString("1.797693134862315708145274237317043567981e+308").Text('f') + ".123",
				intLen: 309,
			},
			want: &BigIntStr{
				value: testDecimalFormString("1.797693134862315708145274237317043567981e+308").Text('f'),
			},
		},
		{
			name: "2",
			d: &DecimalStr{
				value:  FormatInt64(math.MaxInt64),
				intLen: len(FormatInt64(math.MaxInt64)),
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
	d := &DecimalStr{
		value:  testDecimalFormString("1.797693134862315708145274237317043567981e+308").String() + ".123",
		intLen: 309,
	}
	tests := []struct {
		name string
		d    *DecimalStr
		want DecimalNumber
	}{
		{
			name: "1",
			d:    d,
			want: d,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Decimal(); got != tt.want {
				t.Errorf("DecimalStr.Decimal() = %p, want %p", got, tt.want)
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
		want *apd.Decimal
	}{
		{
			name: "1",
			d: &DecimalStr{
				value:  testDecimalFormString("1.797693134862315708145274237317043567981e+308").Text('f') + ".123",
				intLen: 309,
			},
			want: testDecimalFormString(testDecimalFormString("1.797693134862315708145274237317043567981e+308").Text('f') + ".123"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.AsDecimal(); got.Cmp(tt.want) != 0 {
				t.Errorf("DecimalStr.AsDecimal() = %v, want %v", *got, *(tt.want))
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
				value: _DecimalZero,
			},
			want: false,
		},
		{
			name: "2",
			d: &Decimal{
				value: NewApdDecimalFromFloat(1e32),
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
				value: NewApdDecimalFromFloat(math.MaxFloat64),
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
				value: apd.NewBigInt(123456232542542525),
			},
		},
		{
			name: "2",
			d: &Decimal{
				value: testDecimalFormString("-123456232542542525.525254252524"),
			},
			want: &BigInt{
				value: apd.NewBigInt(-123456232542542525),
			},
		},
		{
			name: "3",
			d: &Decimal{
				value: testDecimalFormString("0.00122323123"),
			},
			want: &BigInt{
				value: apd.NewBigInt(0),
			},
		},
		{
			name: "4",
			d: &Decimal{
				value: testDecimalFormString("123450000"),
			},
			want: &BigInt{
				value: apd.NewBigInt(123450000),
			},
		},

		{
			name: "5",
			d: &Decimal{
				value: testDecimalFormString("1.00122323123"),
			},
			want: &BigInt{
				value: apd.NewBigInt(1),
			},
		},
		{
			name: "6",
			d: &Decimal{
				value: testDecimalFormString("12345.67e-3"),
			},
			want: &BigInt{
				value: apd.NewBigInt(12),
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
	d := &Decimal{
		value: testDecimalFormString("123456232542542525.525254252524"),
	}
	tests := []struct {
		name string
		d    *Decimal
		want DecimalNumber
	}{
		{
			name: "1",
			d:    d,
			want: d,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Decimal(); got != tt.want {
				t.Errorf("*apd.Decimal() = %p, want %p", got, tt.want)
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
		want *apd.Decimal
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
