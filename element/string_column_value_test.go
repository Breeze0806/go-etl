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
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNilStringColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		n    *NilStringColumnValue
		want ColumnType
	}{
		{
			name: "1",
			n:    NewNilStringColumnValue().(*NilStringColumnValue),
			want: TypeString,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilStringColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilStringColumnValue_Clone(t *testing.T) {
	tests := []struct {
		name string
		n    *NilStringColumnValue
		want ColumnValue
	}{
		{
			name: "1",
			n:    NewNilStringColumnValue().(*NilStringColumnValue),
			want: NewNilStringColumnValue(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.Clone()
			if got == tt.n {
				t.Errorf("NilStringColumnValue.Clone() = %p, n %p", got, tt.n)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilStringColumnValue.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		s    *StringColumnValue
		want ColumnType
	}{
		{
			name: "1",
			s:    NewStringColumnValue("").(*StringColumnValue),
			want: TypeString,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringColumnValue_AsBool(t *testing.T) {
	tests := []struct {
		name    string
		s       *StringColumnValue
		want    bool
		wantErr bool
	}{
		//1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
		{
			name: "1",
			s:    NewStringColumnValue("1").(*StringColumnValue),
			want: true,
		},
		{
			name: "t",
			s:    NewStringColumnValue("t").(*StringColumnValue),
			want: true,
		},
		{
			name: "T",
			s:    NewStringColumnValue("T").(*StringColumnValue),
			want: true,
		},
		{
			name: "TRUE",
			s:    NewStringColumnValue("TRUE").(*StringColumnValue),
			want: true,
		},
		{
			name: "true",
			s:    NewStringColumnValue("true").(*StringColumnValue),
			want: true,
		},
		{
			name: "True",
			s:    NewStringColumnValue("True").(*StringColumnValue),
			want: true,
		},
		{
			name: "0",
			s:    NewStringColumnValue("0").(*StringColumnValue),
			want: false,
		},
		{
			name: "f",
			s:    NewStringColumnValue("f").(*StringColumnValue),
			want: false,
		},
		{
			name: "FALSE",
			s:    NewStringColumnValue("FALSE").(*StringColumnValue),
			want: false,
		},
		{
			name: "F",
			s:    NewStringColumnValue("F").(*StringColumnValue),
			want: false,
		},
		{
			name: "false",
			s:    NewStringColumnValue("false").(*StringColumnValue),
			want: false,
		},
		{
			name: "False",
			s:    NewStringColumnValue("False").(*StringColumnValue),
			want: false,
		},
		{
			name:    "FAlse",
			s:       NewStringColumnValue("FAlse").(*StringColumnValue),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AsBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("StringColumnValue.AsBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringColumnValue.AsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringColumnValue_AsBigInt(t *testing.T) {
	tests := []struct {
		name    string
		s       *StringColumnValue
		want    *big.Int
		wantErr bool
	}{
		{
			name: "1",
			s:    NewStringColumnValue("12340000.0").(*StringColumnValue),
			want: testBigIntFromString("12340000"),
		},
		{
			name: "2",
			s:    NewStringColumnValue("1234213213214135465736545425353980988.0").(*StringColumnValue),
			want: testBigIntFromString("1234213213214135465736545425353980988"),
		},
		{
			name: "3",
			s:    NewStringColumnValue("-12340000.3").(*StringColumnValue),
			want: testBigIntFromString("-12340000"),
		},
		{
			name: "4",
			s:    NewStringColumnValue("1.12345689").(*StringColumnValue),
			want: testBigIntFromString("1"),
		},
		{
			name: "5",
			s:    NewStringColumnValue("1.23456e4").(*StringColumnValue),
			want: testBigIntFromString("12345"),
		},
		{
			name:    "6",
			s:       NewStringColumnValue("1.23456e4adad").(*StringColumnValue),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AsBigInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("StringColumnValue.AsBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Cmp(tt.want) != 0 {
				t.Errorf("StringColumnValue.AsBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringColumnValue_AsDecimal(t *testing.T) {
	tests := []struct {
		name    string
		s       *StringColumnValue
		want    decimal.Decimal
		wantErr bool
	}{
		{
			name: "1",
			s:    NewStringColumnValue("12340000.0").(*StringColumnValue),
			want: testDecimalFormString("12340000"),
		},
		{
			name: "2",
			s:    NewStringColumnValue("1234213213214135465736545425353980988.0").(*StringColumnValue),
			want: testDecimalFormString("1234213213214135465736545425353980988"),
		},
		{
			name: "3",
			s:    NewStringColumnValue("-12340000.3").(*StringColumnValue),
			want: testDecimalFormString("-12340000.3"),
		},
		{
			name: "4",
			s:    NewStringColumnValue("1.12345689").(*StringColumnValue),
			want: testDecimalFormString("1.12345689"),
		},
		{
			name: "5",
			s:    NewStringColumnValue("1.23456e4").(*StringColumnValue),
			want: testDecimalFormString("12345.6"),
		},
		{
			name: "6",
			s:    NewStringColumnValue("1e100").(*StringColumnValue),
			want: testDecimalFormString("1e100"),
		},
		{
			name:    "7",
			s:       NewStringColumnValue("1.23456e4adad").(*StringColumnValue),
			wantErr: true,
		},
		{
			name:    "8",
			s:       NewStringColumnValue("1.23456e4adad").(*StringColumnValue),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AsDecimal()
			if (err != nil) != tt.wantErr {
				t.Errorf("StringColumnValue.AsDecimal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("StringColumnValue.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringColumnValue_AsString(t *testing.T) {
	tests := []struct {
		name    string
		s       *StringColumnValue
		want    string
		wantErr bool
	}{
		{
			name: "1",
			s:    NewStringColumnValue("中文abc1234<>&*^%$").(*StringColumnValue),
			want: "中文abc1234<>&*^%$",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("StringColumnValue.AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringColumnValue.AsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringColumnValue_AsBytes(t *testing.T) {
	tests := []struct {
		name    string
		s       *StringColumnValue
		want    []byte
		wantErr bool
	}{
		{
			name: "1",
			s:    NewStringColumnValue("中文abc1234<>&*^%$").(*StringColumnValue),
			want: []byte("中文abc1234<>&*^%$"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AsBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("StringColumnValue.AsBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringColumnValue.AsBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringColumnValue_AsTime(t *testing.T) {
	tests := []struct {
		name    string
		s       *StringColumnValue
		wantT   time.Time
		wantErr bool
	}{
		{
			name:  "1",
			s:     NewStringColumnValue(time.Date(2020, 12, 17, 22, 49, 56, 69-999-999, time.Local).Format(DefaultTimeFormat)).(*StringColumnValue),
			wantT: time.Date(2020, 12, 17, 22, 49, 56, 69-999-999, time.Local),
		},
		{
			name:    "2",
			s:       NewStringColumnValue("中文abc1234<>&*^%$").(*StringColumnValue),
			wantErr: true,
		},
		{
			name:  "3",
			s:     NewStringColumnValue(time.Date(2020, 12, 17, 22, 49, 56, 69-999-999, time.UTC).Format("2006-01-02")).(*StringColumnValue),
			wantT: time.Date(2020, 12, 17, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := tt.s.AsTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("StringColumnValue.AsTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !gotT.Equal(tt.wantT) {
				t.Errorf("StringColumnValue.AsTime() = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}

func TestStringColumnValue_String(t *testing.T) {
	tests := []struct {
		name string
		s    *StringColumnValue
		want string
	}{
		{
			name: "1",
			s:    NewStringColumnValue("中文abc1234<>&*^%$").(*StringColumnValue),
			want: "中文abc1234<>&*^%$",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("StringColumnValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringColumnValue_Clone(t *testing.T) {
	tests := []struct {
		name string
		s    *StringColumnValue
		want ColumnValue
	}{
		{
			name: "1",
			s:    NewStringColumnValue("中文abc1234<>&*^%$").(*StringColumnValue),
			want: NewStringColumnValue("中文abc1234<>&*^%$"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Clone()
			if got == tt.s {
				t.Errorf("StringColumnValue.Clone() = %p, s %v", got, tt.s)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringColumnValue.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringColumnValue_Cmp(t *testing.T) {
	type args struct {
		right ColumnValue
	}
	tests := []struct {
		name    string
		s       *StringColumnValue
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "1",
			s:    NewStringColumnValue("123").(*StringColumnValue),
			args: args{
				right: NewNilBigIntColumnValue(),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "2",
			s:    NewStringColumnValue("abc").(*StringColumnValue),
			args: args{
				right: NewStringColumnValue("abcd"),
			},
			want:    -1,
			wantErr: false,
		},
		{
			name: "3",
			s:    NewStringColumnValue("abc").(*StringColumnValue),
			args: args{
				right: NewStringColumnValue("abd"),
			},
			want:    -1,
			wantErr: false,
		},
		{
			name: "4",
			s:    NewStringColumnValue("abc").(*StringColumnValue),
			args: args{
				right: NewStringColumnValue("abc"),
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "5",
			s:    NewStringColumnValue("abcd").(*StringColumnValue),
			args: args{
				right: NewStringColumnValue("abc"),
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "6",
			s:    NewStringColumnValue("abd").(*StringColumnValue),
			args: args{
				right: NewStringColumnValue("abc"),
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Cmp(tt.args.right)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringColumnValue.Cmp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringColumnValue.Cmp() = %v, want %v", got, tt.want)
			}
		})
	}
}
