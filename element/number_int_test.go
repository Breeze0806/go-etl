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
	i := &Int64{value: math.MaxInt64}
	tests := []struct {
		name string
		i    *Int64
		want BigIntNumber
	}{
		{
			name: "1",
			i:    i,
			want: i,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.BigInt(); got != tt.want {
				t.Errorf("Int64.BigInt() = %p, want %p", got, tt.want)
			}
		})
	}
}

func TestInt64_Decimal(t *testing.T) {
	i := &Int64{value: math.MaxInt64}
	tests := []struct {
		name string
		i    *Int64
		want DecimalNumber
	}{
		{
			name: "1",
			i:    i,
			want: i,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Decimal(); got != tt.want {
				t.Errorf("Int64.Decimal() = %p, want %p", got, tt.want)
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
		want *apd.BigInt
	}{
		{
			name: "1",
			i: &Int64{
				value: math.MaxInt64,
			},
			want: apd.NewBigInt(math.MaxInt64),
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
		want *apd.Decimal
	}{
		{
			name: "1",
			i: &Int64{
				value: math.MaxInt64,
			},
			want: testDecimalFormString(FormatInt64(math.MaxInt64)),
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
			name:    "0",
			i:       &Uint64{value: 0},
			want:    0,
			wantErr: false,
		},
		{
			name:    "10",
			i:       &Uint64{value: 10},
			want:    10,
			wantErr: false,
		},
		{
			name:    "maxint64",
			i:       &Uint64{value: uint64(math.MaxInt64)},
			want:    math.MaxInt64,
			wantErr: false,
		},
		{
			name:    "maxint64+1",
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
			name:    "10",
			i:       &Uint64{value: 10},
			want:    10.0,
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
	u := &Uint64{value: math.MaxUint64}
	tests := []struct {
		name string
		i    *Uint64
		want DecimalNumber
	}{
		{
			name: "1",
			i:    u,
			want: u,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Decimal(); got != tt.want {
				t.Errorf("Uint64.Decimal() = %p, want %p", got, tt.want)
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
			name: "10",
			i:    &Uint64{value: 10},
			want: "10",
		},
		{
			name: "1",
			i:    &Uint64{value: 1234567890},
			want: "1234567890",
		},
		{
			name: "1",
			i:    &Uint64{value: 123456789},
			want: "123456789",
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
			if got := tt.i.CloneBigInt(); !reflect.DeepEqual(got, tt.want) || got == tt.want {
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
		wantBi *apd.BigInt
	}{
		{
			name:   "1",
			i:      &Uint64{value: math.MaxInt64 + 1},
			wantBi: new(apd.BigInt).SetUint64(math.MaxInt64 + 1),
		},
		{
			name:   "2",
			i:      &Uint64{value: math.MaxUint64},
			wantBi: new(apd.BigInt).SetUint64(math.MaxUint64),
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
		want *apd.Decimal
	}{
		{
			name: "1",
			i:    &Uint64{value: math.MaxInt64 + 1},
			want: apd.NewWithBigInt(new(apd.BigInt).SetUint64(math.MaxInt64+1), 0),
		},
		{
			name: "2",
			i:    &Uint64{value: math.MaxUint64},
			want: apd.NewWithBigInt(new(apd.BigInt).SetUint64(math.MaxUint64), 0),
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
