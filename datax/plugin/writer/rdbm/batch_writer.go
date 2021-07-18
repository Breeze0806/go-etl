package rdbm

import (
	"context"
	"sync"
	"time"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/core/transport/exchange"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

//BatchWriter 批量写入器
type BatchWriter interface {
	JobID() int64                         //工作编号
	TaskGroupID() int64                   //任务组编号
	TaskID() int64                        //任务编号
	BatchSize() int                       //单次批量写入数
	BatchTimeout() time.Duration          //单次批量写入超时时间
	BatchWrite(ctx context.Context) error //批量写入
	Options() *database.ParameterOptions  //数据库选项
}

//StartWrite 通过批量写入器writer和记录接受器receiver将记录写入数据库
func StartWrite(ctx context.Context, writer BatchWriter,
	receiver plugin.RecordReceiver) (err error) {
	opts := writer.Options()
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
				writer.JobID(), writer.TaskGroupID(), writer.TaskID())
		}()
		log.Debugf("jobID: %v taskgroupID:%v taskID: %v start to get records",
			writer.JobID(), writer.TaskGroupID(), writer.TaskID())
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
	ticker := time.NewTicker(writer.BatchTimeout())
	defer ticker.Stop()
	var records []element.Record
	log.Debugf("jobID: %v taskgroupID:%v taskID: %v  start to BatchExec",
		writer.JobID(), writer.TaskGroupID(), writer.TaskID())
	for {
		select {
		case record, ok := <-recordChan:
			if !ok {
				//当写入结束时，将剩余的记录写入数据库
				opts.Records = records
				if err = writer.BatchWrite(ctx); err != nil {
					log.Errorf("jobID: %v taskgroupID:%v taskID: %v BatchExec error: %v",
						writer.JobID(), writer.TaskGroupID(), writer.TaskID(), err)
				}
				opts.Records = nil
				err = rerr
				goto End
			}
			records = append(records, record)

			//当数据量超过单次批量数时 写入数据库
			if len(records) >= writer.BatchSize() {
				opts.Records = records
				if err = writer.BatchWrite(ctx); err != nil {
					log.Errorf("jobID: %v taskgroupID:%v taskID: %v BatchExec error: %v",
						writer.JobID(), writer.TaskGroupID(), writer.TaskID(), err)
					goto End
				}
				opts.Records = nil
				records = nil
			}
		//当写入数据未达到单次批量数，超时也写入
		case <-ticker.C:
			opts.Records = records
			if err = writer.BatchWrite(ctx); err != nil {
				log.Errorf("jobID: %v taskgroupID:%v taskID: %v BatchExec error: %v",
					writer.JobID(), writer.TaskGroupID(), writer.TaskID(), err)
				goto End
			}
			opts.Records = nil
			records = nil
		}
	}
End:
	cancel()
	log.Debugf("jobID: %v taskgroupID:%v taskID: %v wait all goroutine",
		writer.JobID(), writer.TaskGroupID(), writer.TaskID())
	//等待携程结束
	wg.Wait()
	log.Debugf("jobID: %v taskgroupID:%v taskID: %v wait all goroutine end",
		writer.JobID(), writer.TaskGroupID(), writer.TaskID())
	switch {
	//当外部取消时，开始写入不是错误
	case ctx.Err() != nil:
		return nil
	//当错误是停止是，也不是错误
	case err == exchange.ErrTerminate:
		return nil
	}
	return

}
