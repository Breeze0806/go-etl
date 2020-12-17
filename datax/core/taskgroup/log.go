package taskgroup

import etllog "github.com/Breeze0806/go-etl/log"

var log etllog.Logger

func LogInit() {
	log = etllog.GetLogger()
}
