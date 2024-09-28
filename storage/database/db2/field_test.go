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

package db2

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

type mockFieldType struct {
	name string
}

func newMockFieldType(name string) *mockFieldType {
	return &mockFieldType{
		name: name,
	}
}

func (m *mockFieldType) Name() string {
	return ""
}

func (m *mockFieldType) ScanType() reflect.Type {
	return nil
}

func (m *mockFieldType) Length() (length int64, ok bool) {
	return
}

func (m *mockFieldType) DecimalSize() (precision, scale int64, ok bool) {
	return
}

func (m *mockFieldType) Nullable() (nullable, ok bool) {
	return
}

func (m *mockFieldType) DatabaseTypeName() string {
	return m.name
}

func (m *mockFieldType) IsSupported() bool {
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
			f:    NewField(database.NewBaseField(0, "f1", database.NewBaseFieldType(&sql.ColumnType{}))),
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
		in0 int
	}
	tests := []struct {
		name string
		f    *Field
		args args
		want string
	}{
		{
			name: "1",
			f:    NewField(database.NewBaseField(0, "f1", database.NewBaseFieldType(&sql.ColumnType{}))),
			args: args{
				in0: 0,
			},
			want: `?`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.BindVar(tt.args.in0); got != tt.want {
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
			f:    NewField(database.NewBaseField(0, "f1", database.NewBaseFieldType(&sql.ColumnType{}))),
			want: `"f1"`,
		},
	}
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
			f:    NewField(database.NewBaseField(0, "f1", newMockFieldType("DECIMAL"))),
			want: newMockFieldType("DECIMAL"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Type(); !reflect.DeepEqual(got.DatabaseTypeName(), tt.want.DatabaseTypeName()) {
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
			f:    NewField(database.NewBaseField(0, "f1", newMockFieldType("DECIMAL"))),
			want: NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("DECIMAL")))),
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
			f:    NewField(database.NewBaseField(0, "f1", newMockFieldType("BIGINT"))),
			args: args{
				c: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(int64(123)), "f1", 0),
			},
			want: database.NewGoValuer(NewField(database.NewBaseField(0, "f1", newMockFieldType("BIGINT"))), element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(int64(123)), "f1", 0)),
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
		//"BOOLEAN"
		{
			name: "BOOLEAN",
			f:    NewFieldType(newMockFieldType("BOOLEAN")),
			want: true,
		},
		//"BIGINT", "INTEGER", "SMALLINT"
		{
			name: "BIGINT",
			f:    NewFieldType(newMockFieldType("BIGINT")),
			want: true,
		},
		{
			name: "INTEGER",
			f:    NewFieldType(newMockFieldType("INTEGER")),
			want: true,
		},
		{
			name: "SMALLINT",
			f:    NewFieldType(newMockFieldType("SMALLINT")),
			want: true,
		},
		//"DOUBLE", "REAL"
		{
			name: "DOUBLE",
			f:    NewFieldType(newMockFieldType("DOUBLE")),
			want: true,
		},
		{
			name: "REAL",
			f:    NewFieldType(newMockFieldType("REAL")),
			want: true,
		},
		//"DATE", "TIME", "TIMESTAMP"
		{
			name: "DATE",
			f:    NewFieldType(newMockFieldType("DATE")),
			want: true,
		},
		{
			name: "TIME",
			f:    NewFieldType(newMockFieldType("TIME")),
			want: true,
		},
		{
			name: "TIMESTAMP",
			f:    NewFieldType(newMockFieldType("TIMESTAMP")),
			want: true,
		},
		//"VARCHAR", "CHAR", "DECIMAL"
		{
			name: "VARCHAR",
			f:    NewFieldType(newMockFieldType("VARCHAR")),
			want: true,
		},
		{
			name: "CHAR",
			f:    NewFieldType(newMockFieldType("CHAR")),
			want: true,
		},
		{
			name: "DECIMAL",
			f:    NewFieldType(newMockFieldType("DECIMAL")),
			want: true,
		},
		//"BLOB" "CLOB"
		{
			name: "BLOB",
			f:    NewFieldType(newMockFieldType("BLOB")),
			want: true,
		},
		{
			name: "CLOB",
			f:    NewFieldType(newMockFieldType("CLOB")),
			want: true,
		},
		//unknown
		{
			name: "unknown",
			f:    NewFieldType(&sql.ColumnType{}),
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

func TestFieldType_GoType(t *testing.T) {
	tests := []struct {
		name string
		f    *FieldType
		want database.GoType
	}{
		//"BOOLEAN"
		{
			name: "BOOLEAN",
			f:    NewFieldType(newMockFieldType("BOOLEAN")),
			want: database.GoTypeBool,
		},
		//"BIGINT", "INTEGER", "SMALLINT"
		{
			name: "BIGINT",
			f:    NewFieldType(newMockFieldType("BIGINT")),
			want: database.GoTypeInt64,
		},
		{
			name: "INTEGER",
			f:    NewFieldType(newMockFieldType("INTEGER")),
			want: database.GoTypeInt64,
		},
		{
			name: "SMALLINT",
			f:    NewFieldType(newMockFieldType("SMALLINT")),
			want: database.GoTypeInt64,
		},
		//"DOUBLE", "REAL"
		{
			name: "DOUBLE",
			f:    NewFieldType(newMockFieldType("DOUBLE")),
			want: database.GoTypeFloat64,
		},
		{
			name: "REAL",
			f:    NewFieldType(newMockFieldType("REAL")),
			want: database.GoTypeFloat64,
		},
		//"DATE", "TIME", "TIMESTAMP"
		{
			name: "DATE",
			f:    NewFieldType(newMockFieldType("DATE")),
			want: database.GoTypeTime,
		},
		{
			name: "TIME",
			f:    NewFieldType(newMockFieldType("TIME")),
			want: database.GoTypeTime,
		},
		{
			name: "TIMESTAMP",
			f:    NewFieldType(newMockFieldType("TIMESTAMP")),
			want: database.GoTypeTime,
		},
		//"VARCHAR", "CHAR", "DECIMAL"
		{
			name: "VARCHAR",
			f:    NewFieldType(newMockFieldType("VARCHAR")),
			want: database.GoTypeString,
		},
		{
			name: "CHAR",
			f:    NewFieldType(newMockFieldType("CHAR")),
			want: database.GoTypeString,
		},
		{
			name: "DECIMAL",
			f:    NewFieldType(newMockFieldType("DECIMAL")),
			want: database.GoTypeString,
		},
		//"BLOB" "CLOB"
		{
			name: "BLOB",
			f:    NewFieldType(newMockFieldType("BLOB")),
			want: database.GoTypeBytes,
		},
		{
			name: "CLOB",
			f:    NewFieldType(newMockFieldType("CLOB")),
			want: database.GoTypeBytes,
		},
		//unknown
		{
			name: "unknown",
			f:    NewFieldType(&sql.ColumnType{}),
			want: database.GoTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.GoType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FieldType.GoType() = %v, want %v", got, tt.want)
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
		wantErr bool
		want    element.Column
	}{
		//"BOOLEAN"
		{
			name: "BOOLEAN",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("BOOLEAN")))),
			args: args{
				src: true,
			},
			want: element.NewDefaultColumn(element.NewBoolColumnValue(true), "test", 1),
		},
		{
			name: "BOOLEAN error",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("BOOLEAN")))),
			args: args{
				src: 1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "BOOLEAN nil",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("BOOLEAN")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBoolColumnValue(), "test", 0),
		},
		//"BIGINT", "INTEGER", "SMALLINT"
		{
			name: "BIGINT",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("BIGINT")))),
			args: args{
				src: int64(1),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1), "test", element.ByteSize(int64(1))),
		},
		{
			name: "INTEGER",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("INTEGER")))),
			args: args{
				src: int32(1),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1), "test", element.ByteSize(int32(1))),
		},
		{
			name: "SMALLINT",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("SMALLINT")))),
			args: args{
				src: int16(1),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1), "test", element.ByteSize(int16(1))),
		},
		{
			name: "SMALLINT nil",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("SMALLINT")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBigIntColumnValue(), "test", 0),
		},
		{
			name: "SMALLINT error",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("SMALLINT")))),
			args: args{
				src: "",
			},
			want:    nil,
			wantErr: true,
		},
		//"DOUBLE", "REAL","DECIMAL"
		{
			name: "DOUBLE",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DOUBLE")))),
			args: args{
				src: 1.01,
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1.01), "test", element.ByteSize(1.01)),
		},
		{
			name: "REAL nil",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("REAL")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilDecimalColumnValue(), "test", 0),
		},
		{
			name: "REAL error",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("REAL")))),
			args: args{
				src: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "DECIMAL",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DECIMAL")))),
			args: args{
				src: []byte("1.01"),
			},
			want: element.NewDefaultColumn(mustDecimalColumnValueFromString("1.01"), "test", element.ByteSize([]byte("1.01"))),
		},
		{
			name: "DECIMAL error",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DECIMAL")))),
			args: args{
				src: []byte("1.01a"),
			},
			want:    nil,
			wantErr: true,
		},
		//"BLOB"
		{
			name: "BLOB",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("BLOB")))),
			args: args{
				src: []byte("1.01a"),
			},
			want: element.NewDefaultColumn(element.NewBytesColumnValue([]byte("1.01a")),
				"test", element.ByteSize([]byte("1.01a"))),
		},
		{
			name: "CLOB",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("CLOB")))),
			args: args{
				src: []byte("中文abc"),
			},
			want: element.NewDefaultColumn(element.NewBytesColumnValue([]byte("中文abc")),
				"test", element.ByteSize([]byte("中文abc"))),
		},
		{
			name: "BLOB nil",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("BLOB")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBytesColumnValue(), "test", 0),
		},
		{
			name: "BLOB error",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("BLOB")))),
			args: args{
				src: "",
			},
			want:    nil,
			wantErr: true,
		},
		//"DATE"
		{
			name: "DATE",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DATE")))),
			args: args{
				src: time.Date(2022, 5, 1, 0, 0, 0, 0, time.Local),
			},
			want: element.NewDefaultColumn(element.NewTimeColumnValueWithDecoder(
				time.Date(2022, 5, 1, 0, 0, 0, 0, time.Local),
				element.NewStringTimeDecoder(dateLayout)), "test", element.ByteSize(time.Date(2022, 5, 1, 0, 0, 0, 0, time.Local))),
		},
		{
			name: "DATE nil",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DATE")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilTimeColumnValue(), "test", 0),
		},
		{
			name: "DATE error",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DATE")))),
			args: args{
				src: "",
			},
			want:    nil,
			wantErr: true,
		},
		//"TIME"
		{
			name: "TIME",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("TIME")))),
			args: args{
				src: time.Date(2022, 5, 1, 14, 57, 11, 111, time.Local),
			},
			want: element.NewDefaultColumn(element.NewTimeColumnValueWithDecoder(
				time.Date(2022, 5, 1, 14, 57, 11, 111, time.Local),
				element.NewStringTimeDecoder(timeLayout)), "test", element.ByteSize(time.Date(2022, 5, 1, 0, 0, 0, 0, time.Local))),
		},
		{
			name: "TIME nil",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("TIME")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilTimeColumnValue(), "test", 0),
		},
		{
			name: "TIME error",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("TIME")))),
			args: args{
				src: "",
			},
			want:    nil,
			wantErr: true,
		},
		//"TIMESTAMP"
		{
			name: "TIMESTAMP",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("TIMESTAMP")))),
			args: args{
				src: time.Date(2022, 5, 1, 14, 57, 11, 111, time.Local),
			},
			want: element.NewDefaultColumn(element.NewTimeColumnValueWithDecoder(
				time.Date(2022, 5, 1, 14, 57, 11, 111, time.Local),
				element.NewStringTimeDecoder(timestampLayout)), "test", element.ByteSize(time.Date(2022, 5, 1, 0, 0, 0, 0, time.Local))),
		},
		{
			name: "TIMESTAMP nil",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("TIMESTAMP")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilTimeColumnValue(), "test", 0),
		},
		{
			name: "TIMESTAMP error",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("TIMESTAMP")))),
			args: args{
				src: "",
			},
			want:    nil,
			wantErr: true,
		},
		//"CHAR"  "VARCHAR"
		{
			name: "CHAR",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("CHAR")))),
			args: args{
				src: []byte("abc"),
			},
			want: element.NewDefaultColumn(element.NewStringColumnValue("abc"), "test", element.ByteSize([]byte("abc"))),
		},
		{
			name: "CHARTrim",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("CHAR")))),
			conf: testJSONFromString(`{"trimChar":true}`),
			args: args{
				src: []byte("    abc   "),
			},
			want: element.NewDefaultColumn(element.NewStringColumnValue("abc"), "test", element.ByteSize([]byte("    abc   "))),
		},
		{
			name: "CHAR nil",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("CHAR")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilStringColumnValue(), "test", 0),
		},
		{
			name: "CHAR error",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("CHAR")))),
			args: args{
				src: 1,
			},
			want:    nil,
			wantErr: true,
		},
		//unknown
		{
			name: "unknown",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("")))),
			args: args{
				src: 1,
			},
			want:    nil,
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
				t.Errorf("Scanner.Column() = %v, want %v", tt.s.Column(), tt.want)
			}
		})
	}
}
