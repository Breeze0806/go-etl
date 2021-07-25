package rdbm

import (
	"reflect"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
	rdbmreader "github.com/Breeze0806/go-etl/datax/plugin/reader/rdbm"
)

func testBaseConfig(conf *config.JSON) (bc *BaseConfig) {
	var err error
	bc, err = NewBaseConfig(conf)
	if err != nil {
		panic(err)
	}
	return bc
}

func TestBaseConfig_GetColumns(t *testing.T) {
	tests := []struct {
		name        string
		b           *BaseConfig
		wantColumns []rdbmreader.Column
	}{
		{
			name: "1",
			b: &BaseConfig{
				Column: []string{"f1", "f2", "f3", "f4"},
			},
			wantColumns: []rdbmreader.Column{
				&rdbmreader.BaseColumn{
					Name: "f1",
				},
				&rdbmreader.BaseColumn{
					Name: "f2",
				},
				&rdbmreader.BaseColumn{
					Name: "f3",
				},
				&rdbmreader.BaseColumn{
					Name: "f4",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotColumns := tt.b.GetColumns(); !reflect.DeepEqual(gotColumns, tt.wantColumns) {
				t.Errorf("BaseConfig.GetColumns() = %v, want %v", gotColumns, tt.wantColumns)
			}
		})
	}
}

func TestBaseConfig_GetBatchTimeout(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseConfig
		want time.Duration
	}{
		{
			name: "1",
			b:    testBaseConfig(TestJSONFromString("{}")),
			want: defalutBatchTimeout,
		},
		{
			name: "2",
			b:    testBaseConfig(TestJSONFromString(`{"batchTimeout":"100ms"}`)),
			want: 100 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.GetBatchTimeout(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseConfig.GetBatchTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseConfig_GetBatchSize(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseConfig
		want int
	}{
		{
			name: "1",
			b:    testBaseConfig(TestJSONFromString("{}")),
			want: defalutBatchSize,
		},

		{
			name: "2",
			b:    testBaseConfig(TestJSONFromString(`{"batchSize":30000}`)),
			want: 30000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.GetBatchSize(); got != tt.want {
				t.Errorf("BaseConfig.GetBatchSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
