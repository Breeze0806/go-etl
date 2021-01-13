package element

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNilBoolColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		n    *NilBoolColumnValue
		want ColumnType
	}{
		{
			name: "1",
			n:    NewNilBoolColumnValue().(*NilBoolColumnValue),
			want: TypeBool,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilBoolColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilBoolColumnValue_Clone(t *testing.T) {
	tests := []struct {
		name string
		n    *NilBoolColumnValue
		want ColumnValue
	}{
		{
			name: "1",
			n:    NewNilBoolColumnValue().(*NilBoolColumnValue),
			want: NewNilBoolColumnValue(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.Clone()
			if got == tt.n {
				t.Errorf("NilBigIntColumnValue.Clone() = %p, n %p want %p", got, tt.n, tt.want)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilBoolColumnValue.clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		b    *BoolColumnValue
		want ColumnType
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: TypeBool,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoolColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_AsBool(t *testing.T) {
	tests := []struct {
		name    string
		b       *BoolColumnValue
		want    bool
		wantErr bool
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: true,
		},
		{
			name: "false",
			b:    NewBoolColumnValue(false).(*BoolColumnValue),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolColumnValue.AsBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BoolColumnValue.AsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_AsBigInt(t *testing.T) {
	tests := []struct {
		name    string
		b       *BoolColumnValue
		want    *big.Int
		wantErr bool
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: big.NewInt(1),
		},
		{
			name: "false",
			b:    NewBoolColumnValue(false).(*BoolColumnValue),
			want: big.NewInt(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsBigInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolColumnValue.AsBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoolColumnValue.AsBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_AsDecimal(t *testing.T) {
	tests := []struct {
		name    string
		b       *BoolColumnValue
		want    decimal.Decimal
		wantErr bool
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: decimal.New(1, 0),
		},
		{
			name: "false",
			b:    NewBoolColumnValue(false).(*BoolColumnValue),
			want: decimal.Zero,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsDecimal()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolColumnValue.AsDecimal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoolColumnValue.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_AsString(t *testing.T) {
	tests := []struct {
		name    string
		b       *BoolColumnValue
		want    string
		wantErr bool
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: "true",
		},
		{
			name: "false",
			b:    NewBoolColumnValue(false).(*BoolColumnValue),
			want: "false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolColumnValue.AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BoolColumnValue.AsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_AsBytes(t *testing.T) {
	tests := []struct {
		name    string
		b       *BoolColumnValue
		want    string
		wantErr bool
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: "true",
		},
		{
			name: "false",
			b:    NewBoolColumnValue(false).(*BoolColumnValue),
			want: "false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolColumnValue.AsBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, []byte(tt.want)) {
				t.Errorf("BoolColumnValue.AsBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_AsTime(t *testing.T) {
	tests := []struct {
		name    string
		b       *BoolColumnValue
		want    time.Time
		wantErr bool
	}{
		{
			name:    "true",
			b:       NewBoolColumnValue(true).(*BoolColumnValue),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.AsTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolColumnValue.AsTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoolColumnValue.AsTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_String(t *testing.T) {
	tests := []struct {
		name string
		b    *BoolColumnValue
		want string
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: "true",
		},
		{
			name: "false",
			b:    NewBoolColumnValue(false).(*BoolColumnValue),
			want: "false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("BoolColumnValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolColumnValue_Clone(t *testing.T) {
	tests := []struct {
		name string
		b    *BoolColumnValue
		want ColumnValue
	}{
		{
			name: "true",
			b:    NewBoolColumnValue(true).(*BoolColumnValue),
			want: NewBoolColumnValue(true),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.Clone()
			if got == tt.b {
				t.Errorf("BoolColumnValue.Clone() = %p, b %p", got, tt.b)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoolColumnValue.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}
