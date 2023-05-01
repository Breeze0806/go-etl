// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package job

import (
	"context"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/common/util"
	"github.com/Breeze0806/go-etl/datax/core"
	statplugin "github.com/Breeze0806/go-etl/datax/core/statistics/container/plugin"
	"github.com/Breeze0806/go-etl/datax/core/taskgroup"
	"github.com/Breeze0806/go-etl/schedule"
	"github.com/pingcap/errors"
)

//Container 工作容器环境，所有的工作都在本容器环境中执行
type Container struct {
	ctx context.Context
	*core.BaseCotainer
	jobID                  int64
	readerPluginName       string
	writerPluginName       string
	jobReader              reader.Job
	jobWriter              writer.Job
	userConf               *config.JSON
	startTimestamp         int64
	endTimestamp           int64
	startTransferTimeStamp int64
	endTransferTimeStamp   int64
	needChannelNumber      int64
	totalStage             int
	//todo ErrorRecordChecker未使用
	errorLimit   util.ErrorRecordChecker
	taskSchduler *schedule.TaskSchduler
	wg           sync.WaitGroup
}

//NewContainer 通过上下文ctx和JSON配置conf生成工作容器环境
//当container job id小于0时，会报错
func NewContainer(ctx context.Context, conf *config.JSON) (c *Container, err error) {
	c = &Container{
		BaseCotainer: core.NewBaseCotainer(),
		ctx:          ctx,
	}
	c.SetConfig(conf)
	c.jobID = c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerJobID, -1)
	if c.jobID < 0 {
		return nil, errors.New("container job id is invalid")
	}
	return
}

//Start 工作容器开始工作
func (c *Container) Start() (err error) {
	log.Infof("DataX jobContainer %v starts job.", c.jobID)
	defer c.destroy()
	c.userConf = c.Config().CloneConfig()

	log.Debugf("DataX jobContainer %v starts to preHandle.", c.jobID)
	if err = c.preHandle(); err != nil {
		return
	}

	log.Infof("DataX jobContainer %v starts to init.", c.jobID)
	if err = c.init(); err != nil {
		return
	}
	log.Infof("DataX jobContainer %v starts to prepare.", c.jobID)
	if err = c.prepare(); err != nil {
		return
	}
	log.Infof("DataX jobContainer %v starts to split.", c.jobID)
	if err = c.split(); err != nil {
		return
	}
	log.Infof("DataX jobContainer %v starts to schedule.", c.jobID)
	if err = c.schedule(); err != nil {
		return
	}
	log.Infof("DataX jobContainer %v starts to post.", c.jobID)
	if err = c.post(); err != nil {
		return
	}
	log.Debugf("DataX jobContainer %v starts to postHandle.", c.jobID)
	if err = c.postHandle(); err != nil {
		return
	}

	return nil
}

//destroy 销毁，在jobReader不为空时进行销毁
//在jobWriter不为空时进行销毁
func (c *Container) destroy() (err error) {
	if c.jobReader != nil {
		if rerr := c.jobReader.Destroy(c.ctx); rerr != nil {
			log.Errorf("DataX jobContainer %v jobReader %s destroy error: %v",
				c.jobID, c.readerPluginName, rerr)
			err = rerr
		}
	}

	if c.jobWriter != nil {
		if werr := c.jobWriter.Destroy(c.ctx); werr != nil {
			log.Errorf("DataX jobContainer %v jobWriter %s destroy error: %v",
				c.jobID, c.writerPluginName, werr)
			err = werr
		}
	}
	return
}

