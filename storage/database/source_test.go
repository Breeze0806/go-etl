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
		want *config.JSON
	}{
		{
			name: "1",
			b:    NewBaseSource(testJSONFromString(`{}`)),
			want: testJSONFromString(`{}`),
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

func TestNewSource(t *testing.T) {
	registerMock()
	type args struct {
		name string
		conf *config.JSON
	}
	tests := []struct {
		name       string
		args       args
		wantSource Source
		wantErr    bool
	}{
		{
			name: "1",
			args: args{
				name: "mock",
				conf: testJSONFromString("{}"),
			},
			wantSource: &mockSource{
				BaseSource: NewBaseSource(testJSONFromString("{}")),
				name:       "mock",
			},
		},
		{
			name: "2",
			args: args{
				name: "test?",
				conf: testJSONFromString("{}"),
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				name: "mockErr",
				conf: testJSONFromString("{}"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSource, err := NewSource(tt.args.name, tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSource, tt.wantSource) {
				t.Errorf("NewSource() = %v, want %v", gotSource, tt.wantSource)
			}
		})
	}
}
