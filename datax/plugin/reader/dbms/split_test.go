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

package dbms

import (
	"context"
	"database/sql"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

type MockFieldTypeWithGoType struct {
	*database.BaseFieldType
}

func NewMockFieldTypeWithGoType() *MockFieldType {
	return &MockFieldType{
		BaseFieldType: database.NewBaseFieldType(&sql.ColumnType{}),
	}
}

func TestSplitRange_fetchColumn(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		s       SplitRange
		args    args
		want    element.Column
		wantErr bool
	}{
		{
			name: "1",
			s: SplitRange{
				Type: string(element.TypeBigInt),
			},
			args: args{
				"1234567890",
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValue(big.NewInt(1234567890)), "", 0),
		},
		{
			name: "2",
			s: SplitRange{
				Type: string(element.TypeBigInt),
			},
			args: args{
				"123456789a",
			},
			wantErr: true,
		},
		{
			name: "3",
			s: SplitRange{
				Type: string(element.TypeString),
			},
			args: args{
				"1234567890",
			},
			want: element.NewDefaultColumn(element.NewStringColumnValue("1234567890"), "", 0),
		},
		{
			name: "4",
			s: SplitRange{
				Type:   string(element.TypeTime),
				Layout: element.DefaultTimeFormat[:10],
			},
			args: args{
				"2009-12-11",
			},
			want: element.NewDefaultColumn(element.NewTimeColumnValueWithDecoder(time.Date(2009, 12, 11,
				0, 0, 0, 0, time.UTC), element.NewStringTimeDecoder(element.DefaultTimeFormat[:10])), "", 0),
		},
		{
			name: "5",
			s: SplitRange{
				Type:   string(element.TypeTime),
				Layout: element.DefaultTimeFormat[:10],
			},
			args: args{
				"2009-12-11 12",
			},
			wantErr: true,
		},
		{
			name: "6",
			s: SplitRange{
				Type: string(element.TypeBool),
			},
			args: args{
				"2009-12-11 12",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.fetchColumn("", tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitRange.fetchColumn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitRange.fetchColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newConvertor(t *testing.T) {
	type args struct {
		min          element.Column
		timeAccuracy string
	}
	tests := []struct {
		name    string
		args    args
		want    convertor
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				min: element.NewDefaultColumn(element.NewBigIntColumnValue(big.NewInt(1234567890)), "", 0),
			},
			want: &bigIntConvertor{},
		},
		{
			name: "2",
			args: args{
				min: element.NewDefaultColumn(element.NewStringColumnValue("1234567890"), "", 0),
			},
			want: &stringConvertor{},
		},
		{
			name: "3",
			args: args{
				min: element.NewDefaultColumn(element.NewTimeColumnValueWithDecoder(time.Date(2009, 12, 11,
					0, 0, 0, 0, time.UTC), element.NewStringTimeDecoder(element.DefaultTimeFormat[:10])), "", 0),
			},
			want: &timeConvertor{
				layout: &timeLayout{layout: element.DefaultTimeFormat[:10]},
				min:    time.Date(2009, 12, 11, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "4",
			args: args{
				min: element.NewDefaultColumn(element.NewNilBoolColumnValue(), "", 0),
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				min: element.NewDefaultColumn(element.NewNilBoolColumnValue(), "", 0),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConvertor(tt.args.min, tt.args.timeAccuracy)
			if (err != nil) != tt.wantErr {
				t.Errorf("newConvertor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newConvertor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bigIntConvertor_splitConfig(t *testing.T) {
	tests := []struct {
		name       string
		b          *bigIntConvertor
		wantTyp    string
		wantLayout string
	}{
		{
			name:    "1",
			b:       &bigIntConvertor{},
			wantTyp: element.TypeBigInt.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTyp, gotLayout := tt.b.splitConfig()
			if gotTyp != tt.wantTyp {
				t.Errorf("bigIntConvertor.splitConfig() gotTyp = %v, want %v", gotTyp, tt.wantTyp)
			}
			if gotLayout != tt.wantLayout {
				t.Errorf("bigIntConvertor.splitConfig() gotLayout = %v, want %v", gotLayout, tt.wantLayout)
			}
		})
	}
}

func Test_bigIntConvertor_fromBigInt(t *testing.T) {
	type args struct {
		bi *big.Int
	}
	tests := []struct {
		name  string
		b     *bigIntConvertor
		args  args
		wantV string
	}{
		{
			name: "1",
			b:    &bigIntConvertor{},
			args: args{
				bi: big.NewInt(1234567890),
			},
			wantV: "1234567890",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotV := tt.b.fromBigInt(tt.args.bi); gotV != tt.wantV {
				t.Errorf("bigIntConvertor.fromBigInt() = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

func Test_bigIntConvertor_toBigInt(t *testing.T) {
	type args struct {
		c element.Column
	}
	tests := []struct {
		name    string
		b       *bigIntConvertor
		args    args
		wantBi  *big.Int
		wantErr bool
	}{
		{
			name: "1",
			b:    &bigIntConvertor{},
			args: args{
				c: element.NewDefaultColumn(element.NewBigIntColumnValue(big.NewInt(1234567890)), "", 0),
			},
			wantBi: big.NewInt(1234567890),
		},
		{
			name: "2",
			b:    &bigIntConvertor{},
			args: args{
				c: element.NewDefaultColumn(element.NewStringColumnValue("abc"), "", 0),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBi, err := tt.b.toBigInt(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("bigIntConvertor.toBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBi, tt.wantBi) {
				t.Errorf("bigIntConvertor.toBigInt() = %v, want %v", gotBi, tt.wantBi)
			}
		})
	}
}

func Test_stringConvertor_splitConfig(t *testing.T) {
	tests := []struct {
		name       string
		s          *stringConvertor
		wantTyp    string
		wantLayout string
	}{
		{
			name:    "1",
			s:       &stringConvertor{},
			wantTyp: element.TypeString.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTyp, gotLayout := tt.s.splitConfig()
			if gotTyp != tt.wantTyp {
				t.Errorf("stringConvertor.splitConfig() gotTyp = %v, want %v", gotTyp, tt.wantTyp)
			}
			if gotLayout != tt.wantLayout {
				t.Errorf("stringConvertor.splitConfig() gotLayout = %v, want %v", gotLayout, tt.wantLayout)
			}
		})
	}
}

func Test_stringConvertor_fromBigInt(t *testing.T) {
	type args struct {
		bi *big.Int
	}
	tests := []struct {
		name  string
		s     *stringConvertor
		args  args
		wantV string
	}{
		{
			name: "1",
			s:    &stringConvertor{},
			args: args{
				bi: big.NewInt(1601891),
			},
			wantV: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotV := tt.s.fromBigInt(tt.args.bi); gotV != tt.wantV {
				t.Errorf("stringConvertor.fromBigInt() = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

func Test_stringConvertor_toBigInt(t *testing.T) {
	type args struct {
		c element.Column
	}
	tests := []struct {
		name    string
		s       *stringConvertor
		args    args
		wantBi  *big.Int
		wantErr bool
	}{
		{
			name: "1",
			s:    &stringConvertor{},
			args: args{
				c: element.NewDefaultColumn(element.NewStringColumnValue("abc"), "", 0),
			},
			wantBi: big.NewInt(1601891),
		},
		{
			name: "2",
			s:    &stringConvertor{},
			args: args{
				c: element.NewDefaultColumn(element.NewNilStringColumnValue(), "", 0),
			},
			wantErr: true,
		},
		{
			name: "3",
			s:    &stringConvertor{},
			args: args{
				c: element.NewDefaultColumn(element.NewStringColumnValue("中文"), "", 0),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBi, err := tt.s.toBigInt(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringConvertor.toBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBi, tt.wantBi) {
				t.Errorf("stringConvertor.toBigInt() = %v, want %v", gotBi, tt.wantBi)
			}
		})
	}
}

func Test_timeLayout_unit(t *testing.T) {
	tests := []struct {
		name string
		tr   *timeLayout
		want time.Duration
	}{
		{
			name: "day",
			tr:   &timeLayout{layout: element.DefaultTimeFormat[:10]},
			want: 24 * time.Hour,
		},
		{
			name: "min",
			tr:   &timeLayout{layout: element.DefaultTimeFormat[:16]},
			want: 1 * time.Minute,
		},
		{
			name: "s",
			tr:   &timeLayout{layout: element.DefaultTimeFormat[:19]},
			want: 1 * time.Second,
		},
		{
			name: "ms",
			tr:   &timeLayout{layout: element.DefaultTimeFormat[:23]},
			want: 1 * time.Millisecond,
		},
		{
			name: "us",
			tr:   &timeLayout{layout: element.DefaultTimeFormat[:26]},
			want: 1 * time.Microsecond,
		},
		{
			name: "ns",
			tr:   &timeLayout{layout: element.DefaultTimeFormat[:29]},
			want: 1 * time.Nanosecond,
		},
		{
			name: "hours",
			tr:   &timeLayout{layout: element.DefaultTimeFormat[:28]},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.unit(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeLayout.unit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeLayout_getLayout(t *testing.T) {
	type args struct {
		timeAccuracy string
	}
	tests := []struct {
		name string
		tr   *timeLayout
		args args
		want string
	}{
		{
			name: "day",
			tr:   &timeLayout{},
			args: args{
				timeAccuracy: "day",
			},
			want: element.DefaultTimeFormat[:10],
		},
		{
			name: "min",
			tr:   &timeLayout{},
			args: args{
				timeAccuracy: "min",
			},
			want: element.DefaultTimeFormat[:16],
		},
		{
			name: "s",
			tr:   &timeLayout{},
			args: args{
				timeAccuracy: "s",
			},
			want: element.DefaultTimeFormat[:19],
		},
		{
			name: "ms",
			tr:   &timeLayout{},
			args: args{
				timeAccuracy: "ms",
			},
			want: element.DefaultTimeFormat[:23],
		},
		{
			name: "us",
			tr:   &timeLayout{},
			args: args{
				timeAccuracy: "us",
			},
			want: element.DefaultTimeFormat[:26],
		},
		{
			name: "ns",
			tr:   &timeLayout{},
			args: args{
				timeAccuracy: "ns",
			},
			want: element.DefaultTimeFormat[:29],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.getLayout(tt.args.timeAccuracy)
			if tt.tr.layout != tt.want {
				t.Errorf("timeLayout.layout = %v, want %v", tt.tr.layout, tt.want)
			}
		})
	}
}

func Test_timeConvertor_splitConfig(t *testing.T) {
	tests := []struct {
		name       string
		tr         *timeConvertor
		wantTyp    string
		wantLayout string
	}{
		{
			name: "1",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:10],
				},
			},
			wantTyp:    element.TypeTime.String(),
			wantLayout: element.DefaultTimeFormat[:10],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTyp, gotLayout := tt.tr.splitConfig()
			if gotTyp != tt.wantTyp {
				t.Errorf("timeConvertor.splitConfig() gotTyp = %v, want %v", gotTyp, tt.wantTyp)
			}
			if gotLayout != tt.wantLayout {
				t.Errorf("timeConvertor.splitConfig() gotLayout = %v, want %v", gotLayout, tt.wantLayout)
			}
		})
	}
}

func Test_timeConvertor_fromBigInt(t *testing.T) {
	type args struct {
		bi *big.Int
	}
	tests := []struct {
		name  string
		tr    *timeConvertor
		args  args
		wantV string
	}{
		{
			name: "day",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:10],
				},
				min: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				bi: big.NewInt(44925),
			},
			wantV: time.Date(2023, 1, 1, 0, 0, 0, 0,
				time.UTC).Format(element.DefaultTimeFormat[:10]),
		},
		{
			name: "min",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:16],
				},
				min: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				bi: big.NewInt(64693357),
			},
			wantV: time.Date(2023,
				1, 1, 22, 37, 0, 0, time.UTC).Format(element.DefaultTimeFormat[:16]),
		},
		{
			name: "s",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:19],
				},
				min: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				bi: big.NewInt(3881601431),
			},
			wantV: time.Date(2023,
				1, 1, 22, 37, 11, 0, time.UTC).Format(element.DefaultTimeFormat[:19]),
		},
		{
			name: "ms",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:23],
				},
				min: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				bi: big.NewInt(3881601431999),
			},
			wantV: time.Date(2023,
				1, 1, 22, 37, 11, 999000000, time.UTC).Format(element.DefaultTimeFormat[:23]),
		},
		{
			name: "us",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:26],
				},
				min: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				bi: big.NewInt(3881601431999999),
			},
			wantV: time.Date(2023,
				1, 1, 22, 37, 11, 999999000, time.UTC).Format(element.DefaultTimeFormat[:26]),
		},
		{
			name: "ns",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:29],
				},
				min: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				bi: big.NewInt(3881601431999999999),
			},
			wantV: time.Date(2023,
				1, 1, 22, 37, 11, 999999999, time.UTC).Format(element.DefaultTimeFormat[:29]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotV := tt.tr.fromBigInt(tt.args.bi); gotV != tt.wantV {
				t.Errorf("timeConvertor.fromBigInt() = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

func Test_timeConvertor_toBigInt(t *testing.T) {
	type args struct {
		c element.Column
	}
	tests := []struct {
		name    string
		tr      *timeConvertor
		args    args
		wantBi  *big.Int
		wantErr bool
	}{
		{
			name: "day",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:10],
				},
				min: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				c: element.NewDefaultColumn(element.NewTimeColumnValue(time.Date(2023,
					1, 1, 0, 0, 0, 0, time.UTC)), "", 0),
			},
			wantBi: big.NewInt(44925),
		},
		{
			name: "min",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:16],
				},
				min: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				c: element.NewDefaultColumn(element.NewTimeColumnValue(time.Date(2023,
					1, 1, 22, 37, 0, 0, time.UTC)), "", 0),
			},
			wantBi: big.NewInt(64693357),
		},
		{
			name: "s",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:19],
				},
				min: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				c: element.NewDefaultColumn(element.NewTimeColumnValue(time.Date(2023,
					1, 1, 22, 37, 11, 0, time.UTC)), "", 0),
			},
			wantBi: big.NewInt(3881601431),
		},
		{
			name: "ms",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:23],
				},
				min: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				c: element.NewDefaultColumn(element.NewTimeColumnValue(time.Date(2023,
					1, 1, 22, 37, 11, 999000000, time.UTC)), "", 0),
			},
			wantBi: big.NewInt(3881601431999),
		},
		{
			name: "us",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:26],
				},
				min: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				c: element.NewDefaultColumn(element.NewTimeColumnValue(time.Date(2023,
					1, 1, 22, 37, 11, 999999000, time.UTC)), "", 0),
			},
			wantBi: big.NewInt(3881601431999999),
		},
		{
			name: "ns",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:29],
				},
				min: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				c: element.NewDefaultColumn(element.NewTimeColumnValue(time.Date(2023,
					1, 1, 22, 37, 11, 999999999, time.UTC)), "", 0),
			},
			wantBi: big.NewInt(3881601431999999999),
		},
		{
			name: "layout",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:28],
				},
				min: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				c: element.NewDefaultColumn(element.NewTimeColumnValue(time.Date(2023,
					1, 1, 22, 37, 11, 999999999, time.UTC)), "", 0),
			},
			wantErr: true,
		},
		{
			name: "time",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:29],
				},
				min: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				c: element.NewDefaultColumn(element.NewNilTimeColumnValue(), "", 0),
			},
			wantErr: true,
		},
		{
			name: "largeGrap",
			tr: &timeConvertor{
				layout: &timeLayout{
					layout: element.DefaultTimeFormat[:29],
				},
				min: time.Date(1000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				c: element.NewDefaultColumn(element.NewTimeColumnValue(time.Date(2023,
					1, 1, 22, 37, 11, 999999999, time.UTC)), "", 0),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBi, err := tt.tr.toBigInt(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("timeConvertor.toBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBi, tt.wantBi) {
				t.Errorf("timeConvertor.toBigInt() = %v, want %v", gotBi, tt.wantBi)
			}
		})
	}
}

