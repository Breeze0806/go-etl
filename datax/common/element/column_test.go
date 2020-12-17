package element

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func Test_notNilColumnValue_IsNil(t *testing.T) {
	tests := []struct {
		name string
		n    *notNilColumnValue
		want bool
	}{
		{
			name: "1",
			n:    &notNilColumnValue{},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.IsNil(); got != tt.want {
				t.Errorf("notNilColumnValue.IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		n    *nilColumnValue
		want ColumnType
	}{
		{
			name: "1",
			n:    &nilColumnValue{},
			want: TypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nilColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_IsNil(t *testing.T) {
	tests := []struct {
		name string
		n    *nilColumnValue
		want bool
	}{
		{
			name: "1",
			n:    &nilColumnValue{},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.IsNil(); got != tt.want {
				t.Errorf("nilColumnValue.IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_AsBool(t *testing.T) {
	tests := []struct {
		name    string
		n       *nilColumnValue
		want    bool
		wantErr bool
	}{
		{
			name:    "1",
			n:       &nilColumnValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.AsBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("nilColumnValue.AsBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("nilColumnValue.AsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_AsBigInt(t *testing.T) {
	tests := []struct {
		name    string
		n       *nilColumnValue
		want    *big.Int
		wantErr bool
	}{
		{
			name:    "1",
			n:       &nilColumnValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.AsBigInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("nilColumnValue.AsBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nilColumnValue.AsBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_AsDecimal(t *testing.T) {
	tests := []struct {
		name    string
		n       *nilColumnValue
		want    decimal.Decimal
		wantErr bool
	}{
		{
			name:    "1",
			n:       &nilColumnValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.AsDecimal()
			if (err != nil) != tt.wantErr {
				t.Errorf("nilColumnValue.AsDecimal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nilColumnValue.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_AsString(t *testing.T) {
	tests := []struct {
		name    string
		n       *nilColumnValue
		want    string
		wantErr bool
	}{
		{
			name:    "1",
			n:       &nilColumnValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("nilColumnValue.AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("nilColumnValue.AsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_AsBytes(t *testing.T) {
	tests := []struct {
		name    string
		n       *nilColumnValue
		want    []byte
		wantErr bool
	}{
		{
			name:    "1",
			n:       &nilColumnValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.AsBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("nilColumnValue.AsBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nilColumnValue.AsBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_AsTime(t *testing.T) {
	tests := []struct {
		name    string
		n       *nilColumnValue
		want    time.Time
		wantErr bool
	}{
		{
			name:    "1",
			n:       &nilColumnValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.AsTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("nilColumnValue.AsTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nilColumnValue.AsTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nilColumnValue_String(t *testing.T) {
	tests := []struct {
		name string
		n    *nilColumnValue
		want string
	}{
		{
			name: "1",
			n:    &nilColumnValue{},
			want: "<nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.String(); got != tt.want {
				t.Errorf("nilColumnValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultColumn_Name(t *testing.T) {
	tests := []struct {
		name string
		d    *DefaultColumn
		want string
	}{
		{
			name: "1",
			d:    NewDefaultColumn(NewNilBigIntColumnValue(), "test", 12).(*DefaultColumn),
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Name(); got != tt.want {
				t.Errorf("DefaultColumn.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultColumn_Clone(t *testing.T) {
	tests := []struct {
		name string
		d    *DefaultColumn
		want Column
	}{
		{
			name: "1",
			d:    NewDefaultColumn(NewNilBigIntColumnValue(), "test", 12).(*DefaultColumn),
			want: NewDefaultColumn(NewNilBigIntColumnValue(), "test", 12),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultColumn.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultColumn_ByteSize(t *testing.T) {
	tests := []struct {
		name string
		d    *DefaultColumn
		want int64
	}{
		{
			name: "1",
			d:    NewDefaultColumn(NewNilBigIntColumnValue(), "test", 12).(*DefaultColumn),
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.ByteSize(); got != tt.want {
				t.Errorf("DefaultColumn.ByteSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultColumn_MemorySize(t *testing.T) {
	tests := []struct {
		name string
		d    *DefaultColumn
		want int64
	}{
		{
			name: "1",
			d:    NewDefaultColumn(NewNilBigIntColumnValue(), "test", 12).(*DefaultColumn),
			want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.MemorySize(); got != tt.want {
				t.Errorf("DefaultColumn.MemorySize() = %v, want %v", got, tt.want)
			}
		})
	}
}
