package file

import (
	"context"
	"sync"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/core/transport/exchange"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/stream/file"
)

//Task 任务
type Task struct {
	*writer.BaseTask

	streamer  *file.OutStreamer
	conf      Config
	newConfig func(conf *config.JSON) (Config, error)
	content   *config.JSON
}

func NewTask(newConfig func(conf *config.JSON) (Config, error)) *Task {
	return &Task{
		BaseTask:  writer.NewBaseTask(),
		newConfig: newConfig,
	}
}

func (t *Task) Init(ctx context.Context) (err error) {
	var name string
	if name, err = t.PluginConf().GetString("creater"); err != nil {
		return
	}
	var filename string
	if filename, err = t.PluginJobConf().GetString("path"); err != nil {
		return
	}

	if t.content, err = t.PluginJobConf().GetConfig("content"); err != nil {
		return
	}

	if t.conf, err = t.newConfig(t.content); err != nil {
		return
	}

	if t.streamer, err = file.NewOutStreamer(name, filename); err != nil {
		return
	}

	return
}

//Destroy 销毁
func (t *Task) Destroy(ctx context.Context) (err error) {
	if t.streamer != nil {
		err = t.streamer.Close()
	}
	return
}

func (t *Task) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) (err error) {
	var sw file.StreamWriter
	if sw, err = t.streamer.Writer(t.content); err != nil {
		return
	}

	recordChan := make(chan element.Record)
	var rerr error
	afterCtx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	//通过该携程读取记录接受器receiver的记录放入recordChan
	go func() {
		defer func() {
			wg.Done()
			//关闭recordChan
			close(recordChan)
			log.Debugf("jobID: %v taskgroupID:%v taskID: %v get records end",
				t.JobID(), t.TaskGroupID(), t.TaskID())
		}()
		log.Debugf("jobID: %v taskgroupID:%v taskID: %v start to get records",
			t.JobID(), t.TaskGroupID(), t.TaskID())
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

			//当记录接受器receiver返回不为空错误，写入recordChan
			if rerr != exchange.ErrEmpty {
				select {
				//防止在ctx关闭时不写入recordChan
				case <-afterCtx.Done():
					return
				case recordChan <- record:
				}

			}
		}
	}()
	ticker := time.NewTicker(t.conf.GetBatchTimeout())
	defer ticker.Stop()
	cnt := 0
	log.Debugf("jobID: %v taskgroupID: %v taskID: %v  start to write",
		t.JobID(), t.TaskGroupID(), t.TaskID())
	for {
		select {
		case record, ok := <-recordChan:
			if !ok {
				//当写入结束时，将剩余的记录写入数据库
				if cnt > 0 {
					if err = sw.Flush(); err != nil {
						log.Errorf("jobID: %v taskgroupID:%v taskID: %v Flush error: %v",
							t.JobID(), t.TaskGroupID(), t.TaskID(), err)
					}
				}
				err = rerr
				goto End
			}

			if err = sw.Write(record); err != nil {
				log.Errorf("jobID: %v taskgroupID:%v taskID: %v Write error: %v",
					t.JobID(), t.TaskGroupID(), t.TaskID(), err)
				goto End
			}

			//当数据量超过单次批量数时 写入数据库
			if cnt >= t.conf.GetBatchSize() {
				if err = sw.Flush(); err != nil {
					log.Errorf("jobID: %v taskgroupID:%v taskID: %v Flush error: %v",
						t.JobID(), t.TaskGroupID(), t.TaskID(), err)
					goto End
				}
				cnt = 0
			}
		//当写入数据未达到单次批量数，超时也写入
		case <-ticker.C:
			if cnt > 0 {
				if err = sw.Flush(); err != nil {
					log.Errorf("jobID: %v taskgroupID:%v taskID: %v Flush error: %v",
						t.JobID(), t.TaskGroupID(), t.TaskID(), err)
					goto End
				}
			}
			cnt = 0
		}
	}
End:
	if err = sw.Close(); err != nil {
		log.Errorf("jobID: %v taskgroupID:%v taskID: %v Close error: %v",
			t.JobID(), t.TaskGroupID(), t.TaskID(), err)
	}
	cancel()
	log.Debugf("jobID: %v taskgroupID:%v taskID: %v wait all goroutine",
		t.JobID(), t.TaskGroupID(), t.TaskID())
	//等待携程结束
	wg.Wait()
	log.Debugf("jobID: %v taskgroupID:%v taskID: %v wait all goroutine end",
		t.JobID(), t.TaskGroupID(), t.TaskID())
	switch {
	//当外部取消时，开始写入不是错误
	case ctx.Err() != nil:
		return nil
	//当错误是停止时，也不是错误
	case err == exchange.ErrTerminate:
		return nil
	}

	return
}
