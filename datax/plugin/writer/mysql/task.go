package mysql

import (
	"context"
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/rdbm"
	"github.com/Breeze0806/go-etl/storage/database"
)

//Task 任务
type Task struct {
	*writer.BaseTask

	execer    rdbm.Execer
	newExecer func(name string, conf *config.JSON) (rdbm.Execer, error)
	param     *parameter
}

//Init 初始化
func (t *Task) Init(ctx context.Context) (err error) {
	var name string
	if name, err = t.PluginConf().GetString("dialect"); err != nil {
		return
	}

	var paramConfig *paramConfig
	if paramConfig, err = newParamConfig(t.PluginJobConf()); err != nil {
		return
	}

	var jobSettingConf *config.JSON
	if jobSettingConf, err = t.PluginJobConf().GetConfig(coreconst.DataxJobSetting); err != nil {
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

	if t.execer, err = t.newExecer(name, jobSettingConf); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err = t.execer.QueryContext(ctx, "select 1")
	if err != nil {
		return
	}

	t.param = newParameter(paramConfig, t.execer)

	param := newTableParam(t.param)
	if _, err = t.execer.FetchTableWithParam(ctx, param); err != nil {
		return
	}

	return
}

//Destroy 销毁
func (t *Task) Destroy(ctx context.Context) (err error) {
	if t.execer != nil {
		err = t.execer.Close()
	}
	return
}

type batchWriter struct {
	*Task

	opts *database.ParameterOptions
}

func newBatchWriter(t *Task, opts *database.ParameterOptions) *batchWriter {
	return &batchWriter{
		Task: t,
		opts: opts,
	}
}

func (m *batchWriter) BatchSize() int {
	return m.param.paramConfig.getBatchSize()
}

func (m *batchWriter) BatchTimeout() time.Duration {
	return m.param.paramConfig.getBatchTimeout()
}

func (m *batchWriter) BatchWrite(ctx context.Context) error {
	return m.execer.BatchExec(ctx, m.opts)
}

func (m *batchWriter) Options() *database.ParameterOptions {
	return m.opts
}

//StartWrite 开始写
func (t *Task) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) (err error) {
	writer := newBatchWriter(t, &database.ParameterOptions{
		TxOptions: nil,
		Table:     t.param.Table(),
		Mode:      t.param.paramConfig.WriteMode,
	})
	return rdbm.StartWrite(ctx, writer, receiver)
}
