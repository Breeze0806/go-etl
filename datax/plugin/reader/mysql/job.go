package mysql

import (
	"context"
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/storage/database"
)

//Job 工作
type Job struct {
	*plugin.BaseJob

	db *database.DBWrapper
}

//Init 初始化
func (j *Job) Init(ctx context.Context) (err error) {
	var name string
	if name, err = j.PluginConf().GetString("dialect"); err != nil {
		return
	}
	var paramConf *config.JSON
	if paramConf, err = j.PluginJobConf().GetConfig(coreconst.DataxJobContentReaderParameter); err != nil {
		return
	}

	var paramConfig *paramConfig
	if paramConfig, err = newParamConfig(paramConf); err != nil {
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

	if j.db, err = database.Open(name, jobSettingConf); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err = j.db.ExecContext(ctx, "select 1")
	if err != nil {
		return
	}
	return
}

//Destroy 销毁
func (j *Job) Destroy(ctx context.Context) (err error) {
	return j.db.Close()
}

//Split 切分
func (j *Job) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	return []*config.JSON{j.PluginJobConf().CloneConfig()}, nil
}
