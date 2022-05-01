package db2

import "testing"

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
				URL:      "HOSTNAME=192.168.0.1;PORT=50000;DATABASE=testdb",
				Username: "user",
				Password: "passwd",
			},
			wantDsn: "DATABASE=testdb;HOSTNAME=192.168.0.1;PORT=50000;PWD=passwd;UID=user",
		},
		{
			name: "2",
			c: &Config{
				URL:      "HOSTNAME =192.168.0.1;PORT= 50000;DATABASE=testdb",
				Username: "user",
				Password: "passwd",
			},
			wantDsn: "DATABASE=testdb;HOSTNAME=192.168.0.1;PORT=50000;PWD=passwd;UID=user",
		},
		{
			name: "3",
			c: &Config{
				URL:      "PORT=50000;DATABASE=testdb",
				Username: "user",
				Password: "passwd",
			},
			wantErr: true,
		},
		{
			name: "4",
			c: &Config{
				URL:      "HOSTNAME=192.168.0.1;PORT=50000",
				Username: "user",
				Password: "passwd",
			},
			wantErr: true,
		},
		{
			name: "5",
			c: &Config{
				URL:      "HOSTNAME =192.168.0.1;=;DATABASE=testdb",
				Username: "user",
				Password: "passwd",
			},
			wantErr: true,
		},
		{
			name: "6",
			c: &Config{
				URL:      "HOSTNAME =192.168.0.1; testdb",
				Username: "user",
				Password: "passwd",
			},
			wantErr: true,
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
