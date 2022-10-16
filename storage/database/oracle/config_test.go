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

package oracle

import (
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/godror/godror"
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

func TestConfig_FetchConnectionParams(t *testing.T) {
	tests := []struct {
		name    string
		c       *Config
		want    string
		wantErr bool
	}{
		{
			name: "1",
			c: &Config{
				URL:      "oracle://salesserver1/sales.us.example.com&poolSessionTimeout=42s",
				Username: "scott",
				Password: "tiger",
			},
			want: `user=scott password=tiger connectString="salesserver1/sales.us.example.com&poolSessionTimeout=42s"
		configDir= connectionClass= enableEvents=0 externalAuth=1 heterogeneousPool=0
		libDir= newPassword= noTimezoneCheck=0 poolIncrement=1 poolMaxSessions=1000
		poolMinSessions=1 poolSessionMaxLifetime=1h0m0s poolSessionTimeout=5m0s poolWaitTimeout=30s
		prelim=0 standaloneConnection=0 sysasm=0 sysdba=0 sysoper=0 timezone=`,
		},
		{
			name: "2",
			c: &Config{
				URL:      `connectString="salesserver1/sales.us.example.com&poolSessionTimeout=42s`,
				Username: "scott",
				Password: "tiger",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCon, err := tt.c.FetchConnectionParams()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.FetchConnectionParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			wantCon, _ := godror.ParseDSN(tt.want)
			if !reflect.DeepEqual(gotCon, wantCon) {
				t.Errorf("Config.FetchConnectionParams() = %v, want %v", gotCon.StringWithPassword(), wantCon)
			}
		})
	}
}
