package element

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNilTimeColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		n    *NilTimeColumnValue
		want ColumnType
	}{
		{
			name: "1",
			n:    NewNilTimeColumnValue().(*NilTimeColumnValue),
			want: TypeTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilTimeColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilTimeColumnValue_clone(t *testing.T) {
	tests := []struct {
		name string
		n    *NilTimeColumnValue
		want ColumnValue
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilTimeColumnValue.clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeColumnValue_Type(t *testing.T) {
	tests := []struct {
		name string
		t    *TimeColumnValue
		want ColumnType
	}{
		{
			name: "1",
			t:    NewTimeColumnValue(time.Time{}).(*TimeColumnValue),
			want: TypeTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeColumnValue.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeColumnValue_AsBool(t *testing.T) {
	tests := []struct {
		name    string
		t       *TimeColumnValue
		want    bool
		wantErr bool
	}{
		{
			name:    "1",
			t:       NewTimeColumnValue(time.Time{}).(*TimeColumnValue),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.AsBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeColumnValue.AsBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TimeColumnValue.AsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeColumnValue_AsBigInt(t *testing.T) {
	tests := []struct {
		name    string
		t       *TimeColumnValue
		want    *big.Int
		wantErr bool
	}{
		{
			name:    "1",
			t:       NewTimeColumnValue(time.Time{}).(*TimeColumnValue),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.AsBigInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeColumnValue.AsBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeColumnValue.AsBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeColumnValue_AsDecimal(t *testing.T) {
	tests := []struct {
		name    string
		t       *TimeColumnValue
		want    decimal.Decimal
		wantErr bool
	}{
		{
			name:    "1",
			t:       NewTimeColumnValue(time.Time{}).(*TimeColumnValue),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.AsDecimal()
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeColumnValue.AsDecimal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeColumnValue.AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeColumnValue_AsString(t *testing.T) {
	tests := []struct {
		name    string
		t       *TimeColumnValue
		wantS   string
		wantErr bool
	}{
		{
			name:  "1",
			t:     NewTimeColumnValue(time.Date(2020, 12, 17, 22, 49, 56, 69-999-999, time.Local)).(*TimeColumnValue),
			wantS: time.Date(2020, 12, 17, 22, 49, 56, 69-999-999, time.Local).Format(defaultTimeFormat),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := tt.t.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeColumnValue.AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotS != tt.wantS {
				t.Errorf("TimeColumnValue.AsString() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestTimeColumnValue_AsBytes(t *testing.T) {
	tests := []struct {
		name    string
		t       *TimeColumnValue
		wantB   []byte
		wantErr bool
	}{
		{
			name:  "1",
			t:     NewTimeColumnValue(time.Date(2020, 12, 17, 22, 49, 56, 69-999-999, time.Local)).(*TimeColumnValue),
			wantB: []byte(time.Date(2020, 12, 17, 22, 49, 56, 69-999-999, time.Local).Format(defaultTimeFormat)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotB, err := tt.t.AsBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeColumnValue.AsBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("TimeColumnValue.AsBytes() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestTimeColumnValue_AsTime(t *testing.T) {
	tests := []struct {
		name    string
		t       *TimeColumnValue
		want    time.Time
		wantErr bool
	}{
		{
			name: "1",
			t:    NewTimeColumnValue(time.Date(2020, 12, 17, 22, 49, 56, 69-999-999, time.Local)).(*TimeColumnValue),
			want: time.Date(2020, 12, 17, 22, 49, 56, 69-999-999, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.AsTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeColumnValue.AsTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("TimeColumnValue.AsTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
