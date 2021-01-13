package database

import (
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
)

func TestBaseSource_Config(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseSource
		want *config.Json
	}{
		{
			name: "1",
			b:    NewBaseSource(testJsonFromString(`{}`)),
			want: testJsonFromString(`{}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Config(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSource.Config() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetMaxOpenConns(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want int
	}{
		{
			name: "1",
			c:    &Config{},
			want: DefaultMaxOpenConns,
		},
		{
			name: "2",
			c: &Config{
				MaxOpenConns: 10,
			},
			want: 10,
		},
		{
			name: "3",
			c: &Config{
				MaxOpenConns: -10,
			},
			want: DefaultMaxOpenConns,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetMaxOpenConns(); got != tt.want {
				t.Errorf("Config.GetMaxOpenConns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetMaxIdleConns(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want int
	}{
		{
			name: "1",
			c:    &Config{},
			want: DefaultMaxIdleConns,
		},
		{
			name: "2",
			c: &Config{
				MaxIdleConns: -10,
			},
			want: DefaultMaxIdleConns,
		},
		{
			name: "3",
			c: &Config{
				MaxIdleConns: 10,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetMaxIdleConns(); got != tt.want {
				t.Errorf("Config.GetMaxIdleConns() = %v, want %v", got, tt.want)
			}
		})
	}
}
