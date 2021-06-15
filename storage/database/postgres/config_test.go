package postgres

import (
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
)

func testJSONFromString(s string) *config.JSON {
	json, err := config.NewJSONFromString(s)
	if err != nil {
		panic(err)
	}
	return json
}

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
				conf: testJSONFromString(`{
					"url" : "1",
					"username" : "2",
					"password": "3"
				}`),
			},
			wantC: &Config{
				URL:      "1",
				Username: "2",
				Password: "3",
			},
		},
		{
			name: "2",
			args: args{
				conf: testJSONFromString(`{
					"url" : 1,
					"username" : "1",
					"password": "1"
				}`),
			},
			wantErr: true,
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

func TestConfig_FormatDSN(t *testing.T) {
	tests := []struct {
		name    string
		c       *Config
		wantDsn string
		wantErr bool
	}{
		{
			name: "1",
			c: &Config{
				URL:      "postgres://hostname:5432/db",
				Username: "user",
				Password: "password",
			},
			wantDsn: "postgres://user:password@hostname:5432/db",
		},
		{
			name: "2",
			c: &Config{
				URL:      "postgres://hostname:5432/db",
				Username: "user",
				Password: "",
			},
			wantDsn: "postgres://user@hostname:5432/db",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDsn, err := tt.c.FormatDSN()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.FormatDSN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDsn != tt.wantDsn {
				t.Errorf("Config.FormatDSN() = %v, want %v", gotDsn, tt.wantDsn)
			}
		})
	}
}
