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

package plugin

// Type: Plugin Type
type Type string

// Plugin Type Enumeration
var (
	Reader      Type = "reader"      // Reader
	Writer      Type = "writer"      // Writer
	Transformer Type = "transformer" // Converter
	Handler     Type = "handler"     // Processor
)

// NewType: New Type
func NewType(s string) Type {
	return Type(s)
}

func (t Type) String() string {
	return string(t)
}

// IsValid: Is Valid
func (t Type) IsValid() bool {
	switch t {
	case Reader, Writer, Transformer, Handler:
		return true
	}
	return false
}
