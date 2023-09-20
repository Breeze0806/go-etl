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

package dbms

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	dbmsreader "github.com/Breeze0806/go-etl/datax/plugin/reader/dbms"
	"github.com/Breeze0806/go-etl/schedule"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/Breeze0806/go/time2"
)

// 默认参数
var (
	defalutBatchSize    = 1000
	defalutBatchTimeout = 1 * time.Second
)

// Config 关系数据库写入器配置
type Config interface {
	GetUsername() string                                                     //获取用户名
	GetPassword() string                                                     //获取密码
	GetURL() string                                                          //获取连接url
	GetColumns() []dbmsreader.Column                                         //获取列信息
	GetBaseTable() *database.BaseTable                                       //获取表信息
	GetWriteMode() string                                                    //获取写入模式
	GetBatchSize() int                                                       //单次批量写入数
	GetBatchTimeout() time.Duration                                          //单次批量写入超时时间
	GetRetryStrategy(j schedule.RetryJudger) (schedule.RetryStrategy, error) //获取重试策略
	IgnoreOneByOneError() bool                                               //忽略一个个重试的错误
	GetPreSQL() []string                                                     //获取准备的SQL语句
	GetPostSQL() []string                                                    //获取结束的SQL语句
}

// BaseConfig 用于实现基本的关系数据库配置，如无特殊情况采用该配置，帮助快速实现writer
type BaseConfig struct {
	Username            string                `json:"username"`     //用户名
	Password            string                `json:"password"`     //密码
	Column              []string              `json:"column"`       //列信息
	Connection          dbmsreader.ConnConfig `json:"connection"`   //连接信息
	WriteMode           string                `json:"writeMode"`    //写入模式,如插入insert
	BatchSize           int                   `json:"batchSize"`    //单次批量写入数
	BatchTimeout        time2.Duration        `json:"batchTimeout"` //单次批量写入超时时间
	PreSQL              []string              `json:"preSQL"`       //准备的SQL语句
	PostSQL             []string              `json:"postSQL"`      //结束的SQL语句
	ignoreOneByOneError bool                  //忽略一个个重试的错误
	newRetryStrategy    func(j schedule.RetryJudger) (schedule.RetryStrategy, error)
}

// NewBaseConfig 从conf解析出关系数据库配置
func NewBaseConfig(conf *config.JSON) (c *BaseConfig, err error) {
	c = &BaseConfig{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	var jobsetting *config.JSON
	jobsetting, err = conf.GetConfig(coreconst.DataxJobSetting)
	if err != nil {
		jobsetting, err = config.NewJSONFromString("{}")
	}
	c.ignoreOneByOneError, _ = jobsetting.GetBool("retry.ignoreOneByOneError")

	c.newRetryStrategy = func(j schedule.RetryJudger) (schedule.RetryStrategy, error) {
		return schedule.NewRetryStrategy(j, jobsetting)
	}

	if err = checkHasSelect(c.PreSQL); err != nil {
		return nil, fmt.Errorf("check preSQL fail. error: %v", err)
	}

	if err = checkHasSelect(c.PostSQL); err != nil {
		return nil, fmt.Errorf("check postSQL fail. error: %v", err)
	}
	return
}

// GetUsername 获取用户名
func (b *BaseConfig) GetUsername() string {
	return b.Username
}

// GetPassword 获取密码
func (b *BaseConfig) GetPassword() string {
	return b.Password
}

// GetURL 获取连接url
func (b *BaseConfig) GetURL() string {
	return b.Connection.URL
}

// GetColumns 获取列信息
func (b *BaseConfig) GetColumns() (columns []dbmsreader.Column) {
	for _, v := range b.Column {
		columns = append(columns, &dbmsreader.BaseColumn{
			Name: v,
		})
	}
	return
}

// GetBaseTable 获取表信息
func (b *BaseConfig) GetBaseTable() *database.BaseTable {
	return database.NewBaseTable(b.Connection.Table.Db, b.Connection.Table.Schema,
		b.Connection.Table.Name)
}

// GetWriteMode 获取写入模式
func (b *BaseConfig) GetWriteMode() string {
	return b.WriteMode
}

// GetBatchTimeout 单次批量超时时间
func (b *BaseConfig) GetBatchTimeout() time.Duration {
	if b.BatchTimeout.Duration == 0 {
		return defalutBatchTimeout
	}
	return b.BatchTimeout.Duration
}

// GetBatchSize 单次批量写入数
func (b *BaseConfig) GetBatchSize() int {
	if b.BatchSize == 0 {
		return defalutBatchSize
	}

	return b.BatchSize
}

// GetPreSQL 获取准备的SQL语句
func (b *BaseConfig) GetPreSQL() []string {
	return getSQlsWithoutEmpty(b.PreSQL)
}

// GetPostSQL 获取结束的SQL语句
func (b *BaseConfig) GetPostSQL() []string {
	return getSQlsWithoutEmpty(b.PostSQL)
}

// IgnoreOneByOneError 忽略一个个重试的错误
func (b *BaseConfig) IgnoreOneByOneError() bool {
	return b.ignoreOneByOneError
}

// GetRetryStrategy 获取重试策略
func (b *BaseConfig) GetRetryStrategy(j schedule.RetryJudger) (schedule.RetryStrategy,
	error) {
	return b.newRetryStrategy(j)
}

func getSQlsWithoutEmpty(sqls []string) (res []string) {
	for _, v := range sqls {
		if v != "" {
			res = append(res, v)
		}
	}
	return res
}

func checkHasSelect(sqls []string) (err error) {
	for i, v := range sqls {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(v)), "select") {
			err = fmt.Errorf("%vst sql(%v) has select", i, v)
			return
		}
	}
	return
}
