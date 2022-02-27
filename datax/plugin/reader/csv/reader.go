package cvs

import (
	"github.com/Breeze0806/go-etl/config"
	spireader "github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/plugin/reader"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/file"

	_ "github.com/Breeze0806/go-etl/storage/stream/file/csv"
)

var _pluginConfig string

func init() {
	var err error
	if _pluginConfig, err = reader.RegisterReader(
		func(filename string) (reader.Reader, error) {
			return NewReader(filename)
		}); err != nil {
		panic(err)
	}
}

//Reader 读取器
type Reader struct {
	pluginConf *config.JSON
}

//NewReader 创建读取器
func NewReader(filename string) (r *Reader, err error) {
	r = &Reader{}
	r.pluginConf, err = config.NewJSONFromFile(filename)
	if err != nil {
		return nil, err
	}
	return
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
