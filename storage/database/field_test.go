package database

import (
	"database/sql"
	"database/sql/driver"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/element"
)

func TestBaseField_Name(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseField
		want string
	}{
		{
			name: "1",
			b:    NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})),
			want: "f1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Name(); got != tt.want {
				t.Errorf("BaseField.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseField_ColumnType(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseField
		want FieldType
	}{
		{
			name: "1",
			b:    NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})),
			want: NewBaseFieldType(&sql.ColumnType{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.FieldType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseField.ColumnType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseField_String(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseField
		want string
	}{
		{
			name: "1",
			b:    NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})),
			want: "f1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("BaseField.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoType_String(t *testing.T) {
	// GoTypeUnknow:  "unknow",
	// GoTypeBool:    "bool",
	// GoTypeInt8:    "int8",
	// GoTypeInt16:   "int16",
	// GoTypeInt32:   "int32",
	// GoTypeInt64:   "int64",
	// GoTypeFloat32: "float32",
	// GoTypeFloat64: "float64",
	// GoTypeString:  "string",
	// GoTypeBytes:   "bytes",
	// GoTypeTime:    "time",
	tests := []struct {
		name string
		t    GoType
		want string
	}{
		{
			name: "1",
			t:    GoTypeUnknown,
			want: "unknow",
		},
		{
			name: "2",
			t:    GoTypeBool,
			want: "bool",
		},
		{
			name: "3",
			t:    GoTypeInt8,
			want: "int8",
		},
		{
			name: "4",
			t:    GoTypeInt16,
			want: "int16",
		},
		{
			name: "5",
			t:    GoTypeInt32,
			want: "int32",
		},
		{
			name: "6",
			t:    GoTypeInt64,
			want: "int64",
		},
		{
			name: "7",
			t:    GoTypeFloat32,
			want: "float32",
		},
		{
			name: "8",
			t:    GoTypeFloat64,
			want: "float64",
		},
		{
			name: "9",
			t:    GoTypeString,
			want: "string",
		},
		{
			name: "10",
			t:    GoTypeBytes,
			want: "bytes",
		},
		{
			name: "11",
			t:    GoTypeTime,
			want: "time",
		},
		{
			name: "12",
			t:    GoType(math.MaxUint8),
			want: "unknow",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("GoType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseScanner_SetColumn(t *testing.T) {
	type args struct {
		c element.Column
	}
	tests := []struct {
		name string
		b    *BaseScanner
		args args
		want element.Column
	}{
		{
			name: "1",
			b:    &BaseScanner{},
			args: args{
				c: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1e16), "test", 0),
			},
			want: element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1e16), "test", 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetColumn(tt.args.c)
			if got := tt.b.Column(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseScanner.Column() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoValuer_Value(t *testing.T) {
	tests := []struct {
		name    string
		g       *GoValuer
		want    driver.Value
		wantErr bool
	}{
		{
			name: "1",
			g: NewGoValuer(newMockField(NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})), newMockFieldType(GoTypeBool)),
				element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1234567890.23), "test", 0)),
			want: true,
		},
		{
			name: "2",
			g: NewGoValuer(newMockField(NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})), newMockFieldType(GoTypeInt8)),
				element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1234567890.23), "test", 0)),
			wantErr: true,
			want:    int8(0),
		},

		{
			name: "3",
			g: NewGoValuer(newMockField(NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})), newMockFieldType(GoTypeInt32)),
				element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1234567890.23), "test", 0)),
			want: int32(1234567890),
		},
		{
			name: "4",
			g: NewGoValuer(newMockField(NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})), newMockFieldType(GoTypeInt64)),
				element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1234567890.23), "test", 0)),
			want: int64(1234567890),
		},
		{
			name: "5",
			g: NewGoValuer(newMockField(NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})), newMockFieldType(GoTypeFloat32)),
				element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1234567890.23), "test", 0)),
			want: float32(1234567890.23),
		},
		{
			name: "6",
			g: NewGoValuer(newMockField(NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})), newMockFieldType(GoTypeFloat64)),
				element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1234567890.23), "test", 0)),
			want: float64(1234567890.23),
		},
		{
			name: "7",
			g: NewGoValuer(newMockField(NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})), newMockFieldType(GoTypeString)),
				element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1234567890.23), "test", 0)),
			want: "1234567890.23",
		},
		{
			name: "8",
			g: NewGoValuer(newMockField(NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})), newMockFieldType(GoTypeBytes)),
				element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1234567890.23), "test", 0)),
			want: []byte("1234567890.23"),
		},
		{
			name: "9",
			g: NewGoValuer(newMockField(NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})), newMockFieldType(GoTypeTime)),
				element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1234567890.23), "test", 0)),
			wantErr: true,
			want:    time.Time{},
		},
		{
			name: "10",
			g: NewGoValuer(newMockField(NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})), newMockFieldType(GoTypeUnknown)),
				element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1234567890.23), "test", 0)),
			wantErr: true,
		},
		{
			name: "11",
			g: NewGoValuer(newMockField(NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})), NewBaseFieldType(&sql.ColumnType{})),
				element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1234567890.23), "test", 0)),
			wantErr: true,
		},
		{
			name: "12",
			g: NewGoValuer(newMockField(NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})), newMockFieldType(GoTypeInt64)),
				element.NewDefaultColumn(element.NewNilBigIntColumnValue(), "test", 0)),
			wantErr: false,
		},
		{
			name: "13",
			g: NewGoValuer(newMockField(NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})), newMockFieldType(GoTypeInt16)),
				element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1234567890.23), "test", 0)),
			wantErr: true,
			want:    int16(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("GoValuer.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoValuer.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseField_Index(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseField
		want int
	}{
		{
			name: "1",
			b:    NewBaseField(1, "f1", NewBaseFieldType(&sql.ColumnType{})),
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Index(); got != tt.want {
				t.Errorf("BaseField.Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseFieldType_IsSupportted(t *testing.T) {
	type fields struct {
		ColumnType ColumnType
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "1",
			fields: fields{
				ColumnType: &sql.ColumnType{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BaseFieldType{
				ColumnType: tt.fields.ColumnType,
			}
			if got := b.IsSupportted(); got != tt.want {
				t.Errorf("BaseFieldType.IsSupportted() = %v, want %v", got, tt.want)
			}
		})
	}
}
