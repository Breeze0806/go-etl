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

package dm

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
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

func TestField_Quoted(t *testing.T) {
	tests := []struct {
		name string
		f    *Field
		want string
	}{
		{
			name: "1",
			f:    NewField(database.NewBaseField(0, "f1", NewFieldType(&sql.ColumnType{}))),
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
			f:    NewField(database.NewBaseField(0, "f1", NewFieldType(&sql.ColumnType{}))),
			args: args{
				i: 22,
			},
			want: "?",
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
			f:    NewField(database.NewBaseField(0, "f1", NewFieldType(&sql.ColumnType{}))),
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
			f:    NewField(database.NewBaseField(0, "f1", NewFieldType(newMockColumnType("1")))),
			want: NewFieldType(&mockColumnType{
				name: "1",
			}),
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
			f:    NewField(database.NewBaseField(0, "f1", NewFieldType(newMockColumnType("1")))),
			want: NewScanner(NewField(database.NewBaseField(0, "f1", NewFieldType(newMockColumnType("1"))))),
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
			f:    NewField(database.NewBaseField(0, "f1", NewFieldType(newMockColumnType("1")))),
			args: args{
				c: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1), "f1", 0),
			},
			want: database.NewGoValuer(NewField(database.NewBaseField(0, "f1", NewFieldType(newMockColumnType("1")))), element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1), "f1", 0)),
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

