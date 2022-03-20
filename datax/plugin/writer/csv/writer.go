package csv

import (
	"github.com/Breeze0806/go-etl/config"
	spiwriter "github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/file"

	//csv storage
	_ "github.com/Breeze0806/go-etl/storage/stream/file/csv"
)

//Writer 写入器
type Writer struct {
	pluginConf *config.JSON
}

//ResourcesConfig 插件资源配置
func (w *Writer) ResourcesConfig() *config.JSON {
	return w.pluginConf
}

//Job 工作
func (w *Writer) Job() spiwriter.Job {
	job := NewJob()
	job.SetPluginConf(w.pluginConf)
	return job
}

//Task 任务
func (w *Writer) Task() spiwriter.Task {
	task := file.NewTask(func(conf *config.JSON) (file.Config, error) {
		c, err := file.NewBaseConfig(conf)
		if err != nil {
			return nil, err
		}
		return c, nil
	})
	task.SetPluginConf(w.pluginConf)
	return task
}
