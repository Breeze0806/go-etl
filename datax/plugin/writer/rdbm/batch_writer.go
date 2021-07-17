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

type BatchWriter interface {
	JobID() int64
	TaskGroupID() int64
	TaskID() int64
	BatchSize() int
	BatchTimeout() time.Duration
	BatchWrite(ctx context.Context) error
	Options() *database.ParameterOptions
}

func StartWrite(ctx context.Context, writer BatchWriter,
	receiver plugin.RecordReceiver) (err error) {
	opts := writer.Options()
	recordChan := make(chan element.Record)
	var rerr error
	afterCtx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
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

			if rerr != exchange.ErrEmpty {
				select {
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
	wg.Wait()
	log.Debugf("jobID: %v taskgroupID:%v taskID: %v wait all goroutine end",
		writer.JobID(), writer.TaskGroupID(), writer.TaskID())
	switch {
	case ctx.Err() != nil:
		return nil
	case err == exchange.ErrTerminate:
		return nil
	}
	return

}
