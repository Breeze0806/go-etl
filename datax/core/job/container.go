package job

import (
	"fmt"

	"github.com/Breeze0806/go-etl/datax/common/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/common/util"
	"github.com/Breeze0806/go-etl/datax/core"
)

type Container struct {
	*core.BaseCotainer
	jobId                  int64
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

func (c *Container) init() (err error) {
	c.jobId = c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerJobId, -1)
	if c.jobId < 0 {
		return fmt.Errorf("container job id is invalid")
	}

	return
}

func (c *Container) initReaderJob() (err error) {
	return
}

func (c *Container) initWriterJob() (err error) {
	return
}

func (c *Container) preHandle() (err error) {
	if !c.Config().Exists(coreconst.DataxJobPreHandlerPluginType) {
		return
	}
	handlerPluginTypeStr := ""
	handlerPluginTypeStr, err = c.Config().GetString(coreconst.DataxJobPreHandlerPluginType)
	if err != nil {
		return
	}
	handlerPluginType := plugin.Type(handlerPluginTypeStr)
	if !handlerPluginType.IsValid() {
		return fmt.Errorf("%v is not valid plugin Type", handlerPluginTypeStr)
	}
	handlerPluginName := ""
	handlerPluginName, err = c.Config().GetString(coreconst.DataxJobPreHandlerPluginName)
	if err != nil {
		return
	}
	var handler plugin.Job
	handler, err = loader.LoadJobPlugin(handlerPluginType, handlerPluginName)
	if err != nil {
		return
	}
	err = handler.PreHandler(c.Config())
	if err != nil {
		return
	}
	return
}

func (c *Container) postHandle() (err error) {
	if !c.Config().Exists(coreconst.DataxJobPreHandlerPluginType) {
		return
	}
	handlerPluginTypeStr := ""
	handlerPluginTypeStr, err = c.Config().GetString(coreconst.DataxJobPostHandlerPluginType)
	if err != nil {
		return
	}
	handlerPluginType := plugin.Type(handlerPluginTypeStr)
	if !handlerPluginType.IsValid() {
		return fmt.Errorf("%v is not valid plugin Type", handlerPluginTypeStr)
	}
	handlerPluginName := ""
	handlerPluginName, err = c.Config().GetString(coreconst.DataxJobPostHandlerPluginName)
	if err != nil {
		return
	}
	var handler plugin.Job
	handler, err = loader.LoadJobPlugin(handlerPluginType, handlerPluginName)
	if err != nil {
		return
	}
	err = handler.PostHandler(c.Config())
	if err != nil {
		return
	}
	return
}