func Test_doSplit(t *testing.T) {
	type args struct {
		left  *big.Int
		right *big.Int
		num   int
	}
	tests := []struct {
		name        string
		args        args
		wantResults []*big.Int
	}{
		{
			name: "1",
			args: args{
				left:  big.NewInt(12),
				right: big.NewInt(21),
				num:   5,
			},
			wantResults: []*big.Int{
				big.NewInt(12),
				big.NewInt(14),
				big.NewInt(16),
				big.NewInt(18),
				big.NewInt(20),
				big.NewInt(21),
			},
		},
		{
			name: "2",
			args: args{
				left:  big.NewInt(12),
				right: big.NewInt(12),
				num:   5,
			},
			wantResults: []*big.Int{
				big.NewInt(12),
				big.NewInt(12),
			},
		},
		{
			name: "3",
			args: args{
				left:  big.NewInt(21),
				right: big.NewInt(12),
				num:   5,
			},
			wantResults: []*big.Int{
				big.NewInt(12),
				big.NewInt(14),
				big.NewInt(16),
				big.NewInt(18),
				big.NewInt(20),
				big.NewInt(21),
			},
		},
		{
			name: "4",
			args: args{
				left:  big.NewInt(22),
				right: big.NewInt(19),
				num:   11,
			},
			wantResults: []*big.Int{
				big.NewInt(19),
				big.NewInt(20),
				big.NewInt(21),
				big.NewInt(22),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults := doSplit(tt.args.left, tt.args.right, tt.args.num)
			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("doSplit() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}

func Test_split(t *testing.T) {
	type args struct {
		min          element.Column
		max          element.Column
		num          int
		timeAccuracy string
		splitField   database.Field
	}
	tests := []struct {
		name       string
		args       args
		wantRanges []SplitRange
		wantErr    bool
	}{
		{
			name: "1",
			args: args{
				min:        element.NewDefaultColumn(element.NewBigIntColumnValue(big.NewInt(10000)), "", 0),
				max:        element.NewDefaultColumn(element.NewBigIntColumnValue(big.NewInt(50003)), "", 0),
				num:        4,
				splitField: NewMockField(database.NewBaseField(0, "f1", NewMockFieldType(database.GoTypeInt64)), NewMockFieldType(database.GoTypeInt64)),
			},
			wantRanges: []SplitRange{
				{
					Type:  element.TypeBigInt.String(),
					Left:  "10000",
					Right: "20001",
					where: "f1 >= $1 and f1 < $2",
				},
				{
					Type:  element.TypeBigInt.String(),
					Left:  "20001",
					Right: "30002",
					where: "f1 >= $1 and f1 < $2",
				},
				{
					Type:  element.TypeBigInt.String(),
					Left:  "30002",
					Right: "40003",
					where: "f1 >= $1 and f1 < $2",
				},
				{
					Type:  element.TypeBigInt.String(),
					Left:  "40003",
					Right: "50003",
					where: "f1 >= $1 and f1 <= $2",
				},
			},
		},
		{
			name: "2",
			args: args{
				min:        element.NewDefaultColumn(element.NewBigIntColumnValue(big.NewInt(10000)), "", 0),
				max:        element.NewDefaultColumn(element.NewBigIntColumnValue(big.NewInt(50003)), "", 0),
				num:        0,
				splitField: NewMockField(database.NewBaseField(0, "f1", NewMockFieldType(database.GoTypeInt64)), NewMockFieldType(database.GoTypeInt64)),
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				min:        element.NewDefaultColumn(element.NewBigIntColumnValue(big.NewInt(10000)), "", 0),
				max:        nil,
				num:        4,
				splitField: NewMockField(database.NewBaseField(0, "f1", NewMockFieldType(database.GoTypeInt64)), NewMockFieldType(database.GoTypeInt64)),
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				min:        nil,
				max:        element.NewDefaultColumn(element.NewBigIntColumnValue(big.NewInt(50003)), "", 0),
				num:        4,
				splitField: NewMockField(database.NewBaseField(0, "f1", NewMockFieldType(database.GoTypeInt64)), NewMockFieldType(database.GoTypeInt64)),
			},
			wantErr: true,
		},
		{
			name: "5",
			args: args{
				min:        element.NewDefaultColumn(element.NewBoolColumnValue(true), "", 0),
				max:        element.NewDefaultColumn(element.NewBoolColumnValue(true), "", 0),
				num:        4,
				splitField: NewMockField(database.NewBaseField(0, "f1", NewMockFieldType(database.GoTypeInt64)), NewMockFieldType(database.GoTypeInt64)),
			},
			wantErr: true,
		},
		{
			name: "6",
			args: args{
				min:        element.NewDefaultColumn(element.NewStringColumnValue("中文"), "", 0),
				max:        element.NewDefaultColumn(element.NewStringColumnValue("abc"), "", 0),
				num:        4,
				splitField: NewMockField(database.NewBaseField(0, "f1", NewMockFieldType(database.GoTypeInt64)), NewMockFieldType(database.GoTypeInt64)),
			},
			wantErr: true,
		},
		{
			name: "7",
			args: args{
				min:        element.NewDefaultColumn(element.NewStringColumnValue("abc"), "", 0),
				max:        element.NewDefaultColumn(element.NewStringColumnValue("中文"), "", 0),
				num:        4,
				splitField: NewMockField(database.NewBaseField(0, "f1", NewMockFieldType(database.GoTypeInt64)), NewMockFieldType(database.GoTypeInt64)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRanges, err := split(tt.args.min, tt.args.max, tt.args.num, tt.args.timeAccuracy, tt.args.splitField)
			if (err != nil) != tt.wantErr {
				t.Errorf("split() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRanges, tt.wantRanges) {
				t.Errorf("split() = %+v, want %+v", gotRanges, tt.wantRanges)
			}
		})
	}
}

func TestSplitConfig_fetchMin(t *testing.T) {
	type args struct {
		ctx   context.Context
		field database.Field
	}
	tests := []struct {
		name    string
		s       SplitConfig
		args    args
		wantC   element.Column
		wantErr bool
	}{
		{
			name: "1",
			s: SplitConfig{
				Key: "f1",
				Range: SplitRange{
					Type: string(element.TypeBigInt),
					Left: "100000",
				},
			},
			args: args{
				ctx: context.TODO(),
				field: NewMockField(database.NewBaseField(0, "f1",
					NewMockFieldType(database.GoTypeInt64)),
					NewMockFieldType(database.GoTypeInt64)),
			},
			wantC: element.NewDefaultColumn(element.NewBigIntColumnValue(big.NewInt(100000)),
				"f1", 0),
		},
		{
			name: "2",
			s: SplitConfig{
				Key: "f1",
				Range: SplitRange{
					Type: string(element.TypeString),
					Left: "100000",
				},
			},
			args: args{
				ctx: context.TODO(),
				field: NewMockField(database.NewBaseField(0, "f1", NewMockFieldType(database.GoTypeInt64)),
					NewMockFieldType(database.GoTypeInt64)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			splitTable := NewMockTable(database.NewBaseTable("db", "schema", "table"))
			splitTable.AppendField(tt.args.field)
			gotC, err := tt.s.fetchMin(tt.args.ctx, splitTable)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitConfig.fetchMin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("SplitConfig.fetchMin() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestSplitConfig_fetchMax(t *testing.T) {
	type args struct {
		ctx   context.Context
		field database.Field
	}
	tests := []struct {
		name    string
		s       SplitConfig
		args    args
		wantC   element.Column
		wantErr bool
	}{
		{
			name: "1",
			s: SplitConfig{
				Key: "f1",
				Range: SplitRange{
					Type:  string(element.TypeBigInt),
					Right: "100000",
				},
			},
			args: args{
				ctx: context.TODO(),
				field: NewMockField(database.NewBaseField(0, "f1",
					NewMockFieldType(database.GoTypeInt64)),
					NewMockFieldType(database.GoTypeInt64)),
			},
			wantC: element.NewDefaultColumn(element.NewBigIntColumnValue(big.NewInt(100000)),
				"f1", 0),
		},
		{
			name: "2",
			s: SplitConfig{
				Range: SplitRange{
					Type: string(element.TypeString),
				},
			},
			args: args{
				ctx: context.TODO(),
				field: NewMockField(database.NewBaseField(0, "f1", NewMockFieldType(database.GoTypeInt64)),
					NewMockFieldType(database.GoTypeInt64)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			splitTable := NewMockTable(database.NewBaseTable("db", "schema", "table"))
			splitTable.AppendField(tt.args.field)
			gotC, err := tt.s.fetchMax(tt.args.ctx, splitTable)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitConfig.fetchMax() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("SplitConfig.fetchMax() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestSplitConfig_checkType(t *testing.T) {
	type args struct {
		field database.Field
	}
	tests := []struct {
		name    string
		s       SplitConfig
		args    args
		wantErr bool
	}{
		{
			name: "1",
			s: SplitConfig{
				Key: "f1",
				Range: SplitRange{
					Type: string(element.TypeBigInt),
				},
			},
			args: args{
				field: NewMockField(database.NewBaseField(0, "f1",
					NewMockFieldType(database.GoTypeInt64)),
					NewMockFieldType(database.GoTypeInt64)),
			},
		},
		{
			name: "2",
			s: SplitConfig{
				Key: "f1",
				Range: SplitRange{
					Type: string(element.TypeString),
				},
			},
			args: args{
				field: NewMockField(database.NewBaseField(0, "f1",
					NewMockFieldType(database.GoTypeInt64)),
					NewMockFieldType(database.GoTypeInt64)),
			},
			wantErr: true,
		},
		{
			name: "3",
			s: SplitConfig{
				Key: "f1",
				Range: SplitRange{
					Type: string(element.TypeBigInt),
				},
			},
			args: args{
				field: NewMockField(database.NewBaseField(0, "f1",
					NewMockFieldType(database.GoTypeString)),
					NewMockFieldType(database.GoTypeString)),
			},
		},
		{
			name: "4",
			s: SplitConfig{
				Key: "f1",
				Range: SplitRange{
					Type: string(element.TypeString),
				},
			},
			args: args{
				field: NewMockField(database.NewBaseField(0, "f1",
					NewMockFieldType(database.GoTypeString)),
					NewMockFieldType(database.GoTypeString)),
			},
		},
		{
			name: "5",
			s: SplitConfig{
				Key: "f1",
				Range: SplitRange{
					Type: string(element.TypeDecimal),
				},
			},
			args: args{
				field: NewMockField(database.NewBaseField(0, "f1",
					NewMockFieldType(database.GoTypeString)),
					NewMockFieldType(database.GoTypeString)),
			},
			wantErr: true,
		},
		{
			name: "6",
			s: SplitConfig{
				Key: "f1",
				Range: SplitRange{
					Type: string(element.TypeTime),
				},
			},
			args: args{
				field: NewMockField(database.NewBaseField(0, "f1",
					NewMockFieldType(database.GoTypeTime)),
					NewMockFieldType(database.GoTypeTime)),
			},
		},
		{
			name: "7",
			s: SplitConfig{
				Key: "f1",
				Range: SplitRange{
					Type: string(element.TypeString),
				},
			},
			args: args{
				field: NewMockField(database.NewBaseField(0, "f1",
					NewMockFieldType(database.GoTypeTime)),
					NewMockFieldType(database.GoTypeTime)),
			},
			wantErr: true,
		},
		{
			name: "8",
			s: SplitConfig{
				Key: "f1",
				Range: SplitRange{
					Type: string(element.TypeString),
				},
			},
			args: args{
				field: NewMockField(database.NewBaseField(0, "f1",
					NewMockFieldType(database.GoTypeTime)),
					NewMockFieldTypeWithGoType()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			splitTable := NewMockTable(database.NewBaseTable("db", "schema", "table"))
			splitTable.AppendField(tt.args.field)
			if err := tt.s.checkType(splitTable); (err != nil) != tt.wantErr {
				t.Errorf("SplitConfig.checkType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSplitConfig_setLayout(t *testing.T) {
	tests := []struct {
		name    string
		s       SplitConfig
		want    string
		wantErr bool
	}{
		{
			name: "1",
			s: SplitConfig{
				TimeAccuracy: "day",
				Range: SplitRange{
					Type: string(element.TypeTime),
				},
			},
			want:    "2006-01-02",
			wantErr: false,
		},
		{
			name: "2",
			s: SplitConfig{
				TimeAccuracy: "",
				Range: SplitRange{
					Type: string(element.TypeTime),
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.setLayout(); (err != nil) != tt.wantErr {
				t.Errorf("SplitConfig.setLayout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.s.Range.Layout != tt.want {
				t.Errorf("SplitConfig.Range.Layout = %v, want %v", tt.s.Range.Layout, tt.want)
			}
		})
	}
}
