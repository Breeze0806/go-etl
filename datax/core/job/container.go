package job

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"

	"github.com/Breeze0806/go-etl/datax/common/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/schedule"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/common/util"
	"github.com/Breeze0806/go-etl/datax/core"
	statplugin "github.com/Breeze0806/go-etl/datax/core/statistics/container/plugin"
	"github.com/Breeze0806/go-etl/datax/core/taskgroup"
)

type Container struct {
	ctx context.Context
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
	needChannelNumber      int64
	totalStage             int
	errorLimit             util.ErrorRecordChecker
	taskSchduler           *schedule.TaskSchduler
	wg                     sync.WaitGroup
}

func NewContainer(ctx context.Context, conf *config.Json) (c *Container, err error) {
	c = &Container{
		BaseCotainer: core.NewBaseCotainer(),
		ctx:          ctx,
	}
	c.SetConfig(conf)
	c.jobId = c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerJobId, -1)
	if c.jobId < 0 {
		return nil, fmt.Errorf("container job id is invalid")
	}
	return
}

func (c *Container) Start() (err error) {
	log.Infof("DataX jobContainer %v starts job.", c.jobId)
	defer c.destroy()
	c.userConf = c.Config().CloneConfig()

	log.Debugf("DataX jobContainer %v starts to preHandle.", c.jobId)
	if err = c.preHandle(); err != nil {
		return
	}
	log.Infof("DataX jobContainer %v starts to init.", c.jobId)
	if err = c.init(); err != nil {
		return
	}
	log.Infof("DataX jobContainer %v starts to prepare.", c.jobId)
	if err = c.prepare(); err != nil {
		return
	}
	log.Infof("DataX jobContainer %v starts to split.", c.jobId)
	if err = c.split(); err != nil {
		return
	}
	log.Infof("DataX jobContainer %v starts to schedule.", c.jobId)
	if err = c.schedule(); err != nil {
		return
	}
	log.Infof("DataX jobContainer %v starts to post.", c.jobId)
	if err = c.post(); err != nil {
		return
	}
	log.Debugf("DataX jobContainer %v starts to postHandle.", c.jobId)
	if err = c.postHandle(); err != nil {
		return
	}

	return nil
}

func (c *Container) destroy() (err error) {
	if c.jobReader != nil {
		if rerr := c.jobReader.Destroy(c.ctx); rerr != nil {
			log.Errorf("DataX jobContainer %v jobReader %s destroy error: %v",
				c.jobId, c.readerPluginName, rerr)
			err = rerr
		}
	}

	if c.jobWriter != nil {
		if werr := c.jobWriter.Destroy(c.ctx); werr != nil {
			log.Errorf("DataX jobContainer %v jobWriter %s destroy error: %v",
				c.jobId, c.writerPluginName, werr)
			err = werr
		}
	}
	return
}

func (c *Container) init() (err error) {
	c.readerPluginName, err = c.Config().GetString(coreconst.DataxJobContentReaderName)
	if err != nil {
		return
	}

	c.writerPluginName, err = c.Config().GetString(coreconst.DataxJobContentWriterName)
	if err != nil {
		return
	}

	var readerConfig, writerConfig *config.Json
	readerConfig, err = c.Config().GetConfig(coreconst.DataxJobContentReaderParameter)
	if err != nil {
		return
	}

	writerConfig, err = c.Config().GetConfig(coreconst.DataxJobContentWriterParameter)
	if err != nil {
		return
	}

	collector := statplugin.NewDefaultJobCollector(c.Communication())
	c.jobReader, err = c.initReaderJob(collector, readerConfig, writerConfig)
	if err != nil {
		return
	}
	log.Infof("DataX jobContainer %v reader %v inited", c.jobId, c.readerPluginName)
	c.jobWriter, err = c.initWriterJob(collector, readerConfig, writerConfig)
	if err != nil {
		return
	}
	log.Infof("DataX jobContainer %v writer %v inited", c.jobId, c.writerPluginName)
	return
}

func (c *Container) prepare() (err error) {
	if err = c.prepareReaderJob(); err != nil {
		return err
	}
	log.Infof("DataX jobContainer %v reader %v prepared", c.jobId, c.readerPluginName)
	if err = c.prepareWriterJob(); err != nil {
		return err
	}
	log.Infof("DataX jobContainer %v writer %v prepared", c.jobId, c.writerPluginName)
	return
}

