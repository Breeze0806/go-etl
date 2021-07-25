package rdbm

import (
	"context"
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/storage/database"
)

//Task 任务
type Task struct {
	*writer.BaseTask

	Handler DbHandler
	Execer  Execer
	Config  Config
	Table   database.Table
}

func NewTask(handler DbHandler) *Task {
	return &Task{
		BaseTask: writer.NewBaseTask(),
		Handler:  handler,
	}
}

//Init 初始化
func (t *Task) Init(ctx context.Context) (err error) {
	var name string
	if name, err = t.PluginConf().GetString("dialect"); err != nil {
		return
	}

	if t.Config, err = t.Handler.Config(t.PluginJobConf()); err != nil {
		return
	}

	var jobSettingConf *config.JSON
	if jobSettingConf, err = t.PluginJobConf().GetConfig(coreconst.DataxJobSetting); err != nil {
		jobSettingConf, _ = config.NewJSONFromString("{}")
		err = nil
	}
	if err = jobSettingConf.Set("username", t.Config.GetUsername()); err != nil {
		return
	}

	if err = jobSettingConf.Set("password", t.Config.GetPassword()); err != nil {
		return
	}

	if err = jobSettingConf.Set("url", t.Config.GetURL()); err != nil {
		return
	}

	if t.Execer, err = t.Handler.Execer(name, jobSettingConf); err != nil {
		return
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = t.Execer.PingContext(timeoutCtx)
	if err != nil {
		return
	}

	param := t.Handler.TableParam(t.Config, t.Execer)
	if t.Table, err = t.Execer.FetchTableWithParam(ctx, param); err != nil {
		return
	}

	return
}

//Destroy 销毁
func (t *Task) Destroy(ctx context.Context) (err error) {
	if t.Execer != nil {
		err = t.Execer.Close()
	}
	return
}
