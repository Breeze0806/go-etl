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
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNewNilBigIntColumnValue(t *testing.T) {
	tests := []struct {
		name string
		want ColumnValue
	}{
		{
			name: "1",
			want: NewNilBigIntColumnValue(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNilBigIntColumnValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNilBigIntColumnValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilBigIntColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		n    *NilBigIntColumnValue
		want ColumnType
	}{
		{
			name: "1",
			n:    NewNilBigIntColumnValue().(*NilBigIntColumnValue),
			want: TypeBigInt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilBigIntColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBigIntColumnValueFromInt64(t *testing.T) {
	type args struct {
		v int64
	}
	tests := []struct {
		name string
		args args
		want ColumnValue
	}{
		{
			name: "Zero",
			args: args{
				v: 0,
			},
			want: NewBigIntColumnValueFromInt64(0),
		},
		{
			name: "MaxInt",
			args: args{
				v: math.MaxInt64,
			},
			want: NewBigIntColumnValueFromInt64(math.MaxInt64),
		},
		{
			name: "MinInt",
			args: args{
				v: math.MinInt64,
			},
			want: NewBigIntColumnValueFromInt64(math.MinInt64),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBigIntColumnValueFromInt64(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBigIntColumnValueFromInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBigIntColumnValueFromBigInt(t *testing.T) {
	type args struct {
		v *big.Int
	}
	tests := []struct {
		name string
		args args
		want ColumnValue
	}{
		{
			name: "Zero",
			args: args{
				v: big.NewInt(0),
			},
			want: NewBigIntColumnValue(big.NewInt(0)),
		},
		{
			name: "FromMaxInt",
			args: args{
				v: big.NewInt(math.MaxInt64),
			},
			want: NewBigIntColumnValue(big.NewInt(math.MaxInt64)),
		},
		{
			name: "MinInt",
			args: args{
				v: big.NewInt(math.MinInt64),
			},
			want: NewBigIntColumnValue(big.NewInt(math.MinInt64)),
		},
		{
			name: "MaxUint",
			args: args{
				v: new(big.Int).SetUint64(math.MaxUint64),
			},
			want: NewBigIntColumnValue(new(big.Int).SetUint64(math.MaxUint64)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBigIntColumnValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBigIntColumnValueFromBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBigIntColumnValueFromString(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		want    ColumnValue
		wantErr bool
	}{
		{
			name: "Zero",
			args: args{
				v: big.NewInt(0).String(),
			},
			want: NewBigIntColumnValue(big.NewInt(0)),
		},
		{
			name: "MaxInt",
			args: args{
				v: big.NewInt(math.MaxInt64).String(),
			},
			want: NewBigIntColumnValue(big.NewInt(math.MaxInt64)),
		},
		{
			name: "MinInt",
			args: args{
				v: big.NewInt(math.MinInt64).String(),
			},
			want: NewBigIntColumnValue(big.NewInt(math.MinInt64)),
		},
		{
			name: "MaxUint",
			args: args{
				v: new(big.Int).SetUint64(math.MaxUint64).String(),
			},
			want: NewBigIntColumnValue(new(big.Int).SetUint64(math.MaxUint64)),
		},
		{
			name: "NegUint",
			args: args{
				v: "-" + new(big.Int).SetUint64(math.MaxUint64).String(),
			},
			want: NewBigIntColumnValue(new(big.Int).Neg(new(big.Int).SetUint64(math.MaxUint64))),
		},

		{
			name: "NegUint1",
			args: args{
				v: "-0000" + new(big.Int).SetUint64(math.MaxUint64).String(),
			},
			want: NewBigIntColumnValue(new(big.Int).Neg(new(big.Int).SetUint64(math.MaxUint64))),
		},
		{
			name: "BigInt1",
			args: args{
				v: "213231321312391283921481934823742365746982570",
			},
			want: NewBigIntColumnValue(testBigIntFromString("213231321312391283921481934823742365746982570")),
		},
		{
			name: "BigInt2",
			args: args{v: "-213231321312391283921481934823742365746982570"},
			want: NewBigIntColumnValue(testBigIntFromString("-213231321312391283921481934823742365746982570")),
		},
		{
			name: "UnValidNumber",
			args: args{
				v: new(big.Int).SetUint64(math.MaxUint64).String() + ".00000",
			},
			wantErr: true,
		},

		{
			name: "UnValidNumber1",
			args: args{
				v: new(big.Int).SetUint64(math.MaxUint64).String() + "abc",
			},
			wantErr: true,
		},
		{
			name: "UnValidNumber2",
			args: args{
				v: new(big.Int).SetUint64(math.MaxUint64).String() + "e19",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBigIntColumnValueFromString(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBigIntColumnValueFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.want.String() {
				t.Errorf("NewBigIntColumnValueFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		b    *BigIntColumnValue
		want ColumnType
	}{
		{
			name: "1",
			b:    NewBigIntColumnValueFromInt64(0).(*BigIntColumnValue),
			want: TypeBigInt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BigIntColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntColumnValue_AsBool(t *testing.T) {
	tests := []struct {
		name    string
		b       *BigIntColumnValue
		want    bool
		wantErr bool
	}{
		{
			name:    "1",
			b:       NewBigIntColumnValueFromInt64(0).(*BigIntColumnValue),
			want:    false,
			wantErr: false,
		},
		{
			name:    "2",
			b:       NewBigIntColumnValueFromInt64(-3000000000000000).(*BigIntColumnValue),
			want:    true,
			wantErr: false,
		},
		{
			name:    "3",
			b:       NewBigIntColumnValueFromInt64(math.MaxInt64).(*BigIntColumnValue),
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("BigIntColumnValue.AsBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BigIntColumnValue.AsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntColumnValue_AsBigInt(t *testing.T) {
	tests := []struct {
		name    string
		b       *BigIntColumnValue
		want    *big.Int
		wantErr bool
	}{
		{
			name: "1",
			b:    NewBigIntColumnValueFromInt64(0).(*BigIntColumnValue),
			want: _IntZero,
		},

		{
			name: "2",
			b:    NewBigIntColumnValueFromInt64(10).(*BigIntColumnValue),
			want: _IntTen,
		},

		{
			name: "3",
			b:    NewBigIntColumnValueFromInt64(-123456789123456789).(*BigIntColumnValue),
			want: big.NewInt(-123456789123456789),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsBigInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("BigIntColumnValue.AsBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.AsBigInt().Cmp(tt.want) != 0 {
				t.Errorf("BigIntColumnValue.AsBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntColumnValue_AsDecimal(t *testing.T) {
	tests := []struct {
		name    string
		b       *BigIntColumnValue
		want    decimal.Decimal
		wantErr bool
	}{
		{
			name:    "1",
			b:       NewBigIntColumnValueFromInt64(0).(*BigIntColumnValue),
			want:    decimal.Zero,
			wantErr: false,
		},
		{
			name:    "2",
			b:       NewBigIntColumnValueFromInt64(1).(*BigIntColumnValue),
			want:    decimal.NewFromInt(1),
			wantErr: false,
		},

		{
			name:    "3",
			b:       NewBigIntColumnValueFromInt64(123456789123456789).(*BigIntColumnValue),
			want:    decimal.NewFromInt(123456789123456789),
			wantErr: false,
		},
		{
			name:    "4",
			b:       NewBigIntColumnValueFromInt64(-123456789123456789).(*BigIntColumnValue),
			want:    decimal.NewFromInt(-123456789123456789),
			wantErr: false,
		},
		{
			name:    "5",
			b:       testBigIntColumnValueFromString("-893748265832572758237426582375023840295258259025"),
			want:    testDecimalFormString("-893748265832572758237426582375023840295258259025"),
			wantErr: false,
		},
		{
			name:    "6",
			b:       testBigIntColumnValueFromString("89374826583257275823742658237502384029525825900000"),
			want:    testDecimalFormString("89374826583257275823742658237502384029525825900000"),
			wantErr: false,
		},
		{
			name:    "7",
			b:       testBigIntColumnValueFromString("-89374826583257275823742658237502384029525825900000"),
			want:    testDecimalFormString("-89374826583257275823742658237502384029525825900000"),
			wantErr: false,
		},
		{
			name:    "8",
			b:       testBigIntColumnValueFromString("893748265832572758237426582375023840295258259025"),
			want:    testDecimalFormString("893748265832572758237426582375023840295258259025"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsDecimal()
			if (err != nil) != tt.wantErr {
				t.Errorf("BigIntColumnValue.AsDecimal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.AsDecimal().Equal(tt.want) {
				t.Errorf("BigIntColumnValue.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntColumnValue_AsString(t *testing.T) {
	tests := []struct {
		name    string
		b       *BigIntColumnValue
		want    string
		wantErr bool
	}{
		{
			name: "1",
			b:    NewBigIntColumnValueFromInt64(-123456789123456789).(*BigIntColumnValue),
			want: "-123456789123456789",
		},
		{
			name: "2",
			b:    NewBigIntColumnValueFromInt64(123456789123456789).(*BigIntColumnValue),
			want: "123456789123456789",
		},
		{
			name: "3",
			b:    testBigIntColumnValueFromString("00007493653627562385674265283652364628429316574085236974073249567251974816481"),
			want: "7493653627562385674265283652364628429316574085236974073249567251974816481",
		},
		{
			name: "4",
			b:    testBigIntColumnValueFromString("-00007493653627562385674265283652364628429316574085236974073249567251974816481"),
			want: "-7493653627562385674265283652364628429316574085236974073249567251974816481",
		},
		{
			name: "5",
			b:    testBigIntColumnValueFromString("00"),
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("BigIntColumnValue.AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BigIntColumnValue.AsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntColumnValue_AsBytes(t *testing.T) {
	tests := []struct {
		name    string
		b       *BigIntColumnValue
		want    string
		wantErr bool
	}{
		{
			name: "1",
			b:    NewBigIntColumnValueFromInt64(-123456789123456789).(*BigIntColumnValue),
			want: "-123456789123456789",
		},
		{
			name: "2",
			b:    NewBigIntColumnValueFromInt64(123456789123456789).(*BigIntColumnValue),
			want: "123456789123456789",
		},
		{
			name: "3",
			b:    testBigIntColumnValueFromString("00007493653627562385674265283652364628429316574085236974073249567251974816481"),
			want: "7493653627562385674265283652364628429316574085236974073249567251974816481",
		},
		{
			name: "4",
			b:    testBigIntColumnValueFromString("-00007493653627562385674265283652364628429316574085236974073249567251974816481"),
			want: "-7493653627562385674265283652364628429316574085236974073249567251974816481",
		},
		{
			name: "5",
			b:    testBigIntColumnValueFromString("00"),
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("BigIntColumnValue.AsBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, []byte(tt.want)) {
				t.Errorf("BigIntColumnValue.AsBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntColumnValue_AsTime(t *testing.T) {
	tests := []struct {
		name    string
		b       *BigIntColumnValue
		want    time.Time
		wantErr bool
	}{
		{
			name:    "1",
			b:       testBigIntColumnValueFromString("1234"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("BigIntColumnValue.AsTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BigIntColumnValue.AsTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntColumnValue_String(t *testing.T) {
	tests := []struct {
		name string
		b    *BigIntColumnValue
		want string
	}{
		{
			name: "1",
			b:    NewBigIntColumnValueFromInt64(-123456789123456789).(*BigIntColumnValue),
			want: "-123456789123456789",
		},
		{
			name: "2",
			b:    NewBigIntColumnValueFromInt64(123456789123456789).(*BigIntColumnValue),
			want: "123456789123456789",
		},
		{
			name: "3",
			b:    testBigIntColumnValueFromString("00007493653627562385674265283652364628429316574085236974073249567251974816481"),
			want: "7493653627562385674265283652364628429316574085236974073249567251974816481",
		},
		{
			name: "4",
			b:    testBigIntColumnValueFromString("-00007493653627562385674265283652364628429316574085236974073249567251974816481"),
			want: "-7493653627562385674265283652364628429316574085236974073249567251974816481",
		},
		{
			name: "5",
			b:    testBigIntColumnValueFromString("00"),
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("BigIntColumnValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntColumnValue_Clone(t *testing.T) {
	tests := []struct {
		name string
		b    *BigIntColumnValue
		want ColumnValue
	}{
		{
			name: "1",
			b:    testBigIntColumnValueFromString("-00007493653627562385674265283652364628429316574085236974073249567251974816481"),
			want: testBigIntColumnValueFromString("-00007493653627562385674265283652364628429316574085236974073249567251974816481"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.Clone()
			if got == tt.b {
				t.Errorf("NilBigIntColumnValue.Clone() = %p, b %p want %p", got, tt.b, tt.want)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilBigIntColumnValue.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilBigIntColumnValue_Clone(t *testing.T) {
	tests := []struct {
		name string
		n    *NilBigIntColumnValue
		want ColumnValue
	}{
		{
			name: "1",
			n:    NewNilBigIntColumnValue().(*NilBigIntColumnValue),
			want: NewNilBigIntColumnValue(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.Clone()
			if got == tt.n {
				t.Errorf("NilBigIntColumnValue.Clone() = %p, n %p want %p", got, tt.n, tt.want)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilBigIntColumnValue.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilBigIntColumnValue_IsNil(t *testing.T) {
	tests := []struct {
		name string
		n    *NilBigIntColumnValue
		want bool
	}{
		{
			name: "1",
			n:    NewNilBigIntColumnValue().(*NilBigIntColumnValue),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.IsNil()
			if got != tt.want {
				t.Errorf("NilBigIntColumnValue.IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntColumnValue_IsNil(t *testing.T) {
	tests := []struct {
		name string
		n    *BigIntColumnValue
		want bool
	}{
		{
			name: "1",
			n:    NewBigIntColumnValueFromInt64(11).(*BigIntColumnValue),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.IsNil()
			if got != tt.want {
				t.Errorf("NilBigIntColumnValue.IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigIntColumnValue_Cmp(t *testing.T) {
	type args struct {
		right ColumnValue
	}
	tests := []struct {
		name    string
		b       *BigIntColumnValue
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "1",
			b:    NewBigIntColumnValueFromInt64(int64(math.MaxInt64)).(*BigIntColumnValue),
			args: args{
				right: NewNilBigIntColumnValue(),
			},
			want:    0,
			wantErr: true,
		},

		{
			name: "2",
			b:    NewBigIntColumnValueFromInt64(int64(math.MaxInt64)).(*BigIntColumnValue),
			args: args{
				right: NewBigIntColumnValueFromInt64(int64(math.MaxInt64 - 1)),
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "3",
			b:    NewBigIntColumnValueFromInt64(int64(math.MaxInt64)).(*BigIntColumnValue),
			args: args{
				right: NewBigIntColumnValueFromInt64(int64(math.MaxInt64)),
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "4",
			b:    NewBigIntColumnValueFromInt64(int64(math.MinInt64)).(*BigIntColumnValue),
			args: args{
				right: NewBigIntColumnValueFromInt64(int64(math.MinInt64 + 1)),
			},
			want:    -1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Cmp(tt.args.right)
			if (err != nil) != tt.wantErr {
				t.Errorf("BigIntColumnValue.Cmp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BigIntColumnValue.Cmp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBigIntColumnValueFromUint64(t *testing.T) {
	type args struct {
		v uint64
	}
	tests := []struct {
		name string
		args args
		want ColumnValue
	}{
		{
			name: "Zero",
			args: args{
				v: 0,
			},
			want: &BigIntColumnValue{
				val: &Uint64{value: 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBigIntColumnValueFromUint64(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBigIntColumnValueFromUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}
