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
	"fmt"
	"strings"
	"testing"
)

func TestTransformError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *TransformError
		want string
	}{
		{
			name: "1",
			e:    NewTransformErrorFormColumnTypes(TypeString, TypeBigInt, fmt.Errorf("test")),
			want: "test",
		},
		{
			name: "2",
			e:    NewTransformErrorFormColumnTypes(TypeString, TypeBigInt, (NewTransformErrorFormColumnTypes(TypeString, TypeBigInt, fmt.Errorf("test1")))),
			want: "test1",
		},

		{
			name: "3",
			e:    NewTransformErrorFormColumnTypes(TypeString, TypeBigInt, nil),
			want: "transform",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); !strings.Contains(got, tt.want) {
				t.Errorf("TransformError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *SetError
		want string
	}{
		{
			name: "1",
			e:    NewSetError(TypeString, TypeBigInt, fmt.Errorf("test")),
			want: "test",
		},
		{
			name: "2",
			e:    NewSetError(TypeString, TypeBigInt, NewSetError(TypeString, TypeBigInt, fmt.Errorf("test1"))),
			want: "test1",
		},
		{
			name: "3",
			e:    NewSetError(TypeString, TypeBigInt, NewSetError(TypeString, TypeBigInt, nil)),
			want: "set",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); !strings.Contains(got, tt.want) {
				t.Errorf("SetError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
