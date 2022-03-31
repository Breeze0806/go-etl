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

package writer

import "testing"

func TestBaseTask_SupportFailOver(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseTask
		want bool
	}{
		{
			name: "1",
			b:    NewBaseTask(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.SupportFailOver(); got != tt.want {
				t.Errorf("BaseTask.SupportFailOver() = %v, want %v", got, tt.want)
			}
		})
	}
}
