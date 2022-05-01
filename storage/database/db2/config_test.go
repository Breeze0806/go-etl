// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
