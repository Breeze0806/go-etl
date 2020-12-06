package job

import (
	"github.com/Breeze0806/go-etl/datax/common/config"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/common/util"
	"github.com/Breeze0806/go-etl/datax/core"
)

type Container struct {
	*core.BaseCotainer
	jobID                  int64
	readerPluginName       string
	writerPluginName       string
	jobReader              reader.Job
	jobWriter              writer.Job
	userConf               *config.Json
	startTimestamp         int64
	endTimestamp           int64
	startTransferTimeStamp int64
	endTransferTimeStamp   int64
	needChannelNumber      int
	totalStage             int
	errorLimit             util.ErrorRecordChecker
}

func (c *Container) Start() error {
	log.Infof("DataX jobContainer starts job.")
	defer c.destroy()
	c.userConf = c.Config().Clone()

	return nil
}

func (c *Container) destroy() (err error) {
	if c.jobReader != nil {
		if err = c.jobReader.Destroy(); err != nil {
			log.Errorf("jobReader %s destroy error: %v", c.readerPluginName, err)
		}
	}

	if c.jobWriter != nil {
		if err = c.jobWriter.Destroy(); err != nil {
			log.Errorf("jobWriter %s destroy error: %v", c.writerPluginName, err)
		}
	}
	return
}

func (c *Container) preHandle() (err error) {
	return
}
