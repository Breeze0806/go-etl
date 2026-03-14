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
)

func TestNewNilJsonColumnValue(t *testing.T) {
	tests := []struct {
		name string
		want ColumnValue
	}{
		{
			name: "1",
			want: &NilJsonColumnValue{
				nilColumnValue: &nilColumnValue{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNilJsonColumnValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNilJsonColumnValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilJsonColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		n    *NilJsonColumnValue
		want ColumnType
	}{
		{
			name: "1",
			n:    NewNilJsonColumnValue().(*NilJsonColumnValue),
			want: TypeJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilJsonColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilJsonColumnValue_Clone(t *testing.T) {
	tests := []struct {
		name string
		n    *NilJsonColumnValue
		want ColumnValue
	}{
		{
			name: "1",
			n:    NewNilJsonColumnValue().(*NilJsonColumnValue),
			want: NewNilJsonColumnValue(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.Clone()
			if got == tt.n {
				t.Errorf("NilJsonColumnValue.Clone() = %p, n %p", got, tt.n)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilJsonColumnValue.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJsonColumnValueFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				s: `{"a":1}`,
			},
			want: `{"a":1}`,
		},
		{
			name: "2",
			args: args{
				s: `[1,2,3]`,
			},
			want: `[1,2,3]`,
		},
		{
			name: "3",
			args: args{
				s: `"hello"`,
			},
			want: `"hello"`,
		},
		{
			name: "4",
			args: args{
				s: `123`,
			},
			want: `123`,
		},
		{
			name:    "5",
			args:    args{s: `invalid json`},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJsonColumnValueFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJsonColumnValueFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.want {
				t.Errorf("NewJsonColumnValueFromString() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestNewJsonColumnValueFromBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				b: []byte(`{"a":1}`),
			},
			want: `{"a":1}`,
		},
		{
			name: "2",
			args: args{
				b: []byte(`[1,2,3]`),
			},
			want: `[1,2,3]`,
		},
		{
			name:    "3",
			args:    args{b: []byte(`invalid json`)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJsonColumnValueFromBytes(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJsonColumnValueFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.want {
				t.Errorf("NewJsonColumnValueFromBytes() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestJsonColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		j    *JsonColumnValue
		want ColumnType
	}{
		{
			name: "1",
			j:    getTestJsonColumnValue(`{"a":1}`),
			want: TypeJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonColumnValue_AsBool(t *testing.T) {
	tests := []struct {
		name    string
		j       *JsonColumnValue
		want    bool
		wantErr bool
	}{
		{
			name:    "1",
			j:       getTestJsonColumnValue(`{"a":1}`),
			wantErr: true,
		},
		{
			name:    "2",
			j:       getTestJsonColumnValue(`true`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.AsBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonColumnValue.AsBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JsonColumnValue.AsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonColumnValue_AsBigInt(t *testing.T) {
	tests := []struct {
		name    string
		j       *JsonColumnValue
		wantErr bool
	}{
		{
			name:    "1",
			j:       getTestJsonColumnValue(`{"a":1}`),
			wantErr: true,
		},
		{
			name:    "2",
			j:       getTestJsonColumnValue(`123`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.AsBigInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonColumnValue.AsBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				t.Errorf("JsonColumnValue.AsBigInt() = %v, want nil", got)
			}
		})
	}
}

func TestJsonColumnValue_AsDecimal(t *testing.T) {
	tests := []struct {
		name    string
		j       *JsonColumnValue
		wantErr bool
	}{
		{
			name:    "1",
			j:       getTestJsonColumnValue(`{"a":1}`),
			wantErr: true,
		},
		{
			name:    "2",
			j:       getTestJsonColumnValue(`123.45`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.AsDecimal()
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonColumnValue.AsDecimal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				t.Errorf("JsonColumnValue.AsDecimal() = %v, want nil", got)
			}
		})
	}
}

func TestJsonColumnValue_AsString(t *testing.T) {
	tests := []struct {
		name    string
		j       *JsonColumnValue
		want    string
		wantErr bool
	}{
		{
			name: "1",
			j:    getTestJsonColumnValue(`{"a":1}`),
			want: `{"a":1}`,
		},
		{
			name: "2",
			j:    getTestJsonColumnValue(`[1,2,3]`),
			want: `[1,2,3]`,
		},
		{
			name: "3",
			j:    getTestJsonColumnValue(`"hello"`),
			want: `"hello"`,
		},
		{
			name: "4",
			j:    getTestJsonColumnValue(`123`),
			want: `123`,
		},
		{
			name: "5",
			j:    getTestJsonColumnValue(`true`),
			want: `true`,
		},
		{
			name: "6",
			j:    getTestJsonColumnValue(`null`),
			want: `null`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonColumnValue.AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JsonColumnValue.AsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonColumnValue_AsBytes(t *testing.T) {
	tests := []struct {
		name    string
		j       *JsonColumnValue
		want    []byte
		wantErr bool
	}{
		{
			name: "1",
			j:    getTestJsonColumnValue(`{"a":1}`),
			want: []byte(`{"a":1}`),
		},
		{
			name: "2",
			j:    getTestJsonColumnValue(`[1,2,3]`),
			want: []byte(`[1,2,3]`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.AsBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonColumnValue.AsBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonColumnValue.AsBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonColumnValue_AsTime(t *testing.T) {
	tests := []struct {
		name    string
		j       *JsonColumnValue
		wantErr bool
	}{
		{
			name:    "1",
			j:       getTestJsonColumnValue(`{"a":1}`),
			wantErr: true,
		},
		{
			name:    "2",
			j:       getTestJsonColumnValue(`"2020-01-01"`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.j.AsTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonColumnValue.AsTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJsonColumnValue_AsJSON(t *testing.T) {
	tests := []struct {
		name    string
		j       *JsonColumnValue
		want    string
		wantErr bool
	}{
		{
			name: "1",
			j:    getTestJsonColumnValue(`{"a":1}`),
			want: `{"a":1}`,
		},
		{
			name: "2",
			j:    getTestJsonColumnValue(`[1,2,3]`),
			want: `[1,2,3]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.AsJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonColumnValue.AsJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ToString() != tt.want {
				t.Errorf("JsonColumnValue.AsJSON() = %v, want %v", got.ToString(), tt.want)
			}
		})
	}
}

func TestJsonColumnValue_String(t *testing.T) {
	tests := []struct {
		name string
		j    *JsonColumnValue
		want string
	}{
		{
			name: "1",
			j:    getTestJsonColumnValue(`{"a":1}`),
			want: `{"a":1}`,
		},
		{
			name: "2",
			j:    getTestJsonColumnValue(`[1,2,3]`),
			want: `[1,2,3]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.String(); got != tt.want {
				t.Errorf("JsonColumnValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonColumnValue_Clone(t *testing.T) {
	tests := []struct {
		name string
		j    *JsonColumnValue
	}{
		{
			name: "1",
			j:    getTestJsonColumnValue(`{"a":1}`),
		},
		{
			name: "2",
			j:    getTestJsonColumnValue(`[1,2,3]`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.j.Clone()
			if got == tt.j {
				t.Errorf("JsonColumnValue.Clone() = %p, j %p", got, tt.j)
			}
			if got.String() != tt.j.String() {
				t.Errorf("JsonColumnValue.Clone() = %v, want %v", got.String(), tt.j.String())
			}
		})
	}
}

func TestJsonColumnValue_Cmp(t *testing.T) {
	type args struct {
		right ColumnValue
	}
	tests := []struct {
		name    string
		j       *JsonColumnValue
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "1",
			j:    getTestJsonColumnValue(`{"a":1}`),
			args: args{
				right: getTestJsonColumnValue(`{"a":1}`),
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "2",
			j:    getTestJsonColumnValue(`{"a":1}`),
			args: args{
				right: getTestJsonColumnValue(`{"a":2}`),
			},
			want:    -1,
			wantErr: false,
		},
		{
			name: "3",
			j:    getTestJsonColumnValue(`{"a":2}`),
			args: args{
				right: getTestJsonColumnValue(`{"a":1}`),
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "4",
			j:    getTestJsonColumnValue(`[1,2,3]`),
			args: args{
				right: NewNilJsonColumnValue(),
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.Cmp(tt.args.right)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonColumnValue.Cmp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JsonColumnValue.Cmp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestJsonColumnValue(s string) *JsonColumnValue {
	cv, _ := NewJsonColumnValueFromString(s)
	return cv.(*JsonColumnValue)
}
