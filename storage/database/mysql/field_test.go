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

package mysql

import (
	"database/sql"
	"math"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/shopspring/decimal"
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

func mustBigIntValueFromString(s string) element.ColumnValue {
	c, err := element.NewBigIntColumnValueFromString(s)
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
			f:    NewField(database.NewBaseField(0, "table", database.NewBaseFieldType(&sql.ColumnType{}))),
			want: "`table`",
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
			f:    NewField(database.NewBaseField(0, "table", database.NewBaseFieldType(&sql.ColumnType{}))),
			args: args{
				i: 0,
			},
			want: "?",
		},
		{
			name: "1",
			f:    NewField(database.NewBaseField(0, "table", database.NewBaseFieldType(&sql.ColumnType{}))),
			args: args{
				i: 100000,
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
			f:    NewField(database.NewBaseField(0, "table", database.NewBaseFieldType(&sql.ColumnType{}))),
			want: "`table`",
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
			f:    NewField(database.NewBaseField(0, "table", database.NewBaseFieldType(&sql.ColumnType{}))),
			want: NewFieldType(&sql.ColumnType{}),
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
			f:    NewField(database.NewBaseField(0, "table", database.NewBaseFieldType(&sql.ColumnType{}))),
			want: NewScanner(NewField(database.NewBaseField(0, "table", database.NewBaseFieldType(&sql.ColumnType{})))),
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
			f:    NewField(database.NewBaseField(0, "f1", NewFieldType(&sql.ColumnType{}))),
			args: args{
				c: element.NewDefaultColumn(nil, "", 0),
			},
			want: database.NewGoValuer(NewField(database.NewBaseField(0, "f1", NewFieldType(&sql.ColumnType{}))), element.NewDefaultColumn(nil, "", 0)),
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
		// "MEDIUMINT", "INT", "BIGINT", "SMALLINT", "TINYINT",
		// "TEXT", "LONGTEXT", "MEDIUMTEXT", "TINYTEXT", "CHAR", "VARCHAR",
		// "TIME", "YEAR",
		// "DECIMAL"
		{
			name: "MEDIUMINT",
			f:    NewFieldType(newMockFieldType("MEDIUMINT")),
			want: database.GoTypeString,
		},
		{
			name: "INT",
			f:    NewFieldType(newMockFieldType("INT")),
			want: database.GoTypeString,
		},
		{
			name: "BIGINT",
			f:    NewFieldType(newMockFieldType("BIGINT")),
			want: database.GoTypeString,
		},
		{
			name: "SMALLINT",
			f:    NewFieldType(newMockFieldType("SMALLINT")),
			want: database.GoTypeString,
		},
		{
			name: "TINYINT",
			f:    NewFieldType(newMockFieldType("TINYINT")),
			want: database.GoTypeString,
		},
		{
			name: "TEXT",
			f:    NewFieldType(newMockFieldType("TEXT")),
			want: database.GoTypeString,
		},
		{
			name: "LONGTEXT",
			f:    NewFieldType(newMockFieldType("LONGTEXT")),
			want: database.GoTypeString,
		},
		{
			name: "MEDIUMTEXT",
			f:    NewFieldType(newMockFieldType("MEDIUMTEXT")),
			want: database.GoTypeString,
		},
		{
			name: "TINYTEXT",
			f:    NewFieldType(newMockFieldType("TINYTEXT")),
			want: database.GoTypeString,
		},
		{
			name: "CHAR",
			f:    NewFieldType(newMockFieldType("CHAR")),
			want: database.GoTypeString,
		},
		{
			name: "VARCHAR",
			f:    NewFieldType(newMockFieldType("VARCHAR")),
			want: database.GoTypeString,
		},
		{
			name: "TIME",
			f:    NewFieldType(newMockFieldType("TIME")),
			want: database.GoTypeString,
		},
		{
			name: "YEAR",
			f:    NewFieldType(newMockFieldType("YEAR")),
			want: database.GoTypeString,
		},
		{
			name: "DECIMAL",
			f:    NewFieldType(newMockFieldType("DECIMAL")),
			want: database.GoTypeString,
		},
		//"BLOB", "LONGBLOB", "MEDIUMBLOB", "BINARY", "TINYBLOB", "VARBINARY"
		{
			name: "BLOB",
			f:    NewFieldType(newMockFieldType("BLOB")),
			want: database.GoTypeBytes,
		},
		{
			name: "LONGBLOB",
			f:    NewFieldType(newMockFieldType("LONGBLOB")),
			want: database.GoTypeBytes,
		},
		{
			name: "MEDIUMBLOB",
			f:    NewFieldType(newMockFieldType("MEDIUMBLOB")),
			want: database.GoTypeBytes,
		},
		{
			name: "BINARY",
			f:    NewFieldType(newMockFieldType("BINARY")),
			want: database.GoTypeBytes,
		},
		{
			name: "TINYBLOB",
			f:    NewFieldType(newMockFieldType("TINYBLOB")),
			want: database.GoTypeBytes,
		},
		{
			name: "VARBINARY",
			f:    NewFieldType(newMockFieldType("VARBINARY")),
			want: database.GoTypeBytes,
		},
		//"DOUBLE", "FLOAT"
		{
			name: "FLOAT",
			f:    NewFieldType(newMockFieldType("FLOAT")),
			want: database.GoTypeFloat64,
		},
		{
			name: "DOUBLE",
			f:    NewFieldType(newMockFieldType("DOUBLE")),
			want: database.GoTypeFloat64,
		},
		//"DATE", "DATETIME", "TIMESTAMP"
		{
			name: "DATE",
			f:    NewFieldType(newMockFieldType("DATE")),
			want: database.GoTypeTime,
		},
		{
			name: "DATETIME",
			f:    NewFieldType(newMockFieldType("DATETIME")),
			want: database.GoTypeTime,
		},
		{
			name: "TIMESTAMP",
			f:    NewFieldType(newMockFieldType("TIMESTAMP")),
			want: database.GoTypeTime,
		},
		{
			name: "NEWDATE",
			f:    NewFieldType(newMockFieldType("NEWDATE")),
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
		src interface{}
	}
	tests := []struct {
		name    string
		s       *Scanner
		conf    *config.JSON
		args    args
		wantErr bool
		want    element.Column
	}{
		//"MEDIUMINT", "INT", "BIGINT", "SMALLINT", "TINYINT", "YEAR"，"UNSIGNED INT", "UNSIGNED BIGINT", "UNSIGNED SMALLINT", "UNSIGNED TINYINT"
		{
			name: "BIGINT",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("BIGINT")))),
			args: args{
				src: []byte("123123456789"),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(123123456789), "test", element.ByteSize([]byte("123123456789"))),
		},
		{
			name: "MEDIUMINT",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("MEDIUMINT")))),
			args: args{
				src: []byte("123123456789e"),
			},
			wantErr: true,
		},
		{
			name: "YEAR",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("YEAR")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBigIntColumnValue(), "test", 0),
		},
		{
			name: "TINYINT",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("TINYINT")))),
			args: args{
				src: int64(123),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(123), "test",
				element.ByteSize(int64(123))),
		},
		{
			name: "SMALLINT",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("SMALLINT")))),
			args: args{
				src: int32(123),
			},
			wantErr: true,
		},
		{
			name: "UNSIGNED BIGINT",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("TINYINT")))),
			args: args{
				src: []byte(strconv.FormatUint(math.MaxUint64, 10)),
			},
			want: element.NewDefaultColumn(mustBigIntValueFromString(strconv.FormatUint(math.MaxUint64, 10)), "test", element.ByteSize(strconv.FormatUint(math.MaxUint64, 10))),
		},
		{
			name: "UNSIGNED INT",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("UNSIGNED INT")))),
			args: args{
				src: int64(math.MaxInt32),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(int64(math.MaxInt32)), "test", element.ByteSize(int64(math.MaxInt32))),
		},
		//"BLOB", "LONGBLOB", "MEDIUMBLOB", "BINARY", "TINYBLOB", "VARBINARY"
		{
			name: "BLOB",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("BLOB")))),
			args: args{
				src: []byte("123123456789"),
			},
			want: element.NewDefaultColumn(element.NewBytesColumnValue([]byte("123123456789")), "test", element.ByteSize([]byte("123123456789"))),
		},
		{
			name: "BINARY",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("BINARY")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBytesColumnValue(), "test", 0),
		},
		{
			name: "VARBINARY",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("BINARY")))),
			args: args{
				src: "nil",
			},
			wantErr: true,
		},
		//"DATE", "DATETIME", "TIMESTAMP"
		{
			name: "DATE",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DATE")))),
			args: args{
				src: time.Date(2021, 1, 13, 18, 43, 12, 0, time.Local),
			},
			want: element.NewDefaultColumn(
				element.NewTimeColumnValueWithDecoder(time.Date(2021, 1, 13, 18, 43, 12, 0, time.Local), element.NewStringTimeDecoder(dateLayout)),
				"test", element.ByteSize(time.Date(2021, 1, 13, 18, 43, 12, 0, time.Local))),
		},
		{
			name: "DATEnull",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DATE")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilTimeColumnValue(), "test", 0),
		},
		{
			name: "DATEerr",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DATE")))),
			args: args{
				src: "123",
			},
			wantErr: true,
		},
		{
			name: "DATETIME",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DATETIME")))),
			args: args{
				src: time.Date(2021, 1, 13, 18, 43, 12, 0, time.Local),
			},
			want: element.NewDefaultColumn(element.NewTimeColumnValueWithDecoder(time.Date(2021, 1, 13, 18, 43, 12, 0, time.Local), element.NewStringTimeDecoder(datetimeLayout)),
				"test", element.ByteSize(time.Date(2021, 1, 13, 18, 43, 12, 0, time.Local))),
		},
		{
			name: "TIMESTAMPErr",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("TIMESTAMP")))),
			args: args{
				src: "nil",
			},
			wantErr: true,
		},
		{
			name: "TIMESTAMPnil",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("TIMESTAMP")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilTimeColumnValue(), "test", 0),
		},
		//"TEXT", "LONGTEXT", "MEDIUMTEXT", "TINYTEXT", "CHAR", "VARCHAR", "TIME"
		{
			name: "TEXT",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("TEXT")))),
			args: args{
				src: []byte("中文abc%$`\""),
			},
			want: element.NewDefaultColumn(element.NewStringColumnValue("中文abc%$`\""), "test",
				element.ByteSize([]byte("中文abc%$`\""))),
		},
		{
			name: "CHARnil",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("CHAR")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilStringColumnValue(), "test", 0),
		},
		{
			name: "CHARTrim",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("CHAR")))),
			args: args{
				src: []byte("   中文abc%$`\"  "),
			},
			conf: testJSONFromString(`{"trimChar":true}`),
			want: element.NewDefaultColumn(element.NewStringColumnValue("中文abc%$`\""), "test",
				element.ByteSize([]byte("   中文abc%$`\"  "))),
		},
		{
			name: "CHAR",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("CHAR")))),
			args: args{
				src: []byte("   中文abc%$`\"  "),
			},
			want: element.NewDefaultColumn(element.NewStringColumnValue("   中文abc%$`\"  "), "test",
				element.ByteSize([]byte("   中文abc%$`\"  "))),
		},
		{
			name: "TIME",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("TIME")))),
			args: args{
				src: int16(0),
			},
			wantErr: true,
		},
		//"DOUBLE", "FLOAT", "DECIMAL"
		{
			name: "DOUBLE",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DOUBLE")))),
			args: args{
				src: []byte("123456.7123456"),
			},
			want: element.NewDefaultColumn(mustDecimalColumnValueFromString("123456.7123456"),
				"test", element.ByteSize([]byte("123456.7123456"))),
		},
		{
			name: "DOUBLE",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DOUBLE")))),
			args: args{
				src: []byte("123456.7123456e"),
			},
			wantErr: true,
		},
		{
			name: "FLOAT",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("FLOAT")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilDecimalColumnValue(), "test", 0),
		},
		{
			name: "FLOAT",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("FLOAT")))),
			args: args{
				src: float32(1.234),
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValue(decimal.NewFromFloat32(float32(1.234))),
				"test", element.ByteSize(float32(1.234))),
		},
		{
			name: "DOUBLE",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("FLOAT")))),
			args: args{
				src: float64(1.23456789),
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(float64(1.23456789)),
				"test", element.ByteSize(float32(1.23456789))),
		},
		{
			name: "DECIMAL",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DECIMAL")))),
			args: args{
				src: int16(0),
			},
			wantErr: true,
		},

		{
			name: "DECIMAL1",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DECIMAL1")))),
			args: args{
				src: int16(0),
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
				t.Errorf("Scanner.Column() = %v, want %v", tt.s.Column(), tt.want)
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
		// "MEDIUMINT", "INT", "BIGINT", "SMALLINT", "TINYINT",
		// "TEXT", "LONGTEXT", "MEDIUMTEXT", "TINYTEXT", "CHAR", "VARCHAR",
		// "TIME", "YEAR",
		// "DECIMAL"
		{
			name: "MEDIUMINT",
			f:    NewFieldType(newMockFieldType("MEDIUMINT")),
			want: true,
		},
		{
			name: "INT",
			f:    NewFieldType(newMockFieldType("INT")),
			want: true,
		},
		{
			name: "BIGINT",
			f:    NewFieldType(newMockFieldType("BIGINT")),
			want: true,
		},
		{
			name: "SMALLINT",
			f:    NewFieldType(newMockFieldType("SMALLINT")),
			want: true,
		},
		{
			name: "TINYINT",
			f:    NewFieldType(newMockFieldType("TINYINT")),
			want: true,
		},
		{
			name: "TEXT",
			f:    NewFieldType(newMockFieldType("TEXT")),
			want: true,
		},
		{
			name: "LONGTEXT",
			f:    NewFieldType(newMockFieldType("LONGTEXT")),
			want: true,
		},
		{
			name: "MEDIUMTEXT",
			f:    NewFieldType(newMockFieldType("MEDIUMTEXT")),
			want: true,
		},
		{
			name: "TINYTEXT",
			f:    NewFieldType(newMockFieldType("TINYTEXT")),
			want: true,
		},
		{
			name: "CHAR",
			f:    NewFieldType(newMockFieldType("CHAR")),
			want: true,
		},
		{
			name: "VARCHAR",
			f:    NewFieldType(newMockFieldType("VARCHAR")),
			want: true,
		},
		{
			name: "TIME",
			f:    NewFieldType(newMockFieldType("TIME")),
			want: true,
		},
		{
			name: "YEAR",
			f:    NewFieldType(newMockFieldType("YEAR")),
			want: true,
		},
		{
			name: "DECIMAL",
			f:    NewFieldType(newMockFieldType("DECIMAL")),
			want: true,
		},
		//"BLOB", "LONGBLOB", "MEDIUMBLOB", "BINARY", "TINYBLOB", "VARBINARY"
		{
			name: "BLOB",
			f:    NewFieldType(newMockFieldType("BLOB")),
			want: true,
		},
		{
			name: "LONGBLOB",
			f:    NewFieldType(newMockFieldType("LONGBLOB")),
			want: true,
		},
		{
			name: "MEDIUMBLOB",
			f:    NewFieldType(newMockFieldType("MEDIUMBLOB")),
			want: true,
		},
		{
			name: "BINARY",
			f:    NewFieldType(newMockFieldType("BINARY")),
			want: true,
		},
		{
			name: "TINYBLOB",
			f:    NewFieldType(newMockFieldType("TINYBLOB")),
			want: true,
		},
		{
			name: "VARBINARY",
			f:    NewFieldType(newMockFieldType("VARBINARY")),
			want: true,
		},
		//"DOUBLE", "FLOAT"
		{
			name: "FLOAT",
			f:    NewFieldType(newMockFieldType("FLOAT")),
			want: true,
		},
		{
			name: "DOUBLE",
			f:    NewFieldType(newMockFieldType("DOUBLE")),
			want: true,
		},
		//"DATE", "DATETIME", "TIMESTAMP"
		{
			name: "DATE",
			f:    NewFieldType(newMockFieldType("DATE")),
			want: true,
		},
		{
			name: "DATETIME",
			f:    NewFieldType(newMockFieldType("DATETIME")),
			want: true,
		},
		{
			name: "TIMESTAMP",
			f:    NewFieldType(newMockFieldType("TIMESTAMP")),
			want: true,
		},
		{
			name: "NEWDATE",
			f:    NewFieldType(newMockFieldType("NEWDATE")),
			want: false,
		},
		//"UNSIGNED INT", "UNSIGNED BIGINT", "UNSIGNED SMALLINT", "UNSIGNED TINYINT"
		{
			name: "UNSIGNED INT",
			f:    NewFieldType(newMockFieldType("UNSIGNED INT")),
			want: true,
		},
		{
			name: "UNSIGNED BIGINT",
			f:    NewFieldType(newMockFieldType("UNSIGNED BIGINT")),
			want: true,
		},
		{
			name: "UNSIGNED SMALLINT",
			f:    NewFieldType(newMockFieldType("UNSIGNED SMALLINT")),
			want: true,
		},
		{
			name: "UNSIGNED TINYINT",
			f:    NewFieldType(newMockFieldType("UNSIGNED TINYINT")),
			want: true,
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
