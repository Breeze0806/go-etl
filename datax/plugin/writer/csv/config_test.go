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
		want    *Config
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				conf: testJSONFromString(`{"encoding":1}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				conf: testJSONFromString(`{"path":[]}`),
			},
			want: &Config{
				Path: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
