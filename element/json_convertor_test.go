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

package element

import (
	"testing"
)

func TestJsonConverter_ConvertFromString(t *testing.T) {
	converter := NewDefaultJSONConverter()
	tests := []struct {
		name    string
		s       string
		want    string
		wantErr bool
	}{
		{
			name: "1",
			s:    `{"a":1}`,
			want: `{"a":1}`,
		},
		{
			name: "2",
			s:    `[1,2,3]`,
			want: `[1,2,3]`,
		},
		{
			name:    "3",
			s:       `invalid json`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.ConvertFromString(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonConverter.ConvertFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ToString() != tt.want {
				t.Errorf("jsonConverter.ConvertFromString() = %v, want %v", got.ToString(), tt.want)
			}
		})
	}
}

func TestJsonConverter_ConvertFromBytes(t *testing.T) {
	converter := NewDefaultJSONConverter()
	tests := []struct {
		name    string
		b       []byte
		want    []byte
		wantErr bool
	}{
		{
			name: "1",
			b:    []byte(`{"a":1}`),
			want: []byte(`{"a":1}`),
		},
		{
			name: "2",
			b:    []byte(`[1,2,3]`),
			want: []byte(`[1,2,3]`),
		},
		{
			name:    "3",
			b:       []byte(`invalid json`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.ConvertFromBytes(tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonConverter.ConvertFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got.ToBytes()) != string(tt.want) {
				t.Errorf("jsonConverter.ConvertFromBytes() = %v, want %v", got.ToBytes(), tt.want)
			}
		})
	}
}
