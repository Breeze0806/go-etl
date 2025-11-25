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

package parquet

import (
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/storage/stream/file/parquet"
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
			name: "Valid config with path array",
			args: args{
				conf: testJSONFromString(`{"path":["/tmp/test.parquet"],"column":[{"name":"id"},{"name":"name"}]}`),
			},
			want: &Config{
				Path:   []string{"/tmp/test.parquet"},
				Column: []parquet.Column{{Name: "id"}, {Name: "name"}},
			},
			wantErr: false,
		},
		{
			name: "Valid config with single path",
			args: args{
				conf: testJSONFromString(`{"path":"/tmp/test.parquet","column":[{"name":"id"},{"name":"name"}]}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Missing path",
			args: args{
				conf: testJSONFromString(`{"column":[{"name":"id"}]}`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := NewConfig(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
