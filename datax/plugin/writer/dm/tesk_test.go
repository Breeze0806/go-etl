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
