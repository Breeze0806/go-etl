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

package oracle

import (
	"testing"

	"github.com/Breeze0806/go-etl/datax/plugin/writer/rdbm"
	"github.com/Breeze0806/go-etl/storage/database/oracle"
)

func Test_execMode(t *testing.T) {
	type args struct {
		writeMode string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				writeMode: oracle.WriteModeInsert,
			},
			want: rdbm.ExecModeStmt,
		},
		{
			name: "2",
			args: args{
				writeMode: "",
			},
			want: rdbm.ExecModeStmt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := execMode(tt.args.writeMode); got != tt.want {
				t.Errorf("execMode() = %v, want %v", got, tt.want)
			}
		})
	}
}
