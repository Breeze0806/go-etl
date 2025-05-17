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

package exporter

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/datax"
	"github.com/Breeze0806/go-etl/datax/core/statistics/container"
	"github.com/Breeze0806/go-etl/datax/core/taskgroup"
	"github.com/Breeze0806/go-etl/datax/core/transport/channel"
)

type MockContainer struct{}

func (m *MockContainer) Start() error {
	return nil
}

func (m *MockContainer) Metrics() (metric *container.Metrics) {
	metric = container.NewMetrics()
	metric.Set("jobID", 10230)
	metric.Set("metrics.0", &TaskGroupMetric{
		TaskGroupID: 0,
		Metrics: []taskgroup.Stats{
			{
				TaskID: 0,
				Channel: channel.StatsJSON{
					TotalByte:   12345678,
					TotalRecord: 123456789,
					Byte:        1987,
					Record:      123,
				},
			},
			{
				TaskID: 1,
				Channel: channel.StatsJSON{
					TotalByte:   888887654321,
					TotalRecord: 98765432,
					Byte:        4312,
					Record:      311,
				},
			},
		},
	})
	return
}

type MockResponseWriter struct {
	header     http.Header
	statusCode int
	data       []byte
}

func (m *MockResponseWriter) Header() http.Header {

	return m.header
}

func (m *MockResponseWriter) Write(buf []byte) (int, error) {
	m.data = append(m.data, buf...)
	return len(buf), nil
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

func TestHandler_ServeHTTP(t *testing.T) {
	URL, _ := url.Parse("http://127.0.0.1:6080/metrics")
	URLJson, _ := url.Parse("http://127.0.0.1:6080/metrics?t=json")
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name     string
		h        *Handler
		args     args
		wantData []byte
	}{
		{
			name: "exporter",
			h: NewHandler(&datax.Engine{
				Container: &MockContainer{},
			}),
			args: args{
				w: &MockResponseWriter{
					header: make(http.Header),
				},
				r: &http.Request{
					URL: URL,
				},
			},
			wantData: []byte(`# HELP datax_channel_byte the number of bytes currently being synchronized in the channel
# TYPE datax_channel_byte gauge
datax_channel_byte{job_id="10230",task_group_id="0",task_id="0"} 1987
datax_channel_byte{job_id="10230",task_group_id="0",task_id="1"} 4312
# HELP datax_channel_record the number of records currently being synchronized in the channel
# TYPE datax_channel_record gauge
datax_channel_record{job_id="10230",task_group_id="0",task_id="0"} 123
datax_channel_record{job_id="10230",task_group_id="0",task_id="1"} 311
# HELP datax_channel_total_byte the total number of bytes synchronized
# TYPE datax_channel_total_byte counter
datax_channel_total_byte{job_id="10230",task_group_id="0",task_id="0"} 1.2345678e+07
datax_channel_total_byte{job_id="10230",task_group_id="0",task_id="1"} 8.88887654321e+11
# HELP datax_channel_total_record the total number of records synchronized
# TYPE datax_channel_total_record counter
datax_channel_total_record{job_id="10230",task_group_id="0",task_id="0"} 1.23456789e+08
datax_channel_total_record{job_id="10230",task_group_id="0",task_id="1"} 9.8765432e+07
`),
		},
		{
			name: "json",
			h: NewHandler(&datax.Engine{
				Container: &MockContainer{},
			}),
			args: args{
				w: &MockResponseWriter{
					header: make(http.Header),
				},
				r: &http.Request{
					URL: URLJson,
				},
			},
			wantData: []byte(`{
    "jobID": 10230,
    "metrics": [
        {
            "taskGroupID": 0,
            "metrics": [
                {
                    "taskID": 0,
                    "channel": {
                        "totalByte": 12345678,
                        "totalRecord": 123456789,
                        "byte": 1987,
                        "record": 123
                    }
                },
                {
                    "taskID": 1,
                    "channel": {
                        "totalByte": 888887654321,
                        "totalRecord": 98765432,
                        "byte": 4312,
                        "record": 311
                    }
                }
            ]
        }
    ]
}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.ServeHTTP(tt.args.w, tt.args.r)
			data := (tt.args.w).(*MockResponseWriter).data
			fmt.Println(string(data))
			if !reflect.DeepEqual(data, tt.wantData) {
				t.Errorf("Engine.Handler() data = %v, want %v", string(data), string(tt.wantData))
			}
		})
	}
}
