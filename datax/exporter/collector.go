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
	"strconv"

	"github.com/Breeze0806/go-etl/datax/core/taskgroup"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	variableLabels = []string{
		"job_id",
		"task_group_id",
		"task_id",
	}
	PrometheusDescs = []*prometheus.Desc{
		prometheus.NewDesc(
			"datax_channel_total_byte",
			"the total number of bytes synchronized",
			variableLabels,
			nil,
		),
		prometheus.NewDesc(
			"datax_channel_total_record",
			"the total number of records synchronized",
			variableLabels,
			nil,
		),
		prometheus.NewDesc(
			"datax_channel_byte",
			"the number of bytes currently being synchronized in the channel",
			variableLabels,
			nil,
		),
		prometheus.NewDesc(
			"datax_channel_record",
			"the number of records currently being synchronized in the channel",
			variableLabels,
			nil,
		),
	}
)

type JobMetric struct {
	JobID   int64             `json:"jobID"`
	Metrics []TaskGroupMetric `json:"metrics"`
}

type TaskGroupMetric struct {
	TaskGroupID int64             `json:"taskGroupID"`
	Metrics     []taskgroup.Stats `json:"metrics"`
}

type JSONMetricCollector struct {
	Metric *JobMetric
}

func (mc JSONMetricCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, d := range PrometheusDescs {
		ch <- d
	}
}

func (mc JSONMetricCollector) Collect(ch chan<- prometheus.Metric) {
	jobID := strconv.FormatInt(mc.Metric.JobID, 10)
	for _, vi := range mc.Metric.Metrics {
		taskGroupID := strconv.FormatInt(vi.TaskGroupID, 10)
		for _, vj := range vi.Metrics {
			taskID := strconv.FormatInt(vj.TaskID, 10)
			ch <- prometheus.MustNewConstMetric(
				PrometheusDescs[0],
				prometheus.CounterValue,
				float64(vj.Channel.TotalByte),
				jobID, taskGroupID, taskID)

			ch <- prometheus.MustNewConstMetric(
				PrometheusDescs[1],
				prometheus.CounterValue,
				float64(vj.Channel.TotalRecord),
				jobID, taskGroupID, taskID)
			ch <- prometheus.MustNewConstMetric(
				PrometheusDescs[2],
				prometheus.GaugeValue,
				float64(vj.Channel.Byte),
				jobID, taskGroupID, taskID)
			ch <- prometheus.MustNewConstMetric(
				PrometheusDescs[3],
				prometheus.GaugeValue,
				float64(vj.Channel.Record),
				jobID, taskGroupID, taskID)
		}
	}
}
