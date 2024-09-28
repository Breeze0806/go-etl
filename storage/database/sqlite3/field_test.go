package sqlite3

import (
	"database/sql"
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"reflect"
	"testing"
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
			want: "`f1`",
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
		//INTEGER
		{
			name: "1",
			f:    NewFieldType(newMockColumnType("INTEGER")),
			want: database.GoTypeString,
		},

		//BLOB
		{
			name: "2",
			f:    NewFieldType(newMockColumnType("BLOB")),
			want: database.GoTypeString,
		},
		//NUMERIC
		{
			name: "3",
			f:    NewFieldType(newMockColumnType("NUMERIC")),
			want: database.GoTypeString,
		},
		//REAL
		{
			name: "4",
			f:    NewFieldType(newMockColumnType("REAL")),
			want: database.GoTypeString,
		},
		//TEXT
		{
			name: "5",
			f:    NewFieldType(newMockColumnType("TEXT")),
			want: database.GoTypeString,
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
		src interface{}
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
			name: "1",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("INTEGER"))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBigIntColumnValue(), "f1", 0),
		},
		{
			name: "2",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("INTEGER"))))),
			args: args{
				src: int64(2 ^ 63 - 1),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(int64(2^63-1)), "f1", element.ByteSize(int64(2^63-1))),
		},
		{
			name: "3",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("BLOB"))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilBytesColumnValue(), "f1", 0),
		},
		{
			name: "4",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("BLOB"))))),
			args: args{
				src: []byte("123"),
			},
			want: element.NewDefaultColumn(element.NewBytesColumnValue([]byte("123")), "f1", element.ByteSize([]byte("123"))),
		},
		{
			name: "5",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("NUMERIC"))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilTimeColumnValue(), "f1", 0),
		},
		{
			name: "6",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("NUMERIC"))))),
			args: args{
				src: int64(2 ^ 63 - 1),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(int64(2^63-1)), "f1", element.ByteSize(int64(2^63-1))),
		},
		{
			name: "7",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("NUMERIC"))))),
			args: args{
				src: 1.23456789,
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1.23456789), "f1", element.ByteSize(1.23456789)),
		},
		{
			name: "8",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("REAL"))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilTimeColumnValue(), "f1", 0),
		},
		{
			name: "9",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("REAL"))))),
			args: args{
				src: int64(2 ^ 63 - 1),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(int64(2^63-1)), "f1", element.ByteSize(int64(2^63-1))),
		},
		{
			name: "10",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("REAL"))))),
			args: args{
				src: 1.23456789,
			},
			want: element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1.23456789), "f1", element.ByteSize(1.23456789)),
		},
		{
			name: "11",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("TEXT"))))),
			args: args{
				src: nil,
			},
			want: element.NewDefaultColumn(element.NewNilStringColumnValue(), "f1", 0),
		},
		{
			name: "12",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("TEXT"))))),
			args: args{
				src: "123",
			},
			want: element.NewDefaultColumn(element.NewStringColumnValue("123"), "f1", element.ByteSize("123")),
		},
		{
			name: "13",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("INTEGER"))))),
			args: args{
				src: "123",
			},
			wantErr: true,
		},
		{
			name: "14",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("BLOB"))))),
			args: args{
				src: 123,
			},
			wantErr: true,
		},
		{
			name: "15",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("NUMERIC"))))),
			args: args{
				src: 123,
			},
			wantErr: true,
		},
		{
			name: "16",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("REAL"))))),
			args: args{
				src: 123,
			},
			wantErr: true,
		},
		{
			name: "17",
			s: NewScanner(NewField(database.NewBaseField(0,
				"f1", NewFieldType(newMockColumnType("TEXT"))))),
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

			if !reflect.DeepEqual(tt.s.Column(), tt.want) {
				t.Errorf("Column() = %v %v, want %v", tt.s.Column().ByteSize(), tt.s.Column(), tt.want)
			}
		})
	}
}
