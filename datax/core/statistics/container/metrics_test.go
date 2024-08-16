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

package container

import (
	"reflect"
	"testing"

	"github.com/Breeze0806/go/encoding"
)

func testJSON(json string) *encoding.JSON {
	j, _ := encoding.NewJSONFromString(json)
	return j
}

func TestMetrics_JSON(t *testing.T) {
	tests := []struct {
		name string
		m    *Metrics
		want *encoding.JSON
	}{
		{
			name: "1",
			m: &Metrics{
				metricJSON: testJSON(`{"test":"metrics"}`),
			},
			want: testJSON(`{"test":"metrics"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.JSON(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Metrics.JSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMetrics(t *testing.T) {
	tests := []struct {
		name string
		want *Metrics
	}{
		{
			name: "1",
			want: &Metrics{
				metricJSON: testJSON(`{}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_Set(t *testing.T) {
	type args struct {
		path  string
		value any
	}
	tests := []struct {
		name    string
		m       *Metrics
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			m:    NewMetrics(),
			args: args{
				path:  "path",
				value: "value",
			},
			want:    "value",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Set(tt.args.path, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Metrics.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got, _ := tt.m.JSON().GetString(tt.args.path); got != tt.want {
				t.Errorf("Metrics.Set() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		m    *Metrics
		args args
		want *encoding.JSON
	}{
		{
			name: "1",
			m: &Metrics{
				metricJSON: testJSON(`{"test":{"path":"value"}}`),
			},
			args: args{
				key: "test",
			},
			want: testJSON(`{"path":"value"}`),
		},
		{
			name: "1",
			m: &Metrics{
				metricJSON: testJSON(`{"test":{"path":"value"}}`),
			},
			args: args{
				key: "test.path",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Get(tt.args.key)
			if tt.want == nil && got == nil {
				return
			}
			if got.String() != tt.want.String() {
				t.Errorf("Metrics.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
