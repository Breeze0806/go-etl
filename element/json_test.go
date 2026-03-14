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

	"github.com/Breeze0806/go/encoding"
)

func TestJsonWrapper_ToString(t *testing.T) {
	tests := []struct {
		name string
		json *encoding.JSON
		want string
	}{
		{
			name: "1",
			json: func() *encoding.JSON {
				j, _ := encoding.NewJSONFromString(`{"a":1}`)
				return j
			}(),
			want: `{"a":1}`,
		},
		{
			name: "2",
			json: func() *encoding.JSON {
				j, _ := encoding.NewJSONFromString(`[1,2,3]`)
				return j
			}(),
			want: `[1,2,3]`,
		},
		{
			name: "3",
			json: nil,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewDefaultJSON(tt.json)
			if got := j.ToString(); got != tt.want {
				t.Errorf("jsonWrapper.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonWrapper_ToBytes(t *testing.T) {
	tests := []struct {
		name string
		json *encoding.JSON
		want []byte
	}{
		{
			name: "1",
			json: func() *encoding.JSON {
				j, _ := encoding.NewJSONFromString(`{"a":1}`)
				return j
			}(),
			want: []byte(`{"a":1}`),
		},
		{
			name: "2",
			json: func() *encoding.JSON {
				j, _ := encoding.NewJSONFromString(`[1,2,3]`)
				return j
			}(),
			want: []byte(`[1,2,3]`),
		},
		{
			name: "3",
			json: nil,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewDefaultJSON(tt.json)
			if got := j.ToBytes(); string(got) != string(tt.want) {
				t.Errorf("jsonWrapper.ToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultJSON_Clone(t *testing.T) {
	tests := []struct {
		name string
		json *encoding.JSON
		want string
	}{
		{
			name: "1",
			json: func() *encoding.JSON {
				j, _ := encoding.NewJSONFromString(`{"a":1}`)
				return j
			}(),
			want: `{"a":1}`,
		},
		{
			name: "2",
			json: func() *encoding.JSON {
				j, _ := encoding.NewJSONFromString(`[1,2,3]`)
				return j
			}(),
			want: `[1,2,3]`,
		},
		{
			name: "3",
			json: nil,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewDefaultJSON(tt.json)
			got := j.Clone()
			if got == j {
				t.Errorf("DefaultJSON.Clone() = %p, j %p", got, j)
			}
			if got.ToString() != tt.want {
				t.Errorf("DefaultJSON.Clone() = %v, want %v", got.ToString(), tt.want)
			}
		})
	}
}

func TestDefaultJSON_GetJSON(t *testing.T) {
	tests := []struct {
		name string
		json *encoding.JSON
		want *encoding.JSON
	}{
		{
			name: "1",
			json: func() *encoding.JSON {
				j, _ := encoding.NewJSONFromString(`{"a":1}`)
				return j
			}(),
			want: func() *encoding.JSON {
				j, _ := encoding.NewJSONFromString(`{"a":1}`)
				return j
			}(),
		},
		{
			name: "2",
			json: nil,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := NewDefaultJSON(tt.json).(*DefaultJSON)
			got := j.GetJSON()
			if tt.want == nil {
				if got != nil {
					t.Errorf("DefaultJSON.GetJSON() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Errorf("DefaultJSON.GetJSON() = nil, want %v", tt.want)
				return
			}
			if got.String() != tt.want.String() {
				t.Errorf("DefaultJSON.GetJSON() = %v, want %v", got.String(), tt.want.String())
			}
		})
	}
}
