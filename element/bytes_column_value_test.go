package element

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNilBytesColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		n    *NilBytesColumnValue
		want ColumnType
	}{
		{
			name: "1",
			n:    NewNilBytesColumnValue().(*NilBytesColumnValue),
			want: TypeBytes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilBytesColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilBytesColumnValue_Clone(t *testing.T) {
	tests := []struct {
		name string
		n    *NilBytesColumnValue
		want ColumnValue
	}{
		{
			name: "1",
			n:    NewNilBytesColumnValue().(*NilBytesColumnValue),
			want: NewNilBytesColumnValue(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.Clone()
			if got == tt.n {
				t.Errorf("NilBytesColumnValue.Clone() = %p, n %p", got, tt.n)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilBytesColumnValue.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBytesColumnValue(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want ColumnValue
	}{
		{
			name: "1",
			args: args{
				s: "中文abc1234<>&*^%$",
			},
			want: NewBytesColumnValue([]byte("中文abc1234<>&*^%$")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBytesColumnValue([]byte(tt.args.s)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBytesColumnValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		s    *BytesColumnValue
		want ColumnType
	}{
		{
			name: "1",
			s:    NewBytesColumnValue([]byte("")).(*BytesColumnValue),
			want: TypeBytes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BytesColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesColumnValue_AsBool(t *testing.T) {
	tests := []struct {
		name    string
		s       *BytesColumnValue
		want    bool
		wantErr bool
	}{
		//1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
		{
			name: "1",
			s:    NewBytesColumnValue([]byte("1")).(*BytesColumnValue),
			want: true,
		},
		{
			name: "t",
			s:    NewBytesColumnValue([]byte("t")).(*BytesColumnValue),
			want: true,
		},
		{
			name: "T",
			s:    NewBytesColumnValue([]byte("T")).(*BytesColumnValue),
			want: true,
		},
		{
			name: "TRUE",
			s:    NewBytesColumnValue([]byte("TRUE")).(*BytesColumnValue),
			want: true,
		},
		{
			name: "true",
			s:    NewBytesColumnValue([]byte("true")).(*BytesColumnValue),
			want: true,
		},
		{
			name: "True",
			s:    NewBytesColumnValue([]byte("True")).(*BytesColumnValue),
			want: true,
		},
		{
			name: "0",
			s:    NewBytesColumnValue([]byte("0")).(*BytesColumnValue),
			want: false,
		},
		{
			name: "f",
			s:    NewBytesColumnValue([]byte("f")).(*BytesColumnValue),
			want: false,
		},
		{
			name: "FALSE",
			s:    NewBytesColumnValue([]byte("FALSE")).(*BytesColumnValue),
			want: false,
		},
		{
			name: "F",
			s:    NewBytesColumnValue([]byte("F")).(*BytesColumnValue),
			want: false,
		},
		{
			name: "false",
			s:    NewBytesColumnValue([]byte("false")).(*BytesColumnValue),
			want: false,
		},
		{
			name: "False",
			s:    NewBytesColumnValue([]byte("False")).(*BytesColumnValue),
			want: false,
		},
		{
			name:    "FAlse",
			s:       NewBytesColumnValue([]byte("FAlse")).(*BytesColumnValue),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AsBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("BytesColumnValue.AsBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BytesColumnValue.AsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesColumnValue_AsBigInt(t *testing.T) {
	tests := []struct {
		name    string
		s       *BytesColumnValue
		want    *big.Int
		wantErr bool
	}{
		{
			name: "1",
			s:    NewBytesColumnValue([]byte("12340000.0")).(*BytesColumnValue),
			want: testBigIntFromString("12340000"),
		},
		{
			name: "2",
			s:    NewBytesColumnValue([]byte("1234213213214135465736545425353980988.0")).(*BytesColumnValue),
			want: testBigIntFromString("1234213213214135465736545425353980988"),
		},
		{
			name: "3",
			s:    NewBytesColumnValue([]byte("-12340000.3")).(*BytesColumnValue),
			want: testBigIntFromString("-12340000"),
		},
		{
			name: "4",
			s:    NewBytesColumnValue([]byte("1.12345689")).(*BytesColumnValue),
			want: testBigIntFromString("1"),
		},
		{
			name: "5",
			s:    NewBytesColumnValue([]byte("1.23456e4")).(*BytesColumnValue),
			want: testBigIntFromString("12345"),
		},
		{
			name:    "6",
			s:       NewBytesColumnValue([]byte("1.23456e4dada")).(*BytesColumnValue),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AsBigInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("BytesColumnValue.AsBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Cmp(tt.want) != 0 {
				t.Errorf("BytesColumnValue.AsBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesColumnValue_AsDecimal(t *testing.T) {
	tests := []struct {
		name    string
		s       *BytesColumnValue
		want    decimal.Decimal
		wantErr bool
	}{
		{
			name: "1",
			s:    NewBytesColumnValue([]byte("12340000.0")).(*BytesColumnValue),
			want: testDecimalFormString("12340000"),
		},
		{
			name: "2",
			s:    NewBytesColumnValue([]byte("1234213213214135465736545425353980988.0")).(*BytesColumnValue),
			want: testDecimalFormString("1234213213214135465736545425353980988"),
		},
		{
			name: "3",
			s:    NewBytesColumnValue([]byte("-12340000.3")).(*BytesColumnValue),
			want: testDecimalFormString("-12340000.3"),
		},
		{
			name: "4",
			s:    NewBytesColumnValue([]byte("1.12345689")).(*BytesColumnValue),
			want: testDecimalFormString("1.12345689"),
		},
		{
			name: "5",
			s:    NewBytesColumnValue([]byte("1.23456e4")).(*BytesColumnValue),
			want: testDecimalFormString("12345.6"),
		},
		{
			name:    "6",
			s:       NewBytesColumnValue([]byte("1.23456e4dad")).(*BytesColumnValue),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AsDecimal()
			if (err != nil) != tt.wantErr {
				t.Errorf("BytesColumnValue.AsDecimal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("BytesColumnValue.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesColumnValue_AsString(t *testing.T) {
	tests := []struct {
		name    string
		s       *BytesColumnValue
		want    string
		wantErr bool
	}{
		{
			name: "1",
			s:    NewBytesColumnValue([]byte("中文abc1234<>&*^%$")).(*BytesColumnValue),
			want: "中文abc1234<>&*^%$",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("BytesColumnValue.AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BytesColumnValue.AsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesColumnValue_AsBytes(t *testing.T) {
	tests := []struct {
		name    string
		s       *BytesColumnValue
		want    []byte
		wantErr bool
	}{
		{
			name: "1",
			s:    NewBytesColumnValue([]byte("中文abc1234<>&*^%$")).(*BytesColumnValue),
			want: []byte("中文abc1234<>&*^%$"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AsBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("BytesColumnValue.AsBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BytesColumnValue.AsBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesColumnValue_AsTime(t *testing.T) {
	tests := []struct {
		name    string
		s       *BytesColumnValue
		wantT   time.Time
		wantErr bool
	}{
		{
			name:  "1",
			s:     NewBytesColumnValue([]byte(time.Date(2020, 12, 17, 22, 49, 56, 69-999-999, time.Local).Format(time.RFC3339Nano))).(*BytesColumnValue),
			wantT: time.Date(2020, 12, 17, 22, 49, 56, 69-999-999, time.Local),
		},
		{
			name:    "2",
			s:       NewBytesColumnValue([]byte("中文abc1234<>&*^%$")).(*BytesColumnValue),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := tt.s.AsTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("BytesColumnValue.AsTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !gotT.Equal(tt.wantT) {
				t.Errorf("BytesColumnValue.AsTime() = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}

func TestBytesColumnValue_String(t *testing.T) {
	tests := []struct {
		name string
		s    *BytesColumnValue
		want string
	}{
		{
			name: "1",
			s:    NewBytesColumnValue([]byte("中文abc1234<>&*^%$")).(*BytesColumnValue),
			want: "中文abc1234<>&*^%$",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("BytesColumnValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesColumnValue_Clone(t *testing.T) {
	tests := []struct {
		name string
		s    *BytesColumnValue
		want ColumnValue
	}{
		{
			name: "1",
			s:    NewBytesColumnValue([]byte("中文abc1234<>&*^%$")).(*BytesColumnValue),
			want: NewBytesColumnValue([]byte("中文abc1234<>&*^%$")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Clone()
			if got == tt.s {
				t.Errorf("BytesColumnValue.Clone() = %p, s %v", got, tt.s)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BytesColumnValue.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesColumnValue_Cmp(t *testing.T) {
	type args struct {
		right ColumnValue
	}
	tests := []struct {
		name    string
		b       *BytesColumnValue
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "1",
			b:    NewBytesColumnValue([]byte("123")).(*BytesColumnValue),
			args: args{
				right: NewNilBigIntColumnValue(),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "2",
			b:    NewBytesColumnValue([]byte("abc")).(*BytesColumnValue),
			args: args{
				right: NewBytesColumnValue([]byte("abcd")),
			},
			want:    -1,
			wantErr: false,
		},
		{
			name: "3",
			b:    NewBytesColumnValue([]byte("abc")).(*BytesColumnValue),
			args: args{
				right: NewBytesColumnValue([]byte("abd")),
			},
			want:    -1,
			wantErr: false,
		},
		{
			name: "4",
			b:    NewBytesColumnValue([]byte("abc")).(*BytesColumnValue),
			args: args{
				right: NewBytesColumnValue([]byte("abc")),
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "5",
			b:    NewBytesColumnValue([]byte("abcd")).(*BytesColumnValue),
			args: args{
				right: NewBytesColumnValue([]byte("abc")),
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "6",
			b:    NewBytesColumnValue([]byte("abd")).(*BytesColumnValue),
			args: args{
				right: NewBytesColumnValue([]byte("abc")),
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Cmp(tt.args.right)
			if (err != nil) != tt.wantErr {
				t.Errorf("BytesColumnValue.Cmp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BytesColumnValue.Cmp() = %v, want %v", got, tt.want)
			}
		})
	}
}
