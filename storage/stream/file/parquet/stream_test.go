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

package parquet

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

func TestWriter_Write(t *testing.T) {
	tmpDir := os.TempDir()
	type args struct {
		records  []element.Record
		outConf  *config.JSON
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				records: []element.Record{
					element.NewDefaultRecord(),
					element.NewDefaultRecord(),
				},
				filename: filepath.Join(tmpDir, "1.parquet"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(tt.args.filename)

			var creator Creator
			out, err := creator.Create(tt.args.filename)
			if err != nil {
				t.Fatal(err)
			}
			defer out.Close()
			w, err := out.Writer(tt.args.outConf)
			if err != nil {
				t.Fatal(err)
			}
			defer w.Close()
			defer w.Flush()
			for _, r := range tt.args.records {
				err = w.Write(r)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("writer.Write() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
