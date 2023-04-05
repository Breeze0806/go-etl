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
	"testing"
	"time"
)

func Test_notNilColumnValue_IsNil(t *testing.T) {
	tests := []struct {
		name string
		n    *notNilColumnValue
		want bool
	}{
		{
			name: "1",
			n:    &notNilColumnValue{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.IsNil(); got != tt.want {
				t.Errorf("notNilColumnValue.IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		n    *nilColumnValue
		want ColumnType
	}{
		{
			name: "1",
			n:    &nilColumnValue{},
			want: TypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nilColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_IsNil(t *testing.T) {
	tests := []struct {
		name string
		n    *nilColumnValue
		want bool
	}{
		{
			name: "1",
			n:    &nilColumnValue{},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.IsNil(); got != tt.want {
				t.Errorf("nilColumnValue.IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_AsBool(t *testing.T) {
	tests := []struct {
		name    string
		n       *nilColumnValue
		want    bool
		wantErr bool
	}{
		{
			name:    "1",
			n:       &nilColumnValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.AsBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("nilColumnValue.AsBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("nilColumnValue.AsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_AsBigInt(t *testing.T) {
	tests := []struct {
		name    string
		n       *nilColumnValue
		want    BigIntNumber
		wantErr bool
	}{
		{
			name:    "1",
			n:       &nilColumnValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.AsBigInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("nilColumnValue.AsBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nilColumnValue.AsBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_AsDecimal(t *testing.T) {
	tests := []struct {
		name    string
		n       *nilColumnValue
		want    DecimalNumber
		wantErr bool
	}{
		{
			name:    "1",
			n:       &nilColumnValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.AsDecimal()
			if (err != nil) != tt.wantErr {
				t.Errorf("nilColumnValue.AsDecimal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nilColumnValue.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_AsString(t *testing.T) {
	tests := []struct {
		name    string
		n       *nilColumnValue
		want    string
		wantErr bool
	}{
		{
			name:    "1",
			n:       &nilColumnValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("nilColumnValue.AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("nilColumnValue.AsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_AsBytes(t *testing.T) {
	tests := []struct {
		name    string
		n       *nilColumnValue
		want    []byte
		wantErr bool
	}{
		{
			name:    "1",
			n:       &nilColumnValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.AsBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("nilColumnValue.AsBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nilColumnValue.AsBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_AsTime(t *testing.T) {
	tests := []struct {
		name    string
		n       *nilColumnValue
		want    time.Time
		wantErr bool
	}{
		{
			name:    "1",
			n:       &nilColumnValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.AsTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("nilColumnValue.AsTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nilColumnValue.AsTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_String(t *testing.T) {
	tests := []struct {
		name string
		n    *nilColumnValue
		want string
	}{
		{
			name: "1",
			n:    &nilColumnValue{},
			want: "<nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.String(); got != tt.want {
				t.Errorf("nilColumnValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultColumn_Name(t *testing.T) {
	tests := []struct {
		name string
		d    *DefaultColumn
		want string
	}{
		{
			name: "1",
			d:    NewDefaultColumn(NewNilBigIntColumnValue(), "test", 12).(*DefaultColumn),
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Name(); got != tt.want {
				t.Errorf("DefaultColumn.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestDefaultColumn_ByteSize(t *testing.T) {
	tests := []struct {
		name string
		d    *DefaultColumn
		want int64
	}{
		{
			name: "1",
			d:    NewDefaultColumn(NewNilBigIntColumnValue(), "test", 12).(*DefaultColumn),
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.ByteSize(); got != tt.want {
				t.Errorf("DefaultColumn.ByteSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultColumn_MemorySize(t *testing.T) {
	tests := []struct {
		name string
		d    *DefaultColumn
		want int64
	}{
		{
			name: "1",
			d:    NewDefaultColumn(NewNilBigIntColumnValue(), "test", 12).(*DefaultColumn),
			want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.MemorySize(); got != tt.want {
				t.Errorf("DefaultColumn.MemorySize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColumnType_String(t *testing.T) {
	tests := []struct {
		name string
		c    ColumnType
		want string
	}{
		{
			name: "1",
			c:    ColumnType("yyy"),
			want: "yyy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("ColumnType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultColumn_AsInt64(t *testing.T) {
	tests := []struct {
		name    string
		d       *DefaultColumn
		want    int64
		wantErr bool
	}{
		{
			name:    "1",
			d:       NewDefaultColumn(NewStringColumnValue("1e10"), "test", 0).(*DefaultColumn),
			want:    10000000000,
			wantErr: false,
		},
		{
			name:    "2",
			d:       NewDefaultColumn(NewStringColumnValue("13.3z"), "test", 0).(*DefaultColumn),
			want:    0,
			wantErr: true,
		},
		{
			name:    "3",
			d:       NewDefaultColumn(NewStringColumnValue("1e22"), "test", 0).(*DefaultColumn),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.AsInt64()
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultColumn.AsInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DefaultColumn.AsInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultColumn_AsFloat64(t *testing.T) {
	tests := []struct {
		name    string
		d       *DefaultColumn
		want    float64
		wantErr bool
	}{
		{
			name:    "1",
			d:       NewDefaultColumn(NewStringColumnValue("-1.23456789e10"), "test", 0).(*DefaultColumn),
			want:    -1.23456789e10,
			wantErr: false,
		},
		{
			name: "2",
			d: NewDefaultColumn(NewStringColumnValue(strconv.FormatFloat(float64(math.MaxFloat64),
				'f', -1, 64)), "test", 0).(*DefaultColumn),
			want:    float64(math.MaxFloat64),
			wantErr: false,
		},
		{
			name:    "3",
			d:       NewDefaultColumn(NewStringColumnValue("13.3z"), "test", 0).(*DefaultColumn),
			want:    0,
			wantErr: true,
		},
		{
			name:    "4",
			d:       NewDefaultColumn(NewStringColumnValue("1e1000"), "test", 0).(*DefaultColumn),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.AsFloat64()
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultColumn.AsFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DefaultColumn.AsFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultColumn_Clone(t *testing.T) {
	tests := []struct {
		name    string
		d       *DefaultColumn
		want    Column
		wantErr bool
	}{
		{
			name:    "1",
			d:       NewDefaultColumn(NewNilBigIntColumnValue(), "test", 0).(*DefaultColumn),
			want:    NewDefaultColumn(NewNilBigIntColumnValue(), "test", 0),
			wantErr: false,
		},

		{
			name:    "2",
			d:       NewDefaultColumn(newMockColumnValue(), "test", 0).(*DefaultColumn),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Clone()
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultColumn.Clone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultColumn.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultColumn_Cmp(t *testing.T) {
	type args struct {
		c Column
	}
	tests := []struct {
		name    string
		d       *DefaultColumn
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "1",
			d:    NewDefaultColumn(newMockColumnValue(), "f1", 0).(*DefaultColumn),
			args: args{
				c: NewDefaultColumn(newMockColumnValue(), "f2", 0),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "2",
			d:    NewDefaultColumn(newMockColumnValue(), "f1", 0).(*DefaultColumn),
			args: args{
				c: NewDefaultColumn(newMockColumnValue(), "f1", 0),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "3",
			d:    NewDefaultColumn(NewBigIntColumnValueFromInt64(1), "f1", 0).(*DefaultColumn),
			args: args{
				c: NewDefaultColumn(NewBigIntColumnValueFromInt64(1), "f1", 0),
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Cmp(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultColumn.Cmp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DefaultColumn.Cmp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteSize(t *testing.T) {
	type args struct {
		src interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				src: nil,
			},
			want: 0,
		},
		{
			name: "2",
			args: args{
				src: true,
			},
			want: 1,
		},
		{
			name: "3",
			args: args{
				src: "src",
			},
			want: 3,
		},
		{
			name: "4",
			args: args{
				src: []byte("src"),
			},
			want: 3,
		},
		{
			name: "5",
			args: args{
				src: 1,
			},
			want: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ByteSize(tt.args.src); got != tt.want {
				t.Errorf("ByteSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
