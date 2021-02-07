package mysql

import (
	"reflect"
	"testing"
	"time"

	"github.com/Breeze0806/go/time2"
)

func Test_paramConfig_getBatchSize(t *testing.T) {
	tests := []struct {
		name string
		p    *paramConfig
		want int
	}{
		{
			name: "1",
			p:    &paramConfig{},
			want: defalutBatchSize,
		},
		{
			name: "2",
			p: &paramConfig{
				BatchSize: 100,
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.getBatchSize(); got != tt.want {
				t.Errorf("paramConfig.getBatchSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_paramConfig_getBatchTimeout(t *testing.T) {
	tests := []struct {
		name string
		p    *paramConfig
		want time.Duration
	}{
		{
			name: "1",
			p:    &paramConfig{},
			want: defalutBatchTimeout,
		},
		{
			name: "2",
			p: &paramConfig{
				BatchTimeout: time2.NewDuration(100 * time.Millisecond),
			},
			want: 100 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.getBatchTimeout(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paramConfig.getBatchTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}
