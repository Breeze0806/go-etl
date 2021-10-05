package writer

import (
	"context"
	"sync"
	"time"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/core/transport/exchange"
	"github.com/Breeze0806/go-etl/element"
)

//BatchWriter 批量写入器
type BatchWriter interface {
	JobID() int64                                                   //工作编号
	TaskGroupID() int64                                             //任务组编号
	TaskID() int64                                                  //任务编号
	BatchSize() int                                                 //单次批量写入数
	BatchTimeout() time.Duration                                    //单次批量写入超时时间
	BatchWrite(ctx context.Context, records []element.Record) error //批量写入
}

//StartWrite 通过批量写入器writer和记录接受器receiver将记录写入数据库
func StartWrite(ctx context.Context, w BatchWriter,
	receiver plugin.RecordReceiver) (err error) {
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
				w.JobID(), w.TaskGroupID(), w.TaskID())
		}()
		log.Debugf("jobID: %v taskgroupID:%v taskID: %v start to get records",
			w.JobID(), w.TaskGroupID(), w.TaskID())
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
	ticker := time.NewTicker(w.BatchTimeout())
	defer ticker.Stop()
	var records []element.Record
	log.Debugf("jobID: %v taskgroupID:%v taskID: %v  start to BatchExec",
		w.JobID(), w.TaskGroupID(), w.TaskID())
	for {
		select {
		case record, ok := <-recordChan:
			if !ok {
				//当写入结束时，将剩余的记录写入数据库
				if err = w.BatchWrite(ctx, records); err != nil {
					log.Errorf("jobID: %v taskgroupID:%v taskID: %v BatchExec(%v) error: %v",
						w.JobID(), w.TaskGroupID(), w.TaskID(), records, err)
				}
				records = nil
				err = rerr
				goto End
			}
			records = append(records, record)

			//当数据量超过单次批量数时 写入数据库
			if len(records) >= w.BatchSize() {
				if err = w.BatchWrite(ctx, records); err != nil {
					log.Errorf("jobID: %v taskgroupID:%v taskID: %v BatchExec(%v) error: %v",
						w.JobID(), w.TaskGroupID(), w.TaskID(), records, err)
					goto End
				}
				records = nil
			}
		//当写入数据未达到单次批量数，超时也写入
		case <-ticker.C:
			if len(records) > 0 {
				if err = w.BatchWrite(ctx, records); err != nil {
					log.Errorf("jobID: %v taskgroupID:%v taskID: %v BatchWrite(%v) error: %v",
						w.JobID(), w.TaskGroupID(), w.TaskID(), records, err)
					goto End
				}
			}
			records = nil
		}
	}
End:
	cancel()
	log.Debugf("jobID: %v taskgroupID:%v taskID: %v wait all goroutine",
		w.JobID(), w.TaskGroupID(), w.TaskID())
	//等待携程结束
	wg.Wait()
	log.Debugf("jobID: %v taskgroupID:%v taskID: %v wait all goroutine end",
		w.JobID(), w.TaskGroupID(), w.TaskID())
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
