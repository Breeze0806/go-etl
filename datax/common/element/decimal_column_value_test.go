package element

import (
	"math"
	"math/big"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNewNilDecimalColumnValue(t *testing.T) {
	tests := []struct {
		name string
		want ColumnValue
	}{
		{
			name: "1",
			want: NewNilDecimalColumnValue(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNilDecimalColumnValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNilDecimalColumnValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilDecimalColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		n    *NilDecimalColumnValue
		want ColumnType
	}{
		{
			name: "1",
			n:    NewNilDecimalColumnValue().(*NilDecimalColumnValue),
			want: TypeDecimal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilDecimalColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilDecimalColumnValue_clone(t *testing.T) {
	tests := []struct {
		name string
		n    *NilDecimalColumnValue
		want ColumnValue
	}{
		{
			name: "1",
			n:    NewNilDecimalColumnValue().(*NilDecimalColumnValue),
			want: NewNilDecimalColumnValue(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.clone()
			if got == tt.n {
				t.Errorf("NilDecimalColumnValue.clone() = %p, n %p", got, tt.n)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilDecimalColumnValue.clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDecimalColumnValueFromFloat(t *testing.T) {
	type args struct {
		f float64
	}
	tests := []struct {
		name string
		args args
		want ColumnValue
	}{
		{
			name: "zero",
			args: args{
				f: 0.0,
			},
			want: NewDecimalColumnValueFromFloat(0.0),
		},

		{
			name: "maxfloat64",
			args: args{
				f: math.MaxFloat64,
			},
			want: NewDecimalColumnValueFromFloat(math.MaxFloat64),
		},

		{
			name: "1",
			args: args{
				f: 0.00000012345,
			},
			want: NewDecimalColumnValue(decimal.New(12345, -11)),
		},

		{
			name: "2",
			args: args{
				f: -0.00000012345,
			},
			want: NewDecimalColumnValue(decimal.New(-12345, -11)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDecimalColumnValueFromFloat(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDecimalColumnValueFromFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDecimalColumnValue(t *testing.T) {
	type args struct {
		d decimal.Decimal
	}
	tests := []struct {
		name string
		args args
		want ColumnValue
	}{
		{
			name: "1",
			args: args{
				d: decimal.New(12345, -11),
			},
			want: NewDecimalColumnValueFromFloat(0.00000012345),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDecimalColumnValue(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDecimalColumnValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDecimalColumnValueFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    ColumnValue
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				s: "-0.00000012345",
			},
			want: NewDecimalColumnValueFromFloat(-0.00000012345),
		},
		{
			name: "1",
			args: args{
				s: "0.00000012345",
			},
			want: NewDecimalColumnValueFromFloat(0.00000012345),
		},
		{
			name: "MaxFloat64",
			args: args{
				s: strconv.FormatFloat(math.MaxFloat64, 'f', -1, 64),
			},
			want: NewDecimalColumnValueFromFloat(math.MaxFloat64),
		},

		{
			name: "MaxFloat32",
			args: args{
				s: strconv.FormatFloat(math.MaxFloat32, 'f', -1, 64),
			},
			want: NewDecimalColumnValueFromFloat(math.MaxFloat32),
		},
		{
			name: "NegMaxFloat32",
			args: args{
				s: "-" + strconv.FormatFloat(math.MaxFloat32, 'f', -1, 64),
			},
			want: NewDecimalColumnValueFromFloat(-math.MaxFloat32),
		},
		{
			name: "2",
			args: args{
				s: "-1232000000000000",
			},
			want: NewDecimalColumnValue(decimal.New(-1232, 12)),
		},
		{
			name: "2",
			args: args{
				s: "1232000000000000",
			},
			want: NewDecimalColumnValue(decimal.New(1232, 12)),
		},
		{
			name: "2.23e10",
			args: args{
				s: "2.23e10",
			},
			want: NewDecimalColumnValue(decimal.New(223, 8)),
		},
		{
			name: "2.23e10",
			args: args{
				s: "2.23e-10",
			},
			want: NewDecimalColumnValue(decimal.New(223, -12)),
		},
		{
			name: "abc",
			args: args{
				s: "abc",
			},
			wantErr: true,
		},
		{
			name: "abc31232131",
			args: args{
				s: "abc31232131",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDecimalColumnValueFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDecimalColumnValueFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("NewDecimalColumnValueFromString() = %v, want %v args: %#v", got, tt.want, tt.args)
			}
		})
	}
}

func TestDecimalColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		d    *DecimalColumnValue
		want ColumnType
	}{
		{
			name: "1",
			d:    NewDecimalColumnValue(decimal.Zero).(*DecimalColumnValue),
			want: TypeDecimal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecimalColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalColumnValue_AsBool(t *testing.T) {
	tests := []struct {
		name    string
		d       *DecimalColumnValue
		want    bool
		wantErr bool
	}{
		{
			name: "zero1",
			d:    NewDecimalColumnValue(decimal.Zero).(*DecimalColumnValue),
			want: false,
		},
		{
			name: "zero2",
			d:    testDecimalColumnValueFormString("0").(*DecimalColumnValue),
			want: false,
		},
		{
			name: "1",
			d:    testDecimalColumnValueFormString("-0.0000000000000000000001").(*DecimalColumnValue),
			want: true,
		},
		{
			name: "2",
			d:    NewDecimalColumnValueFromFloat(math.MaxFloat64).(*DecimalColumnValue),
			want: true,
		},
		{
			name: "3",
			d:    NewDecimalColumnValueFromFloat(-math.MaxFloat32).(*DecimalColumnValue),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.AsBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("DecimalColumnValue.AsBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecimalColumnValue.AsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalColumnValue_AsBigInt(t *testing.T) {
	tests := []struct {
		name    string
		d       *DecimalColumnValue
		want    *big.Int
		wantErr bool
	}{
		{
			name: "zero",
			d:    NewDecimalColumnValue(decimal.Zero).(*DecimalColumnValue),
			want: big.NewInt(0),
		},

		{
			name: "1",
			d:    testDecimalColumnValueFormString("123450000").(*DecimalColumnValue),
			want: big.NewInt(123450000),
		},

		{
			name: "2",
			d:    testDecimalColumnValueFormString("1.00122323123").(*DecimalColumnValue),
			want: big.NewInt(1),
		},
		{
			name: "3",
			d:    testDecimalColumnValueFormString("123456232542542525.525254252524").(*DecimalColumnValue),
			want: big.NewInt(123456232542542525),
		},
		{
			name: "4",
			d:    testDecimalColumnValueFormString("0.00122323123").(*DecimalColumnValue),
			want: big.NewInt(0),
		},
		{
			name: "5",
			d:    testDecimalColumnValueFormString("-123456232542542525.525254252524").(*DecimalColumnValue),
			want: big.NewInt(-123456232542542525),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.AsBigInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("DecimalColumnValue.AsBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Cmp(tt.want) != 0 {
				t.Errorf("DecimalColumnValue.AsBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalColumnValue_AsDecimal(t *testing.T) {
	tests := []struct {
		name    string
		d       *DecimalColumnValue
		want    decimal.Decimal
		wantErr bool
	}{
		{
			name: "zero",
			d:    NewDecimalColumnValue(decimal.Zero).(*DecimalColumnValue),
			want: decimal.Zero,
		},

		{
			name: "1",
			d:    testDecimalColumnValueFormString("123450000").(*DecimalColumnValue),
			want: testDecimalFormString("123450000"),
		},

		{
			name: "2",
			d:    testDecimalColumnValueFormString("1.00122323123").(*DecimalColumnValue),
			want: testDecimalFormString("1.00122323123"),
		},
		{
			name: "3",
			d:    testDecimalColumnValueFormString("123456232542542525.525254252524").(*DecimalColumnValue),
			want: testDecimalFormString("123456232542542525.525254252524"),
		},
		{
			name: "4",
			d:    testDecimalColumnValueFormString("0.00122323123").(*DecimalColumnValue),
			want: testDecimalFormString("0.00122323123"),
		},
		{
			name: "5",
			d:    testDecimalColumnValueFormString("-123456232542542525.525254252524").(*DecimalColumnValue),
			want: testDecimalFormString("-123456232542542525.525254252524"),
		},
		{
			name: "6",
			d:    testDecimalColumnValueFormString("-123456000000").(*DecimalColumnValue),
			want: testDecimalFormString("-123456000000"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.AsDecimal()
			if (err != nil) != tt.wantErr {
				t.Errorf("DecimalColumnValue.AsDecimal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("DecimalColumnValue.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalColumnValue_AsString(t *testing.T) {
	tests := []struct {
		name    string
		d       *DecimalColumnValue
		want    string
		wantErr bool
	}{
		{
			name: "zero",
			d:    NewDecimalColumnValue(decimal.Zero).(*DecimalColumnValue),
			want: "0",
		},

		{
			name: "1",
			d:    testDecimalColumnValueFormString("123450000").(*DecimalColumnValue),
			want: "123450000",
		},

		{
			name: "2",
			d:    testDecimalColumnValueFormString("1.00122323123").(*DecimalColumnValue),
			want: "1.00122323123",
		},
		{
			name: "3",
			d:    testDecimalColumnValueFormString("123456232542542525.525254252524").(*DecimalColumnValue),
			want: "123456232542542525.525254252524",
		},
		{
			name: "4",
			d:    testDecimalColumnValueFormString("0.00122323123").(*DecimalColumnValue),
			want: "0.00122323123",
		},
		{
			name: "5",
			d:    testDecimalColumnValueFormString("-123456232542542525.525254252524").(*DecimalColumnValue),
			want: "-123456232542542525.525254252524",
		},
		{
			name: "6",
			d:    testDecimalColumnValueFormString("-123456000000").(*DecimalColumnValue),
			want: "-123456000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("DecimalColumnValue.AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecimalColumnValue.AsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalColumnValue_AsBytes(t *testing.T) {
	tests := []struct {
		name    string
		d       *DecimalColumnValue
		want    []byte
		wantErr bool
	}{
		{
			name: "zero",
			d:    NewDecimalColumnValue(decimal.Zero).(*DecimalColumnValue),
			want: []byte("0"),
		},

		{
			name: "1",
			d:    testDecimalColumnValueFormString("123450000").(*DecimalColumnValue),
			want: []byte("123450000"),
		},

		{
			name: "2",
			d:    testDecimalColumnValueFormString("1.00122323123").(*DecimalColumnValue),
			want: []byte("1.00122323123"),
		},
		{
			name: "3",
			d:    testDecimalColumnValueFormString("123456232542542525.525254252524").(*DecimalColumnValue),
			want: []byte("123456232542542525.525254252524"),
		},
		{
			name: "4",
			d:    testDecimalColumnValueFormString("0.00122323123").(*DecimalColumnValue),
			want: []byte("0.00122323123"),
		},
		{
			name: "5",
			d:    testDecimalColumnValueFormString("-123456232542542525.525254252524").(*DecimalColumnValue),
			want: []byte("-123456232542542525.525254252524"),
		},
		{
			name: "6",
			d:    testDecimalColumnValueFormString("-123456000000").(*DecimalColumnValue),
			want: []byte("-123456000000"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.AsBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("DecimalColumnValue.AsBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecimalColumnValue.AsBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalColumnValue_AsTime(t *testing.T) {
	tests := []struct {
		name    string
		d       *DecimalColumnValue
		want    time.Time
		wantErr bool
	}{
		{
			name:    "zero",
			d:       NewDecimalColumnValue(decimal.Zero).(*DecimalColumnValue),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.AsTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("DecimalColumnValue.AsTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecimalColumnValue.AsTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalColumnValue_String(t *testing.T) {
	tests := []struct {
		name string
		d    *DecimalColumnValue
		want string
	}{
		{
			name: "zero",
			d:    NewDecimalColumnValue(decimal.Zero).(*DecimalColumnValue),
			want: "0",
		},

		{
			name: "1",
			d:    testDecimalColumnValueFormString("123450000").(*DecimalColumnValue),
			want: "123450000",
		},

		{
			name: "2",
			d:    testDecimalColumnValueFormString("1.00122323123").(*DecimalColumnValue),
			want: "1.00122323123",
		},
		{
			name: "3",
			d:    testDecimalColumnValueFormString("123456232542542525.525254252524").(*DecimalColumnValue),
			want: "123456232542542525.525254252524",
		},
		{
			name: "4",
			d:    testDecimalColumnValueFormString("0.00122323123").(*DecimalColumnValue),
			want: "0.00122323123",
		},
		{
			name: "5",
			d:    testDecimalColumnValueFormString("-123456232542542525.525254252524").(*DecimalColumnValue),
			want: "-123456232542542525.525254252524",
		},
		{
			name: "6",
			d:    testDecimalColumnValueFormString("-123456000000").(*DecimalColumnValue),
			want: "-123456000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("DecimalColumnValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalColumnValue_clone(t *testing.T) {
	tests := []struct {
		name string
		d    *DecimalColumnValue
		want ColumnValue
	}{
		{
			name: "zero",
			d:    NewDecimalColumnValue(decimal.Zero).(*DecimalColumnValue),
			want: NewDecimalColumnValue(decimal.Zero),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.clone()
			if got == tt.d {
				t.Errorf("DecimalColumnValue.clone() = %p, d %p", got, tt.d)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecimalColumnValue.clone() = %v, want %v", got, tt.want)
			}
		})
	}
}
