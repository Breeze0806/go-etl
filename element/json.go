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

// JSON is the JSON interface
type JSON interface {
	ToString() string
	ToBytes() []byte
	Clone() JSON
}

// DefaultJSON is a wrapper for encoding.JSON to implement element.JSON interface
type DefaultJSON struct {
	json *encoding.JSON
}

// NewDefaultJSON creates a new JSON element from encoding.JSON
func NewDefaultJSON(j *encoding.JSON) JSON {
	return &DefaultJSON{json: j}
}

// ToString returns the string representation of the JSON
func (j *DefaultJSON) ToString() string {
	if j.json == nil {
		return ""
	}
	return j.json.String()
}

// ToBytes returns the byte representation of the JSON
func (j *DefaultJSON) ToBytes() []byte {
	if j.json == nil {
		return nil
	}
	b, _ := j.json.MarshalJSON()
	return b
}

// GetJSON returns the underlying encoding.JSON
func (j *DefaultJSON) GetJSON() *encoding.JSON {
	return j.json
}

// Clone clones the JSON
func (j *DefaultJSON) Clone() JSON {
	if j.json == nil {
		return &DefaultJSON{json: nil}
	}
	return &DefaultJSON{json: j.json.Clone()}
}
