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
				value: apd.NewBigInt(math.MaxInt64),
			},
			want: true,
		},
		{
			name: "2",
			b: &BigInt{
				value: apd.NewBigInt(0),
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
				value: apd.NewBigInt(math.MaxInt64),
			},
			want: math.MaxInt64,
		},
		{
			name: "2",
			b: &BigInt{
				value: apd.NewBigInt(0),
			},
			want: 0,
		},
		{
			name: "3",
			b: &BigInt{
				value: testBigIntFromString(testDecimalFormString("1e1000").Text('f')),
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
				value: apd.NewBigInt(math.MaxInt64),
			},
			wantV: 9.223372036854776e+18,
		},
		{
			name: "2",
			b: &BigInt{
				value: apd.NewBigInt(0),
			},
			wantV: 0.0,
		},
		{
			name: "3",
			b: &BigInt{
				value: testBigIntFromString(testDecimalFormString("1e1000").Text('f')),
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
	i := &BigInt{
		value: apd.NewBigInt(math.MaxInt64),
	}
	tests := []struct {
		name string
		b    *BigInt
		want BigIntNumber
	}{
		{
			name: "1",
			b:    i,
			want: i,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.BigInt(); got != tt.want {
				t.Errorf("BigInt.BigInt() = %p, want %p", got, tt.want)
			}
		})
	}
}

func TestBigInt_Decimal(t *testing.T) {
	i := &BigInt{
		value: apd.NewBigInt(math.MaxInt64),
	}
	tests := []struct {
		name string
		b    *BigInt
		want DecimalNumber
	}{
		{
			name: "1",
			b:    i,
			want: i,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Decimal(); got != tt.want {
				t.Errorf("BigInt.Decimal() = %p, want %p", got, tt.want)
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
				value: apd.NewBigInt(math.MaxInt64),
			},
			want: &BigInt{
				value: apd.NewBigInt(math.MaxInt64),
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
				value: apd.NewBigInt(math.MaxInt64),
			},
			want: FormatInt64(math.MaxInt64),
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
				value: apd.NewBigInt(math.MaxInt64),
			},
			want: &BigInt{
				value: apd.NewBigInt(math.MaxInt64),
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
	b := &BigIntStr{
		value: new(apd.BigInt).SetUint64(math.MaxUint16).String(),
	}
	tests := []struct {
		name string
		b    *BigIntStr
		want BigIntNumber
	}{
		{
			name: "1",
			b:    b,
			want: b,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.BigInt(); got != tt.want {
				t.Errorf("BigIntStr.BigInt() = %p, want %p", got, tt.want)
			}
		})
	}
}

func TestBigIntStr_Decimal(t *testing.T) {
	b := &BigIntStr{
		value: new(apd.BigInt).SetUint64(math.MaxUint16).String(),
	}
	tests := []struct {
		name string
		b    *BigIntStr
		want DecimalNumber
	}{
		{
			name: "1",
			b:    b,
			want: b,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Decimal(); got != tt.want {
				t.Errorf("BigIntStr.Decimal() = %p, want %p", got, tt.want)
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
				value: new(apd.BigInt).SetUint64(math.MaxUint16).String(),
			},
			want: &BigIntStr{
				value: new(apd.BigInt).SetUint64(math.MaxUint16).String(),
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
				value: new(apd.BigInt).SetUint64(math.MaxUint16).String(),
			},
			want: &BigIntStr{
				value: new(apd.BigInt).SetUint64(math.MaxUint16).String(),
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
		want *apd.BigInt
	}{
		{
			name: "1",
			b: &BigIntStr{
				value: new(apd.BigInt).SetUint64(math.MaxUint16).String(),
			},
			want: new(apd.BigInt).SetUint64(math.MaxUint16),
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
		want *apd.Decimal
	}{
		{
			name: "1",
			b: &BigIntStr{
				value: new(apd.BigInt).SetUint64(math.MaxUint16).String(),
			},
			want: testDecimalFormString(new(apd.BigInt).SetUint64(math.MaxUint16).String()),
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

func TestBigInt_AsDecimal(t *testing.T) {
	tests := []struct {
		name string
		b    *BigInt
		want *apd.Decimal
	}{
		{
			name: "1",
			b: &BigInt{
				value: apd.NewBigInt(math.MaxInt64),
			},
			want: apd.New(math.MaxInt64, 0),
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.AsDecimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BigInt.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}
