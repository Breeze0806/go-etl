package mysql

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

//Task 任务
type Task struct {
	*plugin.BaseTask

	db    *database.DBWrapper
	param *parameter
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

//StartRead 开始读
func (t *Task) StartRead(ctx context.Context, sender plugin.RecordSender) (err error) {
	handler := database.NewBaseFetchHandler(func() (element.Record, error) {
		return sender.CreateRecord()
	}, func(r element.Record) error {
		return sender.SendWriter(r)
	})

	param := newQueryParam(t.param)
	if err = t.db.FetchRecord(ctx, param, handler); err != nil {
		return
	}
	return nil
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

type queryParam struct {
	*parameter
}

func newQueryParam(p *parameter) *queryParam {
	return &queryParam{
		parameter: p,
	}
}

func (q *queryParam) Query(_ []element.Record) (string, error) {
	buf := bytes.NewBufferString("select ")
	if len(q.paramConfig.Column) == 0 {
		return "", errors.New("column is empty")
	}
	for i, v := range q.paramConfig.Column {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(v)
	}
	buf.WriteString(" from ")
	buf.WriteString(q.Table().Quoted())
	buf.WriteString(" where ")
	buf.WriteString(q.paramConfig.Where)
	return buf.String(), nil
}

func (q *queryParam) Agrs(_ []element.Record) ([]interface{}, error) {
	return nil, nil
}
