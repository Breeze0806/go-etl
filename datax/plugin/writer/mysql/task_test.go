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

package mysql

import (
	"testing"

	"github.com/Breeze0806/go-etl/datax/plugin/writer/dbms"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/Breeze0806/go-etl/storage/database/mysql"
)

type mockTable struct {
	*database.BaseTable

	n int
}

func newMockTable(n int) *mockTable {
	return &mockTable{
		BaseTable: database.NewBaseTable("db", "schema", "name"),
		n:         n,
	}
}

func (m *mockTable) Quoted() string {
	return ""
}

func (m *mockTable) Fields() []database.Field {
	var fields []database.Field
	for i := 0; i < m.n; i++ {
		fields = append(fields, nil)
	}
	return fields
}

func Test_batchWriter_BatchSize(t *testing.T) {
	tests := []struct {
		name     string
		b        *batchWriter
		wantSize int
	}{
		{
			name: "1",
			b: &batchWriter{
				BaseBatchWriter: dbms.NewBaseBatchWriter(&dbms.Task{
					Config: &dbms.BaseConfig{
						BatchSize: 1000,
					},
					Table: newMockTable(maxNumPlaceholder / 1000),
				}, "", nil),
			},
			wantSize: 1000,
		},
		{
			name: "1",
			b: &batchWriter{
				BaseBatchWriter: dbms.NewBaseBatchWriter(&dbms.Task{
					Config: &dbms.BaseConfig{
						BatchSize: 10000,
					},
					Table: newMockTable(32),
				}, "", nil),
			},
			wantSize: maxNumPlaceholder / 32,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSize := tt.b.BatchSize(); gotSize != tt.wantSize {
				t.Errorf("batchWriter.BatchSize() = %v, want %v", gotSize, tt.wantSize)
			}
		})
	}
}

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
				writeMode: database.WriteModeInsert,
			},
			want: dbms.ExecModeNormal,
		},
		{
			name: "2",
			args: args{
				writeMode: mysql.WriteModeReplace,
			},
			want: dbms.ExecModeNormal,
		},
		{
			name: "3",
			args: args{
				writeMode: "",
			},
			want: dbms.ExecModeNormal,
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
