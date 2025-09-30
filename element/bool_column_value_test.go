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
	"reflect"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"
)

func TestNilBoolColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		n    *NilBoolColumnValue
		want ColumnType
	}{
		{
			name: "1",
			n:    NewNilBoolColumnValue().(*NilBoolColumnValue),
			want: TypeBool,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilBoolColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilBoolColumnValue_Clone(t *testing.T) {
	tests := []struct {
		name string
		n    *NilBoolColumnValue
		want ColumnValue
	}{
		{
			name: "1",
			n:    NewNilBoolColumnValue().(*NilBoolColumnValue),
			want: NewNilBoolColumnValue(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.Clone()
			if got == tt.n {
				t.Errorf("NilBigIntColumnValue.Clone() = %p, n %p want %p", got, tt.n, tt.want)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilBoolColumnValue.clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		b    *BoolColumnValue
		want ColumnType
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: TypeBool,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoolColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_AsBool(t *testing.T) {
	tests := []struct {
		name    string
		b       *BoolColumnValue
		want    bool
		wantErr bool
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: true,
		},
		{
			name: "false",
			b:    NewBoolColumnValue(false).(*BoolColumnValue),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolColumnValue.AsBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BoolColumnValue.AsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_AsBigInt(t *testing.T) {
	tests := []struct {
		name    string
		b       *BoolColumnValue
		want    *apd.BigInt
		wantErr bool
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: apd.NewBigInt(1),
		},
		{
			name: "false",
			b:    NewBoolColumnValue(false).(*BoolColumnValue),
			want: apd.NewBigInt(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsBigInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolColumnValue.AsBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.AsBigInt().Cmp(tt.want) != 0 {
				t.Errorf("BoolColumnValue.AsBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_AsDecimal(t *testing.T) {
	tests := []struct {
		name    string
		b       *BoolColumnValue
		want    *apd.Decimal
		wantErr bool
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: apd.New(1, 0),
		},
		{
			name: "false",
			b:    NewBoolColumnValue(false).(*BoolColumnValue),
			want: _DecimalZero,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsDecimal()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolColumnValue.AsDecimal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.AsDecimal().Cmp(tt.want) != 0 {
				t.Errorf("BoolColumnValue.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_AsString(t *testing.T) {
	tests := []struct {
		name    string
		b       *BoolColumnValue
		want    string
		wantErr bool
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: "true",
		},
		{
			name: "false",
			b:    NewBoolColumnValue(false).(*BoolColumnValue),
			want: "false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolColumnValue.AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BoolColumnValue.AsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_AsBytes(t *testing.T) {
	tests := []struct {
		name    string
		b       *BoolColumnValue
		want    string
		wantErr bool
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: "true",
		},
		{
			name: "false",
			b:    NewBoolColumnValue(false).(*BoolColumnValue),
			want: "false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolColumnValue.AsBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, []byte(tt.want)) {
				t.Errorf("BoolColumnValue.AsBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_AsTime(t *testing.T) {
	tests := []struct {
		name    string
		b       *BoolColumnValue
		want    time.Time
		wantErr bool
	}{
		{
			name:    "true",
			b:       NewBoolColumnValue(true).(*BoolColumnValue),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolColumnValue.AsTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoolColumnValue.AsTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_String(t *testing.T) {
	tests := []struct {
		name string
		b    *BoolColumnValue
		want string
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: "true",
		},
		{
			name: "false",
			b:    NewBoolColumnValue(false).(*BoolColumnValue),
			want: "false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("BoolColumnValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_Clone(t *testing.T) {
	tests := []struct {
		name string
		b    *BoolColumnValue
		want ColumnValue
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: NewBoolColumnValue(true),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.Clone()
			if got == tt.b {
				t.Errorf("BoolColumnValue.Clone() = %p, b %p", got, tt.b)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoolColumnValue.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_Cmp(t *testing.T) {
	type args struct {
		right ColumnValue
	}
	tests := []struct {
		name    string
		b       *BoolColumnValue
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "1",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			args: args{
				right: NewNilBoolColumnValue(),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "2",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			args: args{
				right: NewBoolColumnValue(true),
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "3",
			b:    NewBoolColumnValue(false).(*BoolColumnValue),
			args: args{
				right: NewBoolColumnValue(false),
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "4",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			args: args{
				right: NewBoolColumnValue(false),
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "5",
			b:    NewBoolColumnValue(false).(*BoolColumnValue),
			args: args{
				right: NewBoolColumnValue(true),
			},
			want:    -1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Cmp(tt.args.right)
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolColumnValue.Cmp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BoolColumnValue.Cmp() = %v, want %v", got, tt.want)
			}
		})
	}
}
