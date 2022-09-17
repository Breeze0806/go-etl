package csv

import (
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		conf *config.JSON
	}
	tests := []struct {
		name    string
		args    args
		wantC   *Config
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				conf: testJSONFromString(`{"comment":1}`),
			},
			wantC:   nil,
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				conf: testJSONFromString(`{"path":[]}`),
			},
			wantC: &Config{
				Path: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewConfig(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewConfig() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
