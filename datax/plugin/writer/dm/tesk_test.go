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

package dm

import (
	"testing"

	"github.com/Breeze0806/go-etl/datax/plugin/writer/dbms"
	"github.com/Breeze0806/go-etl/storage/database"
)

func TestExecMode(t *testing.T) {
	tests := []struct {
		name      string
		writeMode string
		expected  string
	}{
		{
			name:      "insert mode",
			writeMode: database.WriteModeInsert,
			expected:  dbms.ExecModeNormal,
		},
		{
			name:      "unknown mode",
			writeMode: "unknown_mode",
			expected:  dbms.ExecModeNormal,
		},
		{
			name:      "empty mode",
			writeMode: "",
			expected:  dbms.ExecModeNormal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := execMode(tt.writeMode)
			if result != tt.expected {
				t.Errorf("execMode(%s) = %s, expected %s", tt.writeMode, result, tt.expected)
			}
		})
	}
}

func TestExecModeMap(t *testing.T) {
	// 验证 execModeMap 的映射关系
	if mode, ok := execModeMap[database.WriteModeInsert]; !ok {
		t.Errorf("execModeMap should contain key: %s", database.WriteModeInsert)
	} else if mode != dbms.ExecModeNormal {
		t.Errorf("execModeMap[%s] = %s, expected %s", database.WriteModeInsert, mode, dbms.ExecModeNormal)
	}
}
