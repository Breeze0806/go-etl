package xlsx

import (
	"github.com/Breeze0806/go-etl/config"
	spireader "github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/file"
)

//Reader 读取器
type Reader struct {
	pluginConf *config.JSON
}

//ResourcesConfig 插件资源配置
func (r *Reader) ResourcesConfig() *config.JSON {
	return r.pluginConf
}

//Job 工作
func (r *Reader) Job() spireader.Job {
	job := NewJob()
	job.SetPluginConf(r.pluginConf)
	return job
}

//Task 任务
func (r *Reader) Task() spireader.Task {
	task := file.NewTask()
	task.SetPluginConf(r.pluginConf)
	return task
}