//init 检查并初始化读取器和写入器工作
//当配置文件读取器和写入器的名字和参数不存在的情况下会报错
//另外，读取器和写入器工作初始化失败也会导致报错
func (c *Container) init() (err error) {
	c.readerPluginName, err = c.Config().GetString(coreconst.DataxJobContentReaderName)
	if err != nil {
		return
	}

	c.writerPluginName, err = c.Config().GetString(coreconst.DataxJobContentWriterName)
	if err != nil {
		return
	}

	var readerConfig *config.JSON
	readerConfig, err = c.Config().GetConfig(coreconst.DataxJobContentReaderParameter)
	if err != nil {
		return
	}

	var writerConfig *config.JSON
	writerConfig, err = c.Config().GetConfig(coreconst.DataxJobContentWriterParameter)
	if err != nil {
		return
	}

	var jobSettingConf *config.JSON
	if jobSettingConf, err = c.Config().GetConfig(coreconst.DataxJobSetting); err != nil {
		jobSettingConf, _ = config.NewJSONFromString("{}")
		err = nil
	}

	readerConfig.Set(coreconst.DataxJobSetting, jobSettingConf)

	writerConfig.Set(coreconst.DataxJobSetting, jobSettingConf)

	collector := statplugin.NewDefaultJobCollector(c.Communication())
	c.jobReader, err = c.initReaderJob(collector, readerConfig, writerConfig)
	if err != nil {
		return
	}
	log.Infof("DataX jobContainer %v reader %v inited", c.jobID, c.readerPluginName)
	c.jobWriter, err = c.initWriterJob(collector, readerConfig, writerConfig)
	if err != nil {
		return
	}
	log.Infof("DataX jobContainer %v writer %v inited", c.jobID, c.writerPluginName)
	return
}

//prepare 准备读取器和写入器工作
//如果读取器和写入器工作准备失败就会报错
func (c *Container) prepare() (err error) {
	if err = c.prepareReaderJob(); err != nil {
		return err
	}
	log.Infof("DataX jobContainer %v reader %v prepared", c.jobID, c.readerPluginName)
	if err = c.prepareWriterJob(); err != nil {
		return err
	}
	log.Infof("DataX jobContainer %v writer %v prepared", c.jobID, c.writerPluginName)
	return
}

//prepareReaderJob 准备读取工作
func (c *Container) prepareReaderJob() error {
	return c.jobReader.Prepare(c.ctx)
}

//prepareReaderJob 准备写入工作
func (c *Container) prepareWriterJob() error {
	return c.jobWriter.Prepare(c.ctx)
}

//split 切分读取器和写入器工作
//先进行读取工作切分成多个任务，再根据读取工作切分的结果进行写入工作切分多个任务
//然后逐个将单个读取任务、单个写入任务和转化器组合成完整任务组，由于reader，writer，channel模型
//切分时读取器和写入器的比例为1:1，所以这里可以将reader和writer的配置整合到一起
func (c *Container) split() (err error) {
	if err = c.adjustChannelNumber(); err != nil {
		return
	}

	if c.needChannelNumber <= 0 {
		c.needChannelNumber = 1
	}
	var readerConfs, writerConfs, tasksConfigs []*config.JSON
	readerConfs, err = c.jobReader.Split(c.ctx, int(c.needChannelNumber))
	if err != nil {
		return
	}

	if len(readerConfs) == 0 {
		err = errors.New("reader split fail, config is empty")
		return
	}

	taskNumber := len(readerConfs)
	log.Infof("DataX jobContainer %v reader %v split %v tasks", c.jobID, c.readerPluginName, taskNumber)
	writerConfs, err = c.jobWriter.Split(c.ctx, taskNumber)
	if err != nil {
		return
	}

	if len(writerConfs) == 0 {
		err = errors.New("writer split fail, config is empty")
		return
	}
	log.Infof("DataX jobContainer %v writer %v split %v tasks", c.jobID, c.readerPluginName, len(writerConfs))

	tasksConfigs, err = c.mergeTaskConfigs(readerConfs, writerConfs)
	if err != nil {
		return
	}

	c.Config().Set(coreconst.DataxJobContent, tasksConfigs)

	c.totalStage = len(tasksConfigs)
	return nil
}

//schedule 使用调度器将任务组进行调度，进入执行队列中
func (c *Container) schedule() (err error) {
	var tasksConfigs []*config.JSON
	tasksConfigs, err = c.distributeTaskIntoTaskGroup()
	if err != nil {
		return err
	}

	c.taskSchduler = schedule.NewTaskSchduler(int(c.Config().GetInt64OrDefaullt(
		coreconst.DataxCoreContainerJobMaxWorkerNumber, 4)), len(tasksConfigs))
	defer c.taskSchduler.Stop()
	var taskGroups []*taskgroup.Container
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
		taskGroups = append(taskGroups, taskGroup)
		go func(taskGroup *taskgroup.Container) {
			defer c.wg.Done()
			select {
			case taskGroup.Err = <-errChan:
			case <-c.ctx.Done():
			}
		}(taskGroup)
	}
