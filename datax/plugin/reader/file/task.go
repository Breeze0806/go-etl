package file

import (
	"context"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/stream/file"
)

type Task struct {
	*plugin.BaseTask

	streamer *file.InStreamer
}

func NewTask() *Task {
	return &Task{
		BaseTask: plugin.NewBaseTask(),
	}
}

func (t *Task) Init(ctx context.Context) (err error) {
	var name string
	if name, err = t.PluginConf().GetString("opener"); err != nil {
		return
	}
	var filename string
	if filename, err = t.PluginJobConf().GetString("path"); err != nil {
		return
	}

	if t.streamer, err = file.NewInStreamer(name, filename); err != nil {
		return
	}
	return
}

//Destroy 销毁
func (t *Task) Destroy(ctx context.Context) (err error) {
	if t.streamer != nil {
		err = t.streamer.Close()
	}
	return
}

type handler struct {
	sender plugin.RecordSender
}

func newHander(sender plugin.RecordSender) *handler {
	return &handler{
		sender: sender,
	}
}

func (h *handler) CreateRecord() (element.Record, error) {
	return h.sender.CreateRecord()
}

func (h *handler) OnRecord(r element.Record) error {
	return h.sender.SendWriter(r)
}

func (t *Task) StartRead(ctx context.Context, sender plugin.RecordSender) (err error) {
	handler := newHander(sender)

	log.Infof("jobid %v taskgroupid %v taskid %v startRead begin", t.JobID(), t.TaskGroupID(), t.TaskID())
	defer func() {
		sender.Terminate()
		log.Infof("jobid %v taskgroupid %v taskid %v startRead end", t.JobID(), t.TaskGroupID(), t.TaskID())
	}()
	var configs []*config.JSON
	configs, err = t.PluginJobConf().GetConfigArray("content")
	if err != nil {
		return err
	}

	for _, conf := range configs {
		if err = t.streamer.Read(ctx, conf, handler); err != nil {
			return
		}
	}
	return nil
}
