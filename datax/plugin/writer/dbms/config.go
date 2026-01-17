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

// Default Parameters
var (
	defalutBatchSize    = 1000
	defalutBatchTimeout = 1 * time.Second
)

// Config - Relational Database Writer Configuration
type Config interface {
	GetUsername() string                                                     // Get Username
	GetPassword() string                                                     // Get Password
	GetURL() string                                                          // Get Connection URL
	GetColumns() []dbmsreader.Column                                         // Get Column Information
	GetBaseTable() *database.BaseTable                                       // Get Table Information
	GetWriteMode() string                                                    // Get Write Mode
	GetBatchSize() int                                                       // Batch Size for Single Write
	GetBatchTimeout() time.Duration                                          // Batch Timeout for Single Write
	GetRetryStrategy(j schedule.RetryJudger) (schedule.RetryStrategy, error) // Get Retry Strategy
	IgnoreOneByOneError() bool                                               // Ignore Individual Retry Errors
	GetPreSQL() []string                                                     // Get Prepared SQL Statement
	GetPostSQL() []string                                                    // Get Ending SQL Statement
}

// BaseConfig - Basic Relational Database Configuration for writers. Unless there are special requirements, this configuration can be used to quickly implement writers.
type BaseConfig struct {
	Username            string                `json:"username"`     // Username
	Password            string                `json:"password"`     // Password
	Column              []string              `json:"column"`       // Column Information
	Connection          dbmsreader.ConnConfig `json:"connection"`   // Connection Information
	WriteMode           string                `json:"writeMode"`    // Write Mode, e.g., Insert
	BatchSize           int                   `json:"batchSize"`    // Batch Size for Single Write
	BatchTimeout        time2.Duration        `json:"batchTimeout"` // Batch Timeout for Single Write
	PreSQL              []string              `json:"preSQL"`       // Prepared SQL Statement
	PostSQL             []string              `json:"postSQL"`      // Ending SQL Statement
	ignoreOneByOneError bool                  // Ignore Individual Retry Errors
	newRetryStrategy    func(j schedule.RetryJudger) (schedule.RetryStrategy, error)
}

// NewBaseConfig - Extract relational database configuration from the configuration file.
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

// GetUsername - Retrieve the username.
func (b *BaseConfig) GetUsername() string {
	return b.Username
}

// GetPassword - Retrieve the password.
func (b *BaseConfig) GetPassword() string {
	return b.Password
}

// GetURL - Retrieve the connection URL.
func (b *BaseConfig) GetURL() string {
	return b.Connection.URL
}

// GetColumns - Retrieve column information.
func (b *BaseConfig) GetColumns() (columns []dbmsreader.Column) {
	for _, v := range b.Column {
		columns = append(columns, &dbmsreader.BaseColumn{
			Name: v,
		})
	}
	return
}

// GetBaseTable - Retrieve table information.
func (b *BaseConfig) GetBaseTable() *database.BaseTable {
	return database.NewBaseTable(b.Connection.Table.Db, b.Connection.Table.Schema,
		b.Connection.Table.Name)
}

// GetWriteMode - Retrieve the write mode.
func (b *BaseConfig) GetWriteMode() string {
	if b.WriteMode == "" {
		return database.WriteModeInsert
	}
	return b.WriteMode
}

// GetBatchTimeout - Retrieve the batch timeout for a single write.
func (b *BaseConfig) GetBatchTimeout() time.Duration {
	if b.BatchTimeout.Duration == 0 {
		return defalutBatchTimeout
	}
	return b.BatchTimeout.Duration
}

// GetBatchSize - Retrieve the batch size for a single write.
func (b *BaseConfig) GetBatchSize() int {
	if b.BatchSize == 0 {
		return defalutBatchSize
	}

	return b.BatchSize
}

// GetPreSQL - Retrieve the prepared SQL statement.
func (b *BaseConfig) GetPreSQL() []string {
	return getSQlsWithoutEmpty(b.PreSQL)
}

// GetPostSQL - Retrieve the ending SQL statement.
func (b *BaseConfig) GetPostSQL() []string {
	return getSQlsWithoutEmpty(b.PostSQL)
}

// IgnoreOneByOneError - Ignore individual retry errors.
func (b *BaseConfig) IgnoreOneByOneError() bool {
	return b.ignoreOneByOneError
}

// GetRetryStrategy - Retrieve the retry strategy.
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