End:
	c.wg.Wait()

	b := &strings.Builder{}
	for _, t := range taskGroups {
		if t.Err != nil {
			b.WriteString(t.Err.Error())
			b.WriteByte(' ')
		}
	}
	if b.Len() != 0 {
		return errors.NewNoStackError(b.String())
	}
	return
}

//post 后置通知
func (c *Container) post() (err error) {
	if err = c.jobReader.Post(c.ctx); err != nil {
		return err
	}
	log.Infof("DataX jobContainer %v reader %v posted", c.jobID, c.readerPluginName)
	if err = c.jobWriter.Post(c.ctx); err != nil {
		return err
	}
	log.Infof("DataX jobContainer %v writer %v posted", c.jobID, c.writerPluginName)
	return
}

//mergeTaskConfigs 逐个将单个读取任务、单个写入任务和转化器组合成完整任务组
func (c *Container) mergeTaskConfigs(readerConfs, writerConfs []*config.JSON) (taskConfigs []*config.JSON, err error) {
	if len(readerConfs) != len(writerConfs) {
		err = errors.New("the number of reader tasks are not equal to the number of writer tasks")
		return
	}
	var transformConfs []*config.JSON
	transformConfs, err = c.Config().GetConfigArray(coreconst.DataxJobContentTransformer)
	if err != nil {
		return
	}
	log.Infof("DataX jobContainer %v tansformer config is %v", c.jobID, transformConfs)
	for i := range readerConfs {
		var taskConfig *config.JSON
		taskConfig, _ = config.NewJSONFromString("{}")
		taskConfig.Set(coreconst.JobReaderName, c.readerPluginName)

		taskConfig.SetRawString(coreconst.JobReaderParameter, readerConfs[i].String())

		taskConfig.Set(coreconst.JobWriterName, c.writerPluginName)

		taskConfig.SetRawString(coreconst.JobWriterParameter, writerConfs[i].String())

		if len(transformConfs) != 0 {
			taskConfig.Set(coreconst.JobTransformer, transformConfs)
		}
		taskConfig.Set(coreconst.TaskID, i)
		taskConfigs = append(taskConfigs, taskConfig)
	}
	return
}

//distributeTaskIntoTaskGroup 公平的分配 task 到对应的 taskGroup 中。
//公平体现在：会考虑 task 中对资源负载作的 load 标识进行更均衡的作业分配操作。
func (c *Container) distributeTaskIntoTaskGroup() (confs []*config.JSON, err error) {
	var tasksConfigs []*config.JSON
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
	taskIDMap := parseAndGetResourceMarkAndTaskIDMap(tasksConfigs)
	ss := doAssign(taskIDMap, taskGroupNumber)
	template := c.Config().CloneConfig()
	template.Remove(coreconst.DataxJobContent)

	for i := 0; i < taskGroupNumber; i++ {
		conf := template.CloneConfig()
		conf.Set(coreconst.DataxCoreContainerTaskGroupID, i)
		confs = append(confs, conf)
	}

	for i, v := range ss {
		for j, vj := range v {
			confs[i].Set(coreconst.DataxJobContent+"."+strconv.Itoa(j), tasksConfigs[vj])
		}
	}
	return
}