func TestFieldType_GoType(t *testing.T) {
	tests := []struct {
		name string
		f    *FieldType
		want database.GoType
	}{
		// BOOLEAN 类型
		{
			name: "BIT",
			f:    NewFieldType(newMockColumnType("BIT")),
			want: database.GoTypeBool,
		},
		{
			name: "BOOLEAN",
			f:    NewFieldType(newMockColumnType("BOOLEAN")),
			want: database.GoTypeBool,
		},
		{
			name: "BOOL",
			f:    NewFieldType(newMockColumnType("BOOL")),
			want: database.GoTypeBool,
		},

		// INTEGER 类型
		{
			name: "INTEGER",
			f:    NewFieldType(newMockColumnType("INTEGER")),
			want: database.GoTypeInt64,
		},
		{
			name: "INT",
			f:    NewFieldType(newMockColumnType("INT")),
			want: database.GoTypeInt64,
		},
		{
			name: "BIGINT",
			f:    NewFieldType(newMockColumnType("BIGINT")),
			want: database.GoTypeInt64,
		},
		{
			name: "TINYINT",
			f:    NewFieldType(newMockColumnType("TINYINT")),
			want: database.GoTypeInt64,
		},
		{
			name: "SMALLINT",
			f:    NewFieldType(newMockColumnType("SMALLINT")),
			want: database.GoTypeInt64,
		},
		{
			name: "BYTE",
			f:    NewFieldType(newMockColumnType("BYTE")),
			want: database.GoTypeInt64,
		},

		// NUMERIC/DECIMAL 类型
		{
			name: "NUMERIC",
			f:    NewFieldType(newMockColumnType("NUMERIC")),
			want: database.GoTypeString,
		},
		{
			name: "NUMBER",
			f:    NewFieldType(newMockColumnType("NUMBER")),
			want: database.GoTypeString,
		},
		{
			name: "DECIMAL",
			f:    NewFieldType(newMockColumnType("DECIMAL")),
			want: database.GoTypeString,
		},
		{
			name: "DEC",
			f:    NewFieldType(newMockColumnType("DEC")),
			want: database.GoTypeString,
		},
		{
			name: "FLOAT",
			f:    NewFieldType(newMockColumnType("FLOAT")),
			want: database.GoTypeString,
		},
		{
			name: "DOUBLE",
			f:    NewFieldType(newMockColumnType("DOUBLE")),
			want: database.GoTypeString,
		},
		{
			name: "REAL",
			f:    NewFieldType(newMockColumnType("REAL")),
			want: database.GoTypeString,
		},
		{
			name: "DOUBLE PRECISION",
			f:    NewFieldType(newMockColumnType("DOUBLE PRECISION")),
			want: database.GoTypeString,
		},

		// STRING 类型
		{
			name: "CHAR",
			f:    NewFieldType(newMockColumnType("CHAR")),
			want: database.GoTypeString,
		},
		{
			name: "CHARACTER",
			f:    NewFieldType(newMockColumnType("CHARACTER")),
			want: database.GoTypeString,
		},
		{
			name: "VARCHAR",
			f:    NewFieldType(newMockColumnType("VARCHAR")),
			want: database.GoTypeString,
		},
		{
			name: "TEXT",
			f:    NewFieldType(newMockColumnType("TEXT")),
			want: database.GoTypeString,
		},
		{
			name: "CLOB",
			f:    NewFieldType(newMockColumnType("CLOB")),
			want: database.GoTypeString,
		},
		{
			name: "LONGVARCHAR",
			f:    NewFieldType(newMockColumnType("LONGVARCHAR")),
			want: database.GoTypeString,
		},

		// BINARY 类型
		{
			name: "BINARY",
			f:    NewFieldType(newMockColumnType("BINARY")),
			want: database.GoTypeBytes,
		},
		{
			name: "VARBINARY",
			f:    NewFieldType(newMockColumnType("VARBINARY")),
			want: database.GoTypeBytes,
		},
		{
			name: "BLOB",
			f:    NewFieldType(newMockColumnType("BLOB")),
			want: database.GoTypeBytes,
		},
		{
			name: "BFILE",
			f:    NewFieldType(newMockColumnType("BFILE")),
			want: database.GoTypeBytes,
		},
		{
			name: "IMAGE",
			f:    NewFieldType(newMockColumnType("IMAGE")),
			want: database.GoTypeBytes,
		},
		{
			name: "LONGVARBINARY",
			f:    NewFieldType(newMockColumnType("LONGVARBINARY")),
			want: database.GoTypeBytes,
		},

		// TIME 类型
		{
			name: "DATE",
			f:    NewFieldType(newMockColumnType("DATE")),
			want: database.GoTypeTime,
		},
		{
			name: "TIME",
			f:    NewFieldType(newMockColumnType("TIME")),
			want: database.GoTypeTime,
		},
		{
			name: "DATETIME",
			f:    NewFieldType(newMockColumnType("DATETIME")),
			want: database.GoTypeTime,
		},
		{
			name: "TIMESTAMP",
			f:    NewFieldType(newMockColumnType("TIMESTAMP")),
			want: database.GoTypeTime,
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

func TestFieldType_IsSupportted(t *testing.T) {
	tests := []struct {
		name string
		f    *FieldType
		want bool
	}{
		{
			name: "1",
			f:    NewFieldType(newMockColumnType("INTEGER")),
			want: true,
		},
		{
			name: "2",
			f:    NewFieldType(newMockColumnType("BLOB")),
			want: true,
		},
		{
			name: "3",
			f:    NewFieldType(newMockColumnType("NUMERIC")),
			want: true,
		},
		{
			name: "4",
			f:    NewFieldType(newMockColumnType("REAL")),
			want: true,
		},
		{
			name: "5",
			f:    NewFieldType(newMockColumnType("TEXT")),
			want: true,
		},
		{
			name: "6",
			f:    NewFieldType(newMockColumnType("TEXT1")),
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
		// BOOLEAN 类型测试
		{
			name: "BOOLEAN-nil",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("BOOLEAN"))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBoolColumnValue(), "f1", 0),
		},
		{
			name: "BOOLEAN-bool",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("BOOLEAN"))))),
			args: args{
				src: true,
			},
			want: element.NewDefaultColumn(element.NewBoolColumnValue(true), "f1", element.ByteSize(true)),
		},
		{
			name: "BOOLEAN-int8",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("BOOLEAN"))))),
			args: args{
				src: int8(1),
			},
			want: element.NewDefaultColumn(element.NewBoolColumnValue(true), "f1", element.ByteSize(int8(1))),
		},
		{
			name: "BOOLEAN-invalid",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("BOOLEAN"))))),
			args: args{
				src: "invalid",
			},
			wantErr: true,
		},

		// INTEGER 类型测试
		{
			name: "INT-nil",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("INT"))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBigIntColumnValue(), "f1", 0),
		},
		{
			name: "INT-int64",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("INT"))))),
			args: args{
				src: int64(9223372036854775807), // 2^63-1
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(int64(9223372036854775807)), "f1", element.ByteSize(int64(9223372036854775807))),
		},
		{
			name: "INT-int32",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("INT"))))),
			args: args{
				src: int32(12345),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(int64(12345)), "f1", element.ByteSize(int32(12345))),
		},
		{
			name: "INT-string",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("INT"))))),
			args: args{
				src: "123",
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(123), "f1", element.ByteSize("123")),
		},
		{
			name: "INT-invalid",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("INT"))))),
			args: args{
				src: "invalid",
			},
			wantErr: true,
		},

		// BINARY 类型测试
		{
			name: "BLOB-nil",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("BLOB"))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBytesColumnValue(), "f1", 0),
		},
		{
			name: "BLOB-[]byte",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("BLOB"))))),
			args: args{
				src: []byte("123"),
			},
			want: element.NewDefaultColumn(element.NewBytesColumnValueNoCopy([]byte("123")), "f1", element.ByteSize([]byte("123"))),
		},
		{
			name: "BLOB-invalid",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("BLOB"))))),
			args: args{
				src: "invalid",
			},
			wantErr: true,
		},

		// NUMERIC 类型测试
		{
			name: "NUMERIC-nil",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("NUMERIC"))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilDecimalColumnValue(), "f1", 0),
		},
		{
			name: "NUMERIC-float64",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("NUMERIC"))))),
			args: args{
				src: 1.23456789,
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1.23456789), "f1", element.ByteSize(1.23456789)),
		},
		{
			name: "NUMERIC-float32",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("NUMERIC"))))),
			args: args{
				src: float32(1.234),
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat32(float32(1.234)), "f1", element.ByteSize(float32(1.234))),
		},
		{
			name: "NUMERIC-string",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("NUMERIC"))))),
			args: args{
				src: "1.234",
			},
			want: func() element.Column {
				d, _ := element.NewDecimalColumnValueFromString("1.234")
				return element.NewDefaultColumn(d, "f1", element.ByteSize([]byte("1.234")))
			}(),
		},
		{
			name: "NUMERIC-[]byte",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("NUMERIC"))))),
			args: args{
				src: []byte("1.234"),
			},
			want: func() element.Column {
				d, _ := element.NewDecimalColumnValueFromString("1.234")
				return element.NewDefaultColumn(d, "f1", element.ByteSize([]byte("1.234")))
			}(),
		},
		{
			name: "NUMERIC-invalid",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("NUMERIC"))))),
			args: args{
				src: "invalid",
			},
			wantErr: true,
		},

		// REAL 类型测试
		{
			name: "REAL-nil",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("REAL"))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilDecimalColumnValue(), "f1", 0),
		},
		{
			name: "REAL-float32",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("REAL"))))),
			args: args{
				src: float32(1.234),
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat32(float32(1.234)), "f1", element.ByteSize(float32(1.234))),
		},
		{
			name: "REAL-float64",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("REAL"))))),
			args: args{
				src: 1.23456789,
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1.23456789), "f1", element.ByteSize(1.23456789)),
		},
		{
			name: "REAL-invalid",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("REAL"))))),
			args: args{
				src: "invalid",
			},
			wantErr: true,
		},

		// STRING 类型测试
		{
			name: "TEXT-nil",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("TEXT"))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilStringColumnValue(), "f1", 0),
		},
		{
			name: "TEXT-string",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("TEXT"))))),
			args: args{
				src: "123",
			},
			want: element.NewDefaultColumn(element.NewStringColumnValue("123"), "f1", element.ByteSize("123")),
		},
		{
			name: "TEXT-invalid",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("TEXT"))))),
			args: args{
				src: 123,
			},
			wantErr: true,
		},

		// DATE 类型测试
		{
			name: "DATE-nil",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("DATE"))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilTimeColumnValue(), "f1", 0),
		},
		{
			name: "DATE-time",
			s: func() *Scanner {
				scanner := NewScanner(NewField(database.NewBaseField(0,
					"f1", NewFieldType(newMockColumnType("DATE")))))
				return scanner
			}(),
			args: args{
				src: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: element.NewDefaultColumn(
				element.NewTimeColumnValueWithDecoder(
					time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					element.NewStringTimeDecoder(dateLayout),
				),
				"f1",
				element.ByteSize(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
			),
		},
		{
			name: "DATE-invalid",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("DATE"))))),
			args: args{
				src: "invalid",
			},
			wantErr: true,
		},

		// DATETIME/TIMESTAMP 类型测试
		{
			name: "DATETIME-nil",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("DATETIME"))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilTimeColumnValue(), "f1", 0),
		},
		{
			name: "TIMESTAMP-time",
			s: func() *Scanner {
				scanner := NewScanner(NewField(database.NewBaseField(0,
					"f1", NewFieldType(newMockColumnType("TIMESTAMP")))))
				return scanner
			}(),
			args: args{
				src: time.Date(2023, 1, 1, 12, 30, 45, 123456789, time.UTC),
			},
			want: element.NewDefaultColumn(
				element.NewTimeColumnValueWithDecoder(
					time.Date(2023, 1, 1, 12, 30, 45, 123456789, time.UTC),
					element.NewStringTimeDecoder(datetimeLayout),
				),
				"f1",
				element.ByteSize(time.Date(2023, 1, 1, 12, 30, 45, 123456789, time.UTC)),
			),
		},
		{
			name: "DATETIME-invalid",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("DATETIME"))))),
			args: args{
				src: "invalid",
			},
			wantErr: true,
		},

		// 无效数据库类型测试
		{
			name: "INVALID-type",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("INVALID"))))),
			args: args{
				src: 123,
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

			if !tt.wantErr {
				gotColumn := tt.s.Column()
				if !reflect.DeepEqual(gotColumn, tt.want) {
					t.Errorf("Column() = %v, want %v", gotColumn, tt.want)
				}
			}
		})
	}
}
