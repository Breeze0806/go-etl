package mysql

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

//Task 任务
type Task struct {
	*writer.BaseTask

	db           *database.DBWrapper
	recordTimout time.Duration
	param        *parameter
}

//Init 初始化
func (t *Task) Init(ctx context.Context) (err error) {
	var name string
	if name, err = t.PluginConf().GetString("dialect"); err != nil {
		return
	}
	var paramConf *config.JSON
	if paramConf, err = t.PluginJobConf().GetConfig(coreconst.DataxJobContentReaderParameter); err != nil {
		return
	}

	var paramConfig *paramConfig
	if paramConfig, err = newParamConfig(paramConf); err != nil {
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

	if t.db, err = database.Open(name, jobSettingConf); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err = t.db.ExecContext(ctx, "select 1")
	if err != nil {
		return
	}

	t.param = newParameter(paramConfig, t.db)

	param := newTableParam(t.param)
	if _, err = t.db.FetchTableWithParam(ctx, param); err != nil {
		return
	}

	return
}

//Destroy 销毁
func (t *Task) Destroy(ctx context.Context) (err error) {
	return t.db.Close()
}

//StartWrite 开始写
func (t *Task) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) (err error) {
	opts := &database.ParameterOptions{
		TxOptions: nil,
		Table:     t.param.Table(),
		Mode:      t.param.paramConfig.WriteMode,
	}
	recordChan := make(chan element.Record)
	var rerr error

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(recordChan)
		for {
			select {
			case <-ctx.Done():
			default:
			}
			var record element.Record
			record, rerr = receiver.GetFromReader()
			if rerr != nil {
				return
			}
			recordChan <- record
		}
	}()
	ticker := time.NewTicker(t.recordTimout)
	defer ticker.Stop()
	var records []element.Record
	for {
		select {
		case record, ok := <-recordChan:
			if ok {
				return rerr
			}
			records = append(records, record)
			opts.Records = records
			if len(records) >= 1000 {
				if err = t.db.BatchExec(ctx, opts); err != nil {
					return err
				}
				records = nil
			}
		case <-ticker.C:
			if err = t.db.BatchExec(ctx, opts); err != nil {
				return err
			}
			records = nil
		}
	}
}

type parameter struct {
	*database.BaseParam

	paramConfig *paramConfig
}

func newParameter(paramConfig *paramConfig, db *database.DBWrapper) *parameter {
	p := &parameter{
		BaseParam: database.NewBaseParam(db.Table(database.NewBaseTable(
			paramConfig.Connection.Table.Db, "", paramConfig.Connection.Table.Name)), nil),
	}
	p.paramConfig = paramConfig
	return nil
}

type tableParam struct {
	*parameter
}

func newTableParam(p *parameter) *tableParam {
	return &tableParam{
		parameter: p,
	}
}

func (t *tableParam) Query(_ []element.Record) (string, error) {
	buf := bytes.NewBufferString("select ")
	if len(t.paramConfig.Column) == 0 {
		return "", errors.New("column is empty")
	}
	for i, v := range t.paramConfig.Column {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(v)
	}
	buf.WriteString(" from ")
	buf.WriteString(t.Table().Quoted())
	buf.WriteString(" where 1 = 2")
	return buf.String(), nil
}

func (t *tableParam) Agrs(_ []element.Record) ([]interface{}, error) {
	return nil, nil
}
