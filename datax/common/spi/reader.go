package spi

import "github.com/Breeze0806/go-etl/datax/common/spi/reader"

type Reader interface {
	Job() reader.Job
	Task() reader.Task
}