//adjustChannelNumber 自适应化通道数量
//依次根据字节流大小，记录数大小以及通道数大小生成通道数量
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
			return errors.Errorf("%v should be positive", coreconst.DataxCoreTransportChannelSpeedByte)
		}
		needChannelNumberByByte = globalLimitedByteSpeed / channelLimitedByteSpeed
		if needChannelNumberByByte < 1 {
			needChannelNumberByByte = 1
		}
		log.Infof("DataX jobContainer %v set Max-Byte-Speed to %v bytes", c.jobID, globalLimitedByteSpeed)
	}

	if isRecordLimit := c.Config().GetInt64OrDefaullt(coreconst.DataxJobSettingSpeedRecord, 0) > 0; isRecordLimit {
		globalLimitedRecordSpeed := c.Config().GetInt64OrDefaullt(coreconst.DataxJobSettingSpeedRecord, 10*1024*1024)
		channelLimitedRecordSpeed, err := c.Config().GetInt64(coreconst.DataxCoreTransportChannelSpeedRecord)
		if err != nil {
			return err
		}
		if channelLimitedRecordSpeed <= 0 {
			return errors.Errorf("%v should be positive", coreconst.DataxCoreTransportChannelSpeedByte)
		}

		needChannelNumberByRecord = globalLimitedRecordSpeed / channelLimitedRecordSpeed
		if needChannelNumberByRecord < 1 {
			needChannelNumberByRecord = 1
		}
		log.Infof("DataX jobContainer %v  set Max-Record-Speed to %v records", c.jobID, globalLimitedRecordSpeed)
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
		log.Infof("DataX jobContainer %v set Channel-Number to %v channels.", c.jobID, c.needChannelNumber)
		return nil
	}

	return errors.New("job speed should be setted")
}

//initReaderJob 初始化读取工作
//当读取插件名找不到读取工作或者初始化失败就会报错
func (c *Container) initReaderJob(collector plugin.JobCollector, readerConfig, writerConfig *config.JSON) (job reader.Job, err error) {
	ok := false
	job, ok = loader.LoadReaderJob(c.readerPluginName)
	if !ok {
		err = errors.Errorf("reader %v does not exist", c.readerPluginName)
		return
	}
	job.SetCollector(collector)
	job.SetPluginJobConf(readerConfig)
	job.SetPeerPluginJobConf(writerConfig)
	job.SetPeerPluginName(c.writerPluginName)
	job.SetJobID(c.jobID)
	err = job.Init(c.ctx)
	if err != nil {
		return
	}
	return
}

//initReaderJob 初始化写入工作
//当写入插件名找不到写入工作或者初始化失败就会报错
func (c *Container) initWriterJob(collector plugin.JobCollector, readerConfig, writerConfig *config.JSON) (job writer.Job, err error) {
	ok := false
	job, ok = loader.LoadWriterJob(c.writerPluginName)
	if !ok {
		err = errors.Errorf("writer %v does not exist", c.writerPluginName)
		return
	}
	job.SetCollector(collector)
	job.SetPluginJobConf(writerConfig)
	job.SetPeerPluginJobConf(readerConfig)
	job.SetPeerPluginName(c.readerPluginName)
	job.SetJobID(c.jobID)
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
		return errors.Errorf("%v is not valid plugin Type", handlerPluginTypeStr)
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
		return errors.Errorf("%v is not valid plugin Type", handlerPluginTypeStr)
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

// doAssign 平均分配
//   需要实现的效果通过例子来说是：
//   a 库上有表：0, 1, 2
//   a 库上有表：3, 4
//   c 库上有表：5, 6, 7

//   如果有 4个 taskGroup
//   则 assign 后的结果为：
//   taskGroup-0: 0,  4,
//   taskGroup-1: 3,  6,
//   taskGroup-2: 5,  2,
//   taskGroup-3: 1,  7
func doAssign(taskIDMap map[string][]int, taskGroupNumber int) [][]int {
	taskGroups := make([][]int, taskGroupNumber)
	var taskMasks []string
	var maxLen int = 0
	for k, v := range taskIDMap {
		if maxLen < len(v) {
			maxLen = len(v)
		}
		taskMasks = append(taskMasks, k)
	}

	sort.Sort(sort.StringSlice(taskMasks))

	index := 0
	for i := 0; i < maxLen; i++ {
		for _, v := range taskMasks {
			if len(taskIDMap[v]) > 0 {
				last := 0
				last, taskIDMap[v] = taskIDMap[v][0], taskIDMap[v][1:]
				taskGroups[index%taskGroupNumber] = append(taskGroups[index%taskGroupNumber], last)
				index++
			}
		}
	}
	return taskGroups
}

//parseAndGetResourceMarkAndTaskIDMap 根据task 配置，获取到：
//资源名称 --> taskId(List) 的 map 映射关系(对资源负载作的 load 标识: 任务编号)
func parseAndGetResourceMarkAndTaskIDMap(tasksConfigs []*config.JSON) map[string][]int {
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
