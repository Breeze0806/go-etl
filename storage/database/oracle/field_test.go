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

package oracle

import (
	"database/sql/driver"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/godror/godror"
)

type mockColumnType struct {
	name string
}

func newMockColumnType(name string) *mockColumnType {
	return &mockColumnType{
		name: name,
	}
}

func (m *mockColumnType) Name() string {
	return ""
}

func (m *mockColumnType) ScanType() reflect.Type {
	return nil
}

func (m *mockColumnType) Length() (length int64, ok bool) {
	return
}

func (m *mockColumnType) DecimalSize() (precision, scale int64, ok bool) {
	return
}

func (m *mockColumnType) Nullable() (nullable, ok bool) {
	return
}

func (m *mockColumnType) DatabaseTypeName() string {
	return m.name
}

func (m *mockColumnType) IsSupported() bool {
	return true
}

func mustDecimalColumnValueFromString(s string) element.ColumnValue {
	c, err := element.NewDecimalColumnValueFromString(s)
	if err != nil {
		panic(err)
	}
	return c
}

func TestField_Quoted(t *testing.T) {
	tests := []struct {
		name string
		f    *Field
		want string
	}{
		{
			name: "1",
			f:    NewField(database.NewBaseField(0, "f1", newMockColumnType(""))),
			want: `"f1"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Quoted(); got != tt.want {
				t.Errorf("Field.Quoted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_BindVar(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		f    *Field
		args args
		want string
	}{
		{
			name: "1",
			f:    NewField(database.NewBaseField(0, "f1", newMockColumnType(""))),
			args: args{
				i: 1,
			},
			want: ":1",
		},
		{
			name: "2",
			f:    NewField(database.NewBaseField(0, "f1", newMockColumnType("TIMESTAMP"))),
			args: args{
				i: 1,
			},
			want: "to_timestamp(:1,'yyyy-mm-dd hh24:mi:ss.ff9')",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.BindVar(tt.args.i); got != tt.want {
				t.Errorf("Field.BindVar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_Select(t *testing.T) {
	tests := []struct {
		name string
		f    *Field
		want string
	}{
		{
			name: "1",
			f:    NewField(database.NewBaseField(0, "f1", newMockColumnType(""))),
			want: `"f1"`,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Select(); got != tt.want {
				t.Errorf("Field.Select() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_Type(t *testing.T) {
	tests := []struct {
		name string
		f    *Field
		want database.FieldType
	}{
		{
			name: "1",
			f:    NewField(database.NewBaseField(0, "f1", newMockColumnType("BOOLEAN"))),
			want: NewFieldType(newMockColumnType("BOOLEAN")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Field.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_Scanner(t *testing.T) {
	tests := []struct {
		name string
		f    *Field
		want database.Scanner
	}{
		{
			name: "1",
			f:    NewField(database.NewBaseField(0, "f1", newMockColumnType(""))),
			want: NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("")))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Scanner(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Field.Scanner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_Valuer(t *testing.T) {
	type args struct {
		c element.Column
	}
	tests := []struct {
		name string
		f    *Field
		args args
		want database.Valuer
	}{
		{
			name: "1",
			f:    NewField(database.NewBaseField(0, "f1", newMockColumnType("DOUBLE"))),
			args: args{
				c: element.NewDefaultColumn(element.NewNilBoolColumnValue(), "f1", 0),
			},
			want: NewValuer(NewField(database.NewBaseField(0, "f1", newMockColumnType("DOUBLE"))),
				element.NewDefaultColumn(element.NewNilBoolColumnValue(), "f1", 0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Valuer(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Field.Valuer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldType_IsSupportted(t *testing.T) {
	tests := []struct {
		name string
		f    *FieldType
		want bool
	}{
		{
			name: "1",
			f:    NewFieldType(newMockColumnType("DATE")),
			want: true,
		},
		{
			name: "1",
			f:    NewFieldType(newMockColumnType("DATETIME")),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.IsSupported(); got != tt.want {
				t.Errorf("FieldType.IsSupportted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_Scan(t *testing.T) {
	type args struct {
		src any
	}
	tests := []struct {
		name    string
		s       *Scanner
		conf    *config.JSON
		args    args
		want    element.Column
		wantErr bool
	}{
		{
			name: "BOOLEAN",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("BOOLEAN")))),
			args: args{
				src: true,
			},
			want: element.NewDefaultColumn(element.NewBoolColumnValue(true), "f1", 1),
		},
		{
			name: "BOOLEANnil",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("BOOLEAN")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBoolColumnValue(), "f1", 0),
		},
		{
			name: "BOOLEANerr",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("BOOLEAN")))),
			args: args{
				src: "true",
			},
			wantErr: true,
		},

		{
			name: "BINARY_INTEGERint64",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("BINARY_INTEGER")))),
			args: args{
				src: int64(math.MaxInt64),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(int64(math.MaxInt64)),
				"f1", element.ByteSize(int64(math.MaxInt64))),
		},
		{
			name: "BINARY_INTEGERuint64",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("BINARY_INTEGER")))),
			args: args{
				src: uint64(math.MaxUint64),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromUint64(uint64(math.MaxUint64)),
				"f1", element.ByteSize(uint64(math.MaxUint64))),
		},
		{
			name: "BINARY_INTEGERnil",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("BINARY_INTEGER")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBigIntColumnValue(), "f1", 0),
		},
		{
			name: "BINARY_INTEGERerr",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("BINARY_INTEGER")))),
			args: args{
				src: 123,
			},
			wantErr: true,
		},

		{
			name: "RAW",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("RAW")))),
			args: args{
				src: []byte("中文"),
			},
			want: element.NewDefaultColumn(element.NewBytesColumnValue([]byte("中文")),
				"f1", element.ByteSize([]byte("中文"))),
		},
		{
			name: "RAWnil",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("RAW")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBytesColumnValue(), "f1", 0),
		},
		{
			name: "LONG RAW err",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("LONG RAW")))),
			args: args{
				src: "中文",
			},
			wantErr: true,
		},

		{
			name: "DATE",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("DATE")))),
			args: args{
				src: time.Date(2022, 10, 16, 10, 18, 33, 999999999, time.UTC),
			},
			want: element.NewDefaultColumn(element.NewTimeColumnValueWithDecoder(
				time.Date(2022, 10, 16, 10, 18, 33, 999999999, time.UTC),
				element.NewStringTimeDecoder(dateLayout)), "f1",
				element.ByteSize(time.Date(2022, 10, 16, 10, 18, 33, 999999999, time.UTC))),
		},
		{
			name: "DATE nil",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("DATE")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilTimeColumnValue(), "f1", 0),
		},
		{
			name: "DATE err",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("DATE")))),
			args: args{
				src: "中文",
			},
			wantErr: true,
		},

		{
			name: "TIMESTAMP",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("TIMESTAMP")))),
			args: args{
				src: time.Date(2022, 10, 16, 10, 18, 33, 999999999, time.UTC),
			},
			want: element.NewDefaultColumn(element.NewTimeColumnValueWithDecoder(
				time.Date(2022, 10, 16, 10, 18, 33, 999999999, time.UTC),
				element.NewStringTimeDecoder(datetimeLayout)), "f1",
				element.ByteSize(time.Date(2022, 10, 16, 10, 18, 33, 999999999, time.UTC))),
		},
		{
			name: "TIMESTAMP WITH TIME ZONE nil",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("TIMESTAMP WITH TIME ZONE")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilTimeColumnValue(), "f1", 0),
		},
		{
			name: "TIMESTAMP WITH LOCAL TIME ZONE err",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("TIMESTAMP WITH LOCAL TIME ZONE")))),
			args: args{
				src: "中文",
			},
			wantErr: true,
		},

		{
			name: "VARCHAR2",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("VARCHAR2")))),
			args: args{
				src: "中文abc-123",
			},
			want: element.NewDefaultColumn(element.NewStringColumnValue("中文abc-123"), "f1",
				element.ByteSize("中文abc-123")),
		},
		{
			name: "CHAR nil",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("CHAR")))),
			args: args{
				src: "",
			},
			want: element.NewDefaultColumn(element.NewNilStringColumnValue(), "f1", 0),
		},
		{
			name: "CHARTrim",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("CHAR")))),
			args: args{
				src: " 中文abc-123     ",
			},
			conf: testJSONFromString(`{"trimChar":true}`),
			want: element.NewDefaultColumn(element.NewStringColumnValue("中文abc-123"), "f1",
				element.ByteSize(" 中文abc-123     ")),
		},
		{
			name: "CHAR",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("CHAR")))),
			args: args{
				src: " 中文abc-123     ",
			},
			want: element.NewDefaultColumn(element.NewStringColumnValue(" 中文abc-123     "), "f1",
				element.ByteSize(" 中文abc-123     ")),
		},
		{
			name: "VARCHAR2 err",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("VARCHAR2")))),
			args: args{
				src: 12,
			},
			wantErr: true,
		},

		{
			name: "FLOAT nil",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("FLOAT")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilDecimalColumnValue(), "f1", 0),
		},
		{
			name: "FLOAT float32",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("FLOAT")))),
			args: args{
				src: float32(8.23),
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat32(float32(8.23)),
				"f1", element.ByteSize(float32(8.23))),
		},
		{
			name: "DOUBLE float64",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("DOUBLE")))),
			args: args{
				src: float64(8.23),
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(8.23),
				"f1", element.ByteSize(float64(8.23))),
		},
		{
			name: "NUMBER int64",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("NUMBER")))),
			args: args{
				src: int64(1234567890),
			},
			want: element.NewDefaultColumn(mustDecimalColumnValueFromString("1234567890"),
				"f1", element.ByteSize(int64(1234567890))),
		},
		{
			name: "NUMBER uint64",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("NUMBER")))),
			args: args{
				src: uint64(1234567890),
			},
			want: element.NewDefaultColumn(mustDecimalColumnValueFromString("1234567890"),
				"f1", element.ByteSize(uint64(1234567890))),
		},
		{
			name: "NUMBER bool",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("NUMBER")))),
			args: args{
				src: true,
			},
			want: element.NewDefaultColumn(mustDecimalColumnValueFromString("1"), "f1", 1),
		},
		{
			name: "NUMBER",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("DOUBLE")))),
			args: args{
				src: godror.Number("8.23"),
			},
			want: element.NewDefaultColumn(mustDecimalColumnValueFromString("8.23"), "f1", 4),
		},
		{
			name: "NUMBER err",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("DOUBLE")))),
			args: args{
				src: godror.Number("8.23a"),
			},
			wantErr: true,
		},
		{
			name: "NUMBER type err",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("DOUBLE")))),
			args: args{
				src: "8.23",
			},
			wantErr: true,
		},

		{
			name: "INT",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockColumnType("INT")))),
			args: args{
				src: "8.23",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.conf != nil {
				tt.s.f.SetConfig(tt.conf)
			}

			if err := tt.s.Scan(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("Scanner.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.s.Column(), tt.want) {
				t.Errorf("Column() = %v, want %v", tt.s.Column(), tt.want)
			}
		})
	}
}

func TestValuer_Value(t *testing.T) {
	tests := []struct {
		name    string
		v       *Valuer
		want    driver.Value
		wantErr bool
	}{
		{
			name: "1",
			v: NewValuer(NewField(database.NewBaseField(0, "f1", newMockColumnType("BOOLEAN"))),
				element.NewDefaultColumn(element.NewNilBoolColumnValue(), "f1", 0)),
			want: driver.Value(""),
		},
		{
			name: "2",
			v: NewValuer(NewField(database.NewBaseField(0, "f1", newMockColumnType("BOOLEAN"))),
				element.NewDefaultColumn(element.NewBoolColumnValue(true), "f1", 0)),
			want: driver.Value("1"),
		},
		{
			name: "3",
			v: NewValuer(NewField(database.NewBaseField(0, "f1", newMockColumnType("BOOLEAN"))),
				element.NewDefaultColumn(element.NewBoolColumnValue(false), "f2", 0)),
			want: driver.Value("0"),
		},
		{
			name: "4",
			v: NewValuer(NewField(database.NewBaseField(0, "f1", newMockColumnType("NUMBER"))),
				element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1234567890), "f2", 0)),
			want: driver.Value("1234567890"),
		},
		{
			name: "5",
			v: NewValuer(NewField(database.NewBaseField(0, "f1", newMockColumnType("NUMBER"))),
				element.NewDefaultColumn(element.NewNilBoolColumnValue(), "f2", 0)),
			want: driver.Value(""),
		},
		{
			name: "6",
			v: NewValuer(NewField(database.NewBaseField(0, "f1", newMockColumnType("BLOB"))),
				element.NewDefaultColumn(element.NewNilBoolColumnValue(), "f2", 0)),
			want: driver.Value(nil),
		},
		{
			name: "7",
			v: NewValuer(NewField(database.NewBaseField(0, "f1", newMockColumnType("RAW"))),
				element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1), "f2", 0)),
			want: driver.Value([]byte("1")),
		},
		{
			name: "8",
			v: NewValuer(NewField(database.NewBaseField(0, "f1", newMockColumnType("BOOLEAN"))),
				element.NewDefaultColumn(element.NewStringColumnValue("we"), "f2", 0)),
			wantErr: true,
		},
		{
			name: "9",
			v: NewValuer(NewField(database.NewBaseField(0, "f1", newMockColumnType("DATE"))),
				element.NewDefaultColumn(element.NewTimeColumnValue(time.Date(2025, 12, 1, 23, 47, 11, 21, time.UTC)), "f2", 0)),
			want: driver.Value(time.Date(2025, 12, 1, 23, 47, 11, 21, time.UTC)),
		},
		{
			name: "10",
			v: NewValuer(NewField(database.NewBaseField(0, "f1", newMockColumnType("DATE"))),
				element.NewDefaultColumn(element.NewNilTimeColumnValue(), "f2", 0)),
			want: driver.Value(time.Time{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.v.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Valuer.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Valuer.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
