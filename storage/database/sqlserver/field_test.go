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

package sqlserver

import (
	"database/sql/driver"
	"reflect"
	"testing"
	"time"

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

func (m *mockFieldType) IsSupportted() bool {
	return true
}

func TestField_Quoted(t *testing.T) {
	tests := []struct {
		name string
		f    *Field
		want string
	}{
		{
			name: "1",
			f:    NewField(database.NewBaseField(1, "f1", newMockFieldType("DATE"))),
			want: `[f1]`,
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
			f:    NewField(database.NewBaseField(1, "f1", newMockFieldType("DATE"))),
			args: args{
				i: 12345,
			},
			want: `@p12345`,
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
			f:    NewField(database.NewBaseField(1, "f1", newMockFieldType("DATE"))),
			want: `[f1]`,
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
			f:    NewField(database.NewBaseField(1, "f1", newMockFieldType("DATE"))),
			want: NewFieldType(newMockFieldType("DATE")),
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
			f:    NewField(database.NewBaseField(0, "f1", newMockFieldType("DATE"))),
			want: NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("DATE")))),
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
			f:    NewField(database.NewBaseField(0, "f1", newMockFieldType("DATE"))),
			args: args{
				c: element.NewDefaultColumn(nil, "", 0),
			},
			want: NewValuer(NewField(database.NewBaseField(0, "f1", newMockFieldType("DATE"))), element.NewDefaultColumn(nil, "", 0)),
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
			f:    NewFieldType(newMockFieldType("DATE")),
			want: true,
		},
		{
			name: "2",
			f:    NewFieldType(newMockFieldType("DATE1")),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.IsSupportted(); got != tt.want {
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
		{
			name: "BIT",
			f:    NewFieldType(newMockFieldType("BIT")),
			want: database.GoTypeBool,
		},

		{
			name: "TINYINT",
			f:    NewFieldType(newMockFieldType("TINYINT")),
			want: database.GoTypeInt64,
		},
		{
			name: "SMALLINT",
			f:    NewFieldType(newMockFieldType("SMALLINT")),
			want: database.GoTypeInt64,
		},
		{
			name: "INT",
			f:    NewFieldType(newMockFieldType("INT")),
			want: database.GoTypeInt64,
		},
		{
			name: "BIGINT",
			f:    NewFieldType(newMockFieldType("BIGINT")),
			want: database.GoTypeInt64,
		},

		{
			name: "REAL",
			f:    NewFieldType(newMockFieldType("REAL")),
			want: database.GoTypeFloat64,
		},
		{
			name: "FLOAT",
			f:    NewFieldType(newMockFieldType("FLOAT")),
			want: database.GoTypeFloat64,
		},

		{
			name: "DECIMAL",
			f:    NewFieldType(newMockFieldType("DECIMAL")),
			want: database.GoTypeString,
		},
		{
			name: "VARCHAR",
			f:    NewFieldType(newMockFieldType("VARCHAR")),
			want: database.GoTypeString,
		},
		{
			name: "NVARCHAR",
			f:    NewFieldType(newMockFieldType("NVARCHAR")),
			want: database.GoTypeString,
		},
		{
			name: "CHAR",
			f:    NewFieldType(newMockFieldType("CHAR")),
			want: database.GoTypeString,
		},
		{
			name: "NCHAR",
			f:    NewFieldType(newMockFieldType("NCHAR")),
			want: database.GoTypeString,
		},
		{
			name: "TEXT",
			f:    NewFieldType(newMockFieldType("TEXT")),
			want: database.GoTypeString,
		},
		{
			name: "NTEXT",
			f:    NewFieldType(newMockFieldType("NTEXT")),
			want: database.GoTypeString,
		},

		{
			name: "SMALLDATETIME",
			f:    NewFieldType(newMockFieldType("SMALLDATETIME")),
			want: database.GoTypeTime,
		},
		{
			name: "DATETIME",
			f:    NewFieldType(newMockFieldType("DATETIME")),
			want: database.GoTypeTime,
		},
		{
			name: "DATETIME2",
			f:    NewFieldType(newMockFieldType("DATETIME2")),
			want: database.GoTypeTime,
		},
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
			name: "DATETIMEOFFSET",
			f:    NewFieldType(newMockFieldType("DATETIMEOFFSET")),
			want: database.GoTypeTime,
		},

		{
			name: "VARBINARY",
			f:    NewFieldType(newMockFieldType("VARBINARY")),
			want: database.GoTypeBytes,
		},
		{
			name: "BINARY",
			f:    NewFieldType(newMockFieldType("BINARY")),
			want: database.GoTypeBytes,
		},

		{
			name: "VARBINARY1",
			f:    NewFieldType(newMockFieldType("VARBINARY1")),
			want: database.GoTypeUnknown,
		},
		{
			name: "BINARY1",
			f:    NewFieldType(newMockFieldType("BINARY1")),
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
		args    args
		wantErr bool
		want    element.Column
	}{
		{
			name: "BIT",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("BIT")))),
			args: args{
				true,
			},
			want: element.NewDefaultColumn(element.NewBoolColumnValue(true), "f1", 0),
		},
		{
			name: "BITNull",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("BIT")))),
			args: args{
				nil,
			},
			want: element.NewDefaultColumn(element.NewNilBoolColumnValue(), "f1", 0),
		},
		{
			name: "BITErr",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("BIT")))),
			args: args{
				"true",
			},
			wantErr: true,
		},

		{
			name: "INT",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("INT")))),
			args: args{
				int64(123456789),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(123456789), "f1", 0),
		},
		{
			name: "BIGINTNull",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("BIGINT")))),
			args: args{
				nil,
			},
			want: element.NewDefaultColumn(element.NewNilBigIntColumnValue(), "f1", 0),
		},
		{
			name: "SMALLINTErr",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("SMALLINT")))),
			args: args{
				123456789,
			},
			wantErr: true,
		},

		{
			name: "REAL",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("REAL")))),
			args: args{
				float32(123456789.1),
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(float64(123456789.1)), "f1", 0),
		},
		{
			name: "FLOAT",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("FLOAT")))),
			args: args{
				float64(123456789.1234),
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(float64(123456789.1234)), "f1", 0),
		},
		{
			name: "DECIMAL",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("DECIMAL")))),
			args: args{
				[]byte("123456789.0123456789"),
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(float64(123456789.0123456789)), "f1", 0),
		},
		{
			name: "DECIMALErr",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("DECIMAL")))),
			args: args{
				[]byte("x123456789.0123456789"),
			},
			wantErr: true,
		},
		{
			name: "DECIMALNull",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("DECIMAL")))),
			args: args{
				nil,
			},
			want: element.NewDefaultColumn(element.NewNilDecimalColumnValue(), "f1", 0),
		},
		{
			name: "FLOATErr",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("DECIMAL")))),
			args: args{
				123,
			},
			wantErr: true,
		},

		{
			name: "VARCHAR",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("VARCHAR")))),
			args: args{
				"中文1234abc",
			},
			want: element.NewDefaultColumn(element.NewStringColumnValue("中文1234abc"), "f1", 0),
		},
		{
			name: "NVARCHARNull",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("NVARCHAR")))),
			args: args{
				nil,
			},
			want: element.NewDefaultColumn(element.NewNilStringColumnValue(), "f1", 0),
		},
		{
			name: "TEXTErr",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("TEXT")))),
			args: args{
				[]byte("123"),
			},
			wantErr: true,
		},

		{
			name: "VARBINARY",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("VARBINARY")))),
			args: args{
				[]byte("中文1234abc"),
			},
			want: element.NewDefaultColumn(element.NewBytesColumnValueNoCopy([]byte("中文1234abc")), "f1", 0),
		},
		{
			name: "BINARYNull",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("BINARY")))),
			args: args{
				nil,
			},
			want: element.NewDefaultColumn(element.NewNilBytesColumnValue(), "f1", 0),
		},
		{
			name: "BINARYErr",
			s:    NewScanner(NewField(database.NewBaseField(0, "f1", newMockFieldType("BINARY")))),
			args: args{
				"123",
			},
			wantErr: true,
		},

		{
			name: "DATE",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DATE")))),
			args: args{
				src: time.Date(2022, 9, 4, 14, 56, 0, 0, time.Local),
			},
			want: element.NewDefaultColumn(
				element.NewTimeColumnValueWithDecoder(time.Date(2022, 9, 4, 14, 56, 0, 0, time.Local), element.NewStringTimeDecoder(dateLayout)), "test", 0),
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
			name: "SMALLDATETIME",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("SMALLDATETIME")))),
			args: args{
				src: time.Date(2022, 9, 4, 14, 56, 0, 0, time.Local),
			},
			want: element.NewDefaultColumn(
				element.NewTimeColumnValueWithDecoder(time.Date(2022, 9, 4, 14, 56, 0, 0, time.Local), element.NewStringTimeDecoder(datetimeLayout)), "test", 0),
		},
		{
			name: "DATETIMENull",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DATETIME")))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilTimeColumnValue(), "test", 0),
		},
		{
			name: "DATETIME2err",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DATETIME2")))),
			args: args{
				src: "123",
			},
			wantErr: true,
		},

		{
			name: "err",
			s:    NewScanner(NewField(database.NewBaseField(0, "test", newMockFieldType("DATETIME1")))),
			args: args{
				src: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Scan(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("Scanner.Scan() error = %v, wantErr %v", err, tt.wantErr)
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
			v: NewValuer(NewField(database.NewBaseField(0, "f1", newMockFieldType("VARBINARY"))),
				element.NewDefaultColumn(element.NewNilBoolColumnValue(), "", 0)),
			want: []byte(nil),
		},
		{
			name: "2",
			v: NewValuer(NewField(database.NewBaseField(0, "f1", newMockFieldType("VARCHAR"))),
				element.NewDefaultColumn(element.NewNilBoolColumnValue(), "", 0)),
			want: nil,
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
