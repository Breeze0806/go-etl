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
	"testing"

	"github.com/Breeze0806/go-etl/config"
	spiwriter "github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

func TestWriter_ResourcesConfig(t *testing.T) {
	tests := []struct {
		name string
		w    *Writer
		want *config.JSON
	}{
		{
			name: "1",
			w: &Writer{
				pluginConf: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.ResourcesConfig(); got != tt.want {
				t.Errorf("Writer.ResourcesConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriter_Job(t *testing.T) {
	type fields struct {
		pluginConf *config.JSON
	}
	tests := []struct {
		name   string
		fields fields
		want   spiwriter.Job
	}{
		{
			name: "1",
			fields: fields{
				pluginConf: nil,
			},
			want: &Job{
				Job: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Writer{
				pluginConf: tt.fields.pluginConf,
			}
			if got := w.Job(); got == nil {
				t.Errorf("Writer.Job() = %v, want non-nil", got)
			}
		})
	}
}

func TestWriter_Task(t *testing.T) {
	type fields struct {
		pluginConf *config.JSON
	}
	tests := []struct {
		name   string
		fields fields
		want   spiwriter.Task
	}{
		{
			name: "1",
			fields: fields{
				pluginConf: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Writer{
				pluginConf: tt.fields.pluginConf,
			}
			if got := w.Task(); got == nil {
				t.Errorf("Writer.Task() = %v, want non-nil", got)
			}
		})
	}
}
