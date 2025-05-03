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
			"number of total byte which datax push into channel",
			variableLabels,
			nil,
		),
		prometheus.NewDesc(
			"datax_channel_total_record",
			"number of total record which datax push into channel",
			variableLabels,
			nil,
		),
		prometheus.NewDesc(
			"datax_channel_byte",
			"number of byte which datax now in channel",
			variableLabels,
			nil,
		),
		prometheus.NewDesc(
			"datax_channel_record",
			"number of record which datax now in channel ",
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