func (c *Container) prepareReaderJob() error {
	return c.jobReader.Prepare(c.ctx)
}

func (c *Container) prepareWriterJob() error {
	return c.jobWriter.Prepare(c.ctx)
}

func (c *Container) split() (err error) {
	if err = c.adjustChannelNumber(); err != nil {
		return
	}

	if c.needChannelNumber <= 0 {
		c.needChannelNumber = 1
	}
	var readerConfs, writerConfs, tasksConfigs []*config.Json
	readerConfs, err = c.jobReader.Split(c.ctx, int(c.needChannelNumber))
	if err != nil {
		return
	}

	if len(readerConfs) == 0 {
		err = fmt.Errorf("reader split fail, config is empty")
		return
	}

	taskNumber := len(readerConfs)
	log.Infof("DataX jobContainer %v reader %v split %v tasks", c.jobId, c.readerPluginName, taskNumber)
	writerConfs, err = c.jobWriter.Split(c.ctx, taskNumber)
	if err != nil {
		return
	}

	if len(writerConfs) == 0 {
		err = fmt.Errorf("writer split fail, config is empty")
		return
	}
	log.Infof("DataX jobContainer %v writer %v split %v tasks", c.jobId, c.readerPluginName, len(writerConfs))

	tasksConfigs, err = c.mergeTaskConfigs(readerConfs, writerConfs)
	if err != nil {
		return
	}

	err = c.Config().Set(coreconst.DataxJobContent, tasksConfigs)
	if err != nil {
		return
	}
	c.totalStage = len(tasksConfigs)
	return nil
}

func (c *Container) schedule() (err error) {
	var tasksConfigs []*config.Json
	tasksConfigs, err = c.distributeTaskIntoTaskGroup()
	if err != nil {
		return err
	}

	c.taskSchduler = schedule.NewTaskSchduler(int(c.Config().GetInt64OrDefaullt(
		coreconst.DataxCoreContainerJobMaxWorkerNumber, 4)), len(tasksConfigs))
	defer c.taskSchduler.Stop()
	for i := range tasksConfigs {
		var taskGroup *taskgroup.Container
		taskGroup, err = taskgroup.NewContainer(c.ctx, tasksConfigs[i])
		if err != nil {
			goto End
		}
		c.wg.Add(1)
		var errChan <-chan error
		errChan, err = c.taskSchduler.Push(taskGroup)
		if err != nil {
			c.wg.Done()
			goto End
		}

		go func() {
			defer c.wg.Done()
			select {
			case <-errChan:
			case <-c.ctx.Done():
			}
		}()
	}
End:
	c.wg.Wait()

	return
}

func (c *Container) post() (err error) {
	if err = c.jobReader.Post(c.ctx); err != nil {
		return err
	}
	log.Infof("DataX jobContainer %v reader %v posted", c.jobId, c.readerPluginName)
	if err = c.jobWriter.Post(c.ctx); err != nil {
		return err
	}
	log.Infof("DataX jobContainer %v writer %v posted", c.jobId, c.writerPluginName)
	return
}

func (c *Container) mergeTaskConfigs(readerConfs, writerConfs []*config.Json) (taskConfigs []*config.Json, err error) {
	if len(readerConfs) != len(writerConfs) {
		err = fmt.Errorf("the number of reader tasks are not equal to the number of writer tasks")
		return
	}
	var transformConfs []*config.Json
	transformConfs, err = c.Config().GetConfigArray(coreconst.DataxJobContentTransformer)
	if err != nil {
		return
	}
	log.Infof("DataX jobContainer %v  tansformer config is %v", c.jobId, transformConfs)
	for i := range readerConfs {
		var taskConfig *config.Json
		taskConfig, _ = config.NewJsonFromString("{}")
		err = taskConfig.Set(coreconst.JobReaderName, c.readerPluginName)
		if err != nil {
			return
		}

		err = taskConfig.SetRawString(coreconst.JobReaderParameter, readerConfs[i].String())
		if err != nil {
			return
		}
		err = taskConfig.Set(coreconst.JobWriterName, c.writerPluginName)
		if err != nil {
			return
		}
		err = taskConfig.SetRawString(coreconst.JobWriterParameter, writerConfs[i].String())
		if err != nil {
			return
		}
		if len(transformConfs) != 0 {
			err = taskConfig.Set(coreconst.JobTransformer, transformConfs)
			if err != nil {
				return
			}
		}
		taskConfig.Set(coreconst.TaskId, i)
		taskConfigs = append(taskConfigs, taskConfig)
	}
	return
}

