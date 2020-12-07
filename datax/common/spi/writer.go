package spi

import "github.com/Breeze0806/go-etl/datax/common/spi/writer"

type Writer interface {
	Job() writer.Job
	Task() writer.Task
}
