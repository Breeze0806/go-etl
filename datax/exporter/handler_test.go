package exporter

import (
	"net/http"
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
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		h    *Handler
		args args
	}{
		{
			name: "1",
			h: NewHandler(&datax.Engine{
				Container: &MockContainer{},
			}),
			args: args{
				w: &MockResponseWriter{
					header: make(http.Header),
				},
				r: &http.Request{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.ServeHTTP(tt.args.w, tt.args.r)
			t.Log(string((tt.args.w).(*MockResponseWriter).data))
		})
	}
}