func (c *Container) distributeTaskIntoTaskGroup() (confs []*config.Json, err error) {
	var tasksConfigs []*config.Json
	tasksConfigs, err = c.Config().GetConfigArray(coreconst.DataxJobContent)
	if err != nil {
		return
	}

	channelsPerTaskGroup := c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerTaskgroupChannel, 5)
	channelNumber := c.needChannelNumber
	if channelNumber > int64(len(tasksConfigs)) {
		channelNumber = int64(len(tasksConfigs))
	}
	taskGroupNumber := int(math.Ceil(1.0 * float64(channelNumber) / float64(channelsPerTaskGroup)))
	taskIdMap := parseAndGetResourceMarkAndTaskIdMap(tasksConfigs)
	ss := doAssign(taskIdMap, taskGroupNumber)
	template := c.Config().CloneConfig()
	if err = template.Remove(coreconst.DataxJobContent); err != nil {
		return nil, err
	}

	for i := 0; i < taskGroupNumber; i++ {
		conf := template.CloneConfig()
		if err = conf.Set(coreconst.DataxCoreContainerTaskGroupId, i); err != nil {
			return nil, err
		}
		confs = append(confs, conf)
	}

	for i, v := range ss {
		for j, vj := range v {
			if err = confs[i].Set(coreconst.DataxJobContent+"."+strconv.Itoa(j), tasksConfigs[vj]); err != nil {
				return nil, err
			}
		}
	}
	return
}

func (c *Container) adjustChannelNumber() error {
	var needChannelNumberByByte int64 = math.MaxInt32
	var needChannelNumberByRecord int64 = math.MaxInt32

	if isByteLimit := c.Config().GetInt64OrDefaullt(coreconst.DataxJobSettingSpeedByte, 0) > 0; isByteLimit {
		globalLimitedByteSpeed := c.Config().GetInt64OrDefaullt(coreconst.DataxJobSettingSpeedByte, 10*1024*1024)
		channelLimitedByteSpeed, err := c.Config().GetInt64(coreconst.DataxCoreTransportChannelSpeedByte)
		if err != nil {
			return err
		}
		if channelLimitedByteSpeed <= 0 {
			return fmt.Errorf("%v should be positive", coreconst.DataxCoreTransportChannelSpeedByte)
		}
		needChannelNumberByByte = globalLimitedByteSpeed / channelLimitedByteSpeed
		if needChannelNumberByByte < 1 {
			needChannelNumberByByte = 1
		}
		log.Infof("DataX jobContainer %v set Max-Byte-Speed to %v bytes", c.jobId, globalLimitedByteSpeed)
	}

	if isRecordLimit := c.Config().GetInt64OrDefaullt(coreconst.DataxJobSettingSpeedRecord, 0) > 0; isRecordLimit {
		globalLimitedRecordSpeed := c.Config().GetInt64OrDefaullt(coreconst.DataxJobSettingSpeedRecord, 10*1024*1024)
		channelLimitedRecordSpeed, err := c.Config().GetInt64(coreconst.DataxCoreTransportChannelSpeedRecord)
		if err != nil {
			return err
		}
		if channelLimitedRecordSpeed <= 0 {
			return fmt.Errorf("%v should be positive", coreconst.DataxCoreTransportChannelSpeedByte)
		}

		needChannelNumberByRecord = globalLimitedRecordSpeed / channelLimitedRecordSpeed
		if needChannelNumberByRecord < 1 {
			needChannelNumberByRecord = 1
		}
		log.Infof("DataX jobContainer %v  set Max-Record-Speed to %v records", c.jobId, globalLimitedRecordSpeed)
	}
	if needChannelNumberByByte > needChannelNumberByRecord {
		c.needChannelNumber = needChannelNumberByRecord
	} else {
		c.needChannelNumber = needChannelNumberByByte
	}

	if c.needChannelNumber < math.MaxInt32 {
		return nil
	}

	if isChannelLimit := c.Config().GetInt64OrDefaullt(coreconst.DataxJobSettingSpeedChannel, 0) > 0; isChannelLimit {
		//此时 DataxJobSettingSpeedChannel必然存在
		c.needChannelNumber, _ = c.Config().GetInt64(coreconst.DataxJobSettingSpeedChannel)
		log.Infof("DataX jobContainer %v set Channel-Number to %v channels.", c.jobId, c.needChannelNumber)
		return nil
	}

	return fmt.Errorf("job speed should be setted")
}

