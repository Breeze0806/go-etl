package postgres

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/lib/pq/oid"
)

func testDecimalColumnValueFromString(s string) element.ColumnValue {
	d, err := element.NewDecimalColumnValueFromString(s)
	if err != nil {
		panic(err)
	}
	return d
}

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
			want: "$22",
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
		//bool
		{
			name: "1",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_bool])),
			want: database.GoTypeBool,
		},

		//int64
		{
			name: "2",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_int2])),
			want: database.GoTypeInt64,
		},
		{
			name: "3",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_int4])),
			want: database.GoTypeInt64,
		},
		{
			name: "4",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_int8])),
			want: database.GoTypeInt64,
		},

		//float64
		{
			name: "5",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_float4])),
			want: database.GoTypeFloat64,
		},
		{
			name: "6",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_float8])),
			want: database.GoTypeFloat64,
		},

		//string
		{
			name: "7",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_varchar])),
			want: database.GoTypeString,
		},
		{
			name: "8",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_text])),
			want: database.GoTypeString,
		},
		{
			name: "9",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_bpchar])),
			want: database.GoTypeString,
		},
		{
			name: "10",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_numeric])),
			want: database.GoTypeString,
		},

		//time
		{
			name: "11",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_date])),
			want: database.GoTypeTime,
		},
		{
			name: "12",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_time])),
			want: database.GoTypeTime,
		},
		{
			name: "13",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_timetz])),
			want: database.GoTypeTime,
		},
		{
			name: "14",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_timestamp])),
			want: database.GoTypeTime,
		},
		{
			name: "15",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_timestamptz])),
			want: database.GoTypeTime,
		},

		//bytes
		{
			name: "16",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_bytea])),
			want: database.GoTypeBytes,
		},
		{
			name: "17",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_uuid])),
			want: database.GoTypeBytes,
		},

		//unknown
		{
			name: "18",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T__bool])),
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

func TestFieldType_IsSupportted(t *testing.T) {
	tests := []struct {
		name string
		f    *FieldType
		want bool
	}{
		{
			name: "1",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T_bool])),
			want: true,
		},
		{
			name: "2",
			f:    NewFieldType(newMockColumnType(oid.TypeName[oid.T__bool])),
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

func TestScanner_Scan(t *testing.T) {
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		s       *Scanner
		args    args
		want    element.Column
		wantErr bool
	}{
		{
			name: "1",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_bool]))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBoolColumnValue(), "f1", 0),
		},
		{
			name: "2",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_bool]))))),
			args: args{
				src: true,
			},
			want: element.NewDefaultColumn(element.NewBoolColumnValue(true), "f1", 0),
		},
		{
			name: "3",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_bool]))))),
			args: args{
				src: 1,
			},
			wantErr: true,
		},

		{
			name: "4",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_int4]))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBigIntColumnValue(), "f1", 0),
		},
		{
			name: "5",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_int8]))))),
			args: args{
				src: int64(123456789012),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(int64(123456789012)), "f1", 0),
		},
		{
			name: "6",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_int2]))))),
			args: args{
				src: "1",
			},
			wantErr: true,
		},

		{
			name: "7",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_uuid]))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBytesColumnValue(), "f1", 0),
		},
		{
			name: "8",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_bytea]))))),
			args: args{
				src: []byte("中国"),
			},
			want: element.NewDefaultColumn(element.NewBytesColumnValue([]byte("中国")), "f1", 0),
		},
		{
			name: "9",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_uuid]))))),
			args: args{
				src: "1",
			},
			wantErr: true,
		},

		{
			name: "10",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_date]))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilTimeColumnValue(), "f1", 0),
		},
		{
			name: "11",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_timestamp]))))),
			args: args{
				src: time.Date(2021, 6, 17, 22, 24, 8, 8, time.UTC),
			},
			want: element.NewDefaultColumn(element.NewTimeColumnValue(
				time.Date(2021, 6, 17, 22, 24, 8, 8, time.UTC)), "f1", 0),
		},
		{
			name: "12",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_timestamptz]))))),
			args: args{
				src: "1",
			},
			wantErr: true,
		},

		{
			name: "13",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_varchar]))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilStringColumnValue(), "f1", 0),
		},
		{
			name: "14",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_text]))))),
			args: args{
				src: "中国",
			},
			want: element.NewDefaultColumn(element.NewStringColumnValue("中国"), "f1", 0),
		},
		{
			name: "15",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_bpchar]))))),
			args: args{
				src: 123,
			},
			wantErr: true,
		},

		{
			name: "16",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_float4]))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilStringColumnValue(), "f1", 0),
		},
		{
			name: "17",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_float8]))))),
			args: args{
				src: 1234567890.1231233,
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1234567890.1231233), "f1", 0),
		},
		{
			name: "18",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_numeric]))))),
			args: args{
				src: []byte("1234567890.1231233"),
			},
			want: element.NewDefaultColumn(testDecimalColumnValueFromString("1234567890.1231233"), "f1", 0),
		},
		{
			name: "19",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_numeric]))))),
			args: args{
				src: "1234567890.1231233",
			},
			wantErr: true,
		},
		{
			name: "20",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_numeric]))))),
			args: args{
				src: []byte("1234567890.1231233a"),
			},
			wantErr: true,
		},

		{
			name: "21",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T__bool]))))),
			args: args{
				src: "1234567890.1231233",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Scan(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("Scanner.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got := tt.s.Column(); !tt.wantErr && !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("Column() = %v, want %v", got, tt.want)
			}
		})
	}
}
