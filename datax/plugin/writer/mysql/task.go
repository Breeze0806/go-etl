package mysql

import (
	"context"
	"sync"
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/core/transport/exchange"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

//Task 任务
type Task struct {
	*writer.BaseTask

	execer    Execer
	newExecer func(name string, conf *config.JSON) (Execer, error)
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

//StartWrite 开始写
func (t *Task) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) (err error) {
	opts := &database.ParameterOptions{
		TxOptions: nil,
		Table:     t.param.Table(),
		Mode:      t.param.paramConfig.WriteMode,
	}
	recordChan := make(chan element.Record)
	var rerr error
	afterCtx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			close(recordChan)
			log.Debugf("jobID: %v taskgroupID:%v taskID: %v get records end", t.JobID(), t.TaskGroupID(), t.TaskID())
		}()
		log.Debugf("jobID: %v taskgroupID:%v taskID: %v start to get records", t.JobID(), t.TaskGroupID(), t.TaskID())
		for {
			select {
			case <-afterCtx.Done():
				return
			default:
			}
			var record element.Record
			record, rerr = receiver.GetFromReader()
			if rerr != nil && rerr != exchange.ErrEmpty {
				return
			}

			if rerr != exchange.ErrEmpty {
				select {
				case <-afterCtx.Done():
					return
				case recordChan <- record:
				}

			}
		}
	}()
	ticker := time.NewTicker(t.param.paramConfig.getBatchTimeout())
	defer ticker.Stop()
	var records []element.Record
	log.Debugf("jobID: %v taskgroupID:%v taskID: %v  start to BatchExec", t.JobID(), t.TaskGroupID(), t.TaskID())
	for {
		select {
		case record, ok := <-recordChan:
			if !ok {
				opts.Records = records
				if err = t.execer.BatchExec(ctx, opts); err != nil {
					log.Errorf("jobID: %v taskgroupID:%v taskID: %v BatchExec error: %v", t.JobID(), t.TaskGroupID(), t.TaskID(), err)
				}
				opts.Records = nil
				err = rerr
				goto End
			}
			records = append(records, record)

			if len(records) >= t.param.paramConfig.getBatchSize() {
				opts.Records = records
				if err = t.execer.BatchExec(ctx, opts); err != nil {
					log.Errorf("jobID: %v taskgroupID:%v taskID: %v BatchExec error: %v", t.JobID(), t.TaskGroupID(), t.TaskID(), err)
					goto End
				}
				opts.Records = nil
				records = nil
			}
		case <-ticker.C:
			opts.Records = records
			if err = t.execer.BatchExec(ctx, opts); err != nil {
				log.Errorf("jobID: %v taskgroupID:%v taskID: %v BatchExec error: %v", t.JobID(), t.TaskGroupID(), t.TaskID(), err)
				goto End
			}
			opts.Records = nil
			records = nil
		}
	}
End:
	cancel()
	log.Debugf("jobID: %v taskgroupID:%v taskID: %v wait all goroutine", t.JobID(), t.TaskGroupID(), t.TaskID())
	wg.Wait()
	log.Debugf("jobID: %v taskgroupID:%v taskID: %v wait all goroutine end", t.JobID(), t.TaskGroupID(), t.TaskID())
	switch {
	case ctx.Err() != nil:
		return nil
	case err == exchange.ErrTerminate:
		return nil
	}
	return
}
