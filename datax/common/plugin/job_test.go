package plugin

import (
	"reflect"
	"testing"
)

type mockJobCollector struct {
}

func (m *mockJobCollector) MessageMap() map[string][]string {
	return nil
}

func (m *mockJobCollector) MessageByKey(key string) []string {
	return nil
}

func TestBaseJob_SetCollector(t *testing.T) {
	type args struct {
		collector JobCollector
	}
	tests := []struct {
		name string
		b    *BaseJob
		args args
		want JobCollector
	}{
		{
			name: "1",
			b:    NewBaseJob(),
			args: args{
				collector: &mockJobCollector{},
			},
			want: &mockJobCollector{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetCollector(tt.args.collector)
			if !reflect.DeepEqual(tt.b.Collector(), tt.want) {
				t.Errorf("Collector() = %p want %p", tt.b.Collector(), tt.want)
			}
		})
	}
}
