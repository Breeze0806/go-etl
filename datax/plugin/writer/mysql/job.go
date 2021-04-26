package mysql

import (
	"context"
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

//Job 工作
type Job struct {
	*plugin.BaseJob

	execer    Execer
	newExecer func(name string, conf *config.JSON) (Execer, error)
}

//Init 初始化
func (j *Job) Init(ctx context.Context) (err error) {
	var name string
	if name, err = j.PluginConf().GetString("dialect"); err != nil {
		return
	}

	var paramConfig *paramConfig
	if paramConfig, err = newParamConfig(j.PluginJobConf()); err != nil {
		return
	}

	var jobSettingConf *config.JSON
	if jobSettingConf, err = j.PluginJobConf().GetConfig(coreconst.DataxJobSetting); err != nil {
		jobSettingConf, _ = config.NewJSONFromString("{}")
		err = nil
	}
	if err = jobSettingConf.Set("username", paramConfig.Username); err != nil {
		return
	}

	if err = jobSettingConf.Set("password", paramConfig.Password); err != nil {
		return
	}

	if err = jobSettingConf.Set("url", paramConfig.Connection.URL); err != nil {
		return
	}

	if j.execer, err = j.newExecer(name, jobSettingConf); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err = j.execer.QueryContext(ctx, "select 1")
	if err != nil {
		return
	}
	return
}

//Destroy 销毁
func (j *Job) Destroy(ctx context.Context) (err error) {
	if j.execer != nil {
		err = j.execer.Close()
	}
	return
}

//Split 切分任务
func (j *Job) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	return []*config.JSON{j.PluginJobConf().CloneConfig()}, nil
}
