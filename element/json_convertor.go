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
	"github.com/Breeze0806/go/encoding"
)

var _DefaultJSONConverter JSONConverter = NewDefaultJSONConverter()

// JSONConverter is the JSON converter interface
type JSONConverter interface {
	ConvertFromString(s string) (json JSON, err error)
	ConvertFromBytes(b []byte) (json JSON, err error)
}

// DefaultJSONConverter implements JSONConverter interface
type DefaultJSONConverter struct{}

// NewDefaultJSONConverter creates a new JSON converter
func NewDefaultJSONConverter() JSONConverter {
	return &DefaultJSONConverter{}
}

// ConvertFromString converts a string to JSON
func (c *DefaultJSONConverter) ConvertFromString(s string) (JSON, error) {
	json, err := encoding.NewJSONFromString(s)
	if err != nil {
		return nil, err
	}
	return NewDefaultJSON(json), nil
}

// ConvertFromBytes converts bytes to JSON
func (c *DefaultJSONConverter) ConvertFromBytes(b []byte) (JSON, error) {
	json, err := encoding.NewJSONFromBytes(b)
	if err != nil {
		return nil, err
	}
	return NewDefaultJSON(json), nil
}
