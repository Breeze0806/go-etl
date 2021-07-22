package rdbm

import (
	"context"
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

type Job struct {
	*plugin.BaseJob

	Querier Querier
	Handler DbHandler
}

func (j *Job) Init(ctx context.Context) (err error) {
	var name string
	if name, err = j.PluginConf().GetString("dialect"); err != nil {
		return
	}

	var conf Config
	if conf, err = j.Handler.Config(j.PluginJobConf()); err != nil {
		return
	}

	var jobSettingConf *config.JSON
	if jobSettingConf, err = j.PluginJobConf().GetConfig(coreconst.DataxJobSetting); err != nil {
		jobSettingConf, _ = config.NewJSONFromString("{}")
		err = nil
	}
	if err = jobSettingConf.Set("username", conf.GetUsername()); err != nil {
		return
	}

	if err = jobSettingConf.Set("password", conf.GetPassword()); err != nil {
		return
	}

	if err = jobSettingConf.Set("url", conf.GetURL()); err != nil {
		return
	}

	if j.Querier, err = j.Handler.Querier(name, jobSettingConf); err != nil {
		return
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = j.Querier.PingContext(timeoutCtx)
	if err != nil {
		return
	}
	return
}

//Destroy 销毁
func (j *Job) Destroy(ctx context.Context) (err error) {
	if j.Querier != nil {
		err = j.Querier.Close()
	}
	return
}

//Split 切分
func (j *Job) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	return []*config.JSON{j.PluginJobConf().CloneConfig()}, nil
}