func (c *Container) initReaderJob(collector plugin.JobCollector, readerConfig, writerConfig *config.Json) (job reader.Job, err error) {
	ok := false
	job, ok = loader.LoadReaderJob(c.readerPluginName)
	if !ok {
		err = fmt.Errorf("reader %v does not exist", c.readerPluginName)
		return
	}
	job.SetCollector(collector)
	job.SetPluginJobConf(readerConfig)
	job.SetPeerPluginJobConf(writerConfig)
	job.SetPeerPluginName(c.writerPluginName)
	err = job.Init(c.ctx)
	if err != nil {
		return
	}
	return
}

func (c *Container) initWriterJob(collector plugin.JobCollector, readerConfig, writerConfig *config.Json) (job writer.Job, err error) {
	ok := false
	job, ok = loader.LoadWriterJob(c.writerPluginName)
	if !ok {
		err = fmt.Errorf("writer %v does not exist", c.writerPluginName)
		return
	}
	job.SetCollector(collector)
	job.SetPluginJobConf(writerConfig)
	job.SetPeerPluginJobConf(readerConfig)
	job.SetPeerPluginName(c.readerPluginName)
	err = job.Init(c.ctx)
	if err != nil {
		return
	}
	return
}

//preHandle 事实上对于使用者是空壳，reader和writer未实现对应逻辑PreHandle
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
	err = handler.PreHandler(c.ctx, c.Config())
	if err != nil {
		return
	}
	return
}

//postHandle 事实上对于使用者是空壳，reader和writer未实现对应逻辑PostHandle
func (c *Container) postHandle() (err error) {
	if !c.Config().Exists(coreconst.DataxJobPostHandlerPluginType) {
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
	err = handler.PostHandler(c.ctx, c.Config())
	if err != nil {
		return
	}
	return
}

func doAssign(taskIdMap map[string][]int, taskGroupNumber int) [][]int {
	taskGroups := make([][]int, taskGroupNumber)
	var taskMasks []string
	var maxLen int = 0
	for k, v := range taskIdMap {
		if maxLen < len(v) {
			maxLen = len(v)
		}
		taskMasks = append(taskMasks, k)
	}

	sort.Sort(sort.StringSlice(taskMasks))

	index := 0
	for i := 0; i < maxLen; i++ {
		for _, v := range taskMasks {
			if len(taskIdMap[v]) > 0 {
				last := 0
				last, taskIdMap[v] = taskIdMap[v][0], taskIdMap[v][1:]
				taskGroups[index%taskGroupNumber] = append(taskGroups[index%taskGroupNumber], last)
				index++
			}
		}
	}
	return taskGroups
}

func parseAndGetResourceMarkAndTaskIdMap(tasksConfigs []*config.Json) map[string][]int {
	writerMap := make(map[string][]int)
	readerMap := make(map[string][]int)
	for i, v := range tasksConfigs {
		key := v.GetStringOrDefaullt(coreconst.JobReaderParameterLoadBalanceResourceMark, "aFakeResourceMarkForLoadBalance")
		readerMap[key] = append(readerMap[key], i)
		key = v.GetStringOrDefaullt(coreconst.JobWriterParameterLoadBalanceResourceMark, "aFakeResourceMarkForLoadBalance")
		writerMap[key] = append(writerMap[key], i)
	}
	if len(readerMap) > len(writerMap) {
		return readerMap
	}
	return writerMap
}
