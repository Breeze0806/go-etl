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

package file

import (
	"encoding/json"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go/time2"
)

var (
	defalutBatchSize    = 1000
	defalutBatchTimeout = 1 * time.Second
)

// Config 文件流输出配置
type Config interface {
	GetBatchSize() int              //单次批量写入数
	GetBatchTimeout() time.Duration //单次批量写入超时时间
}

// BaseConfig 基础文件流输出配置
type BaseConfig struct {
	BatchSize    int            `json:"batchSize"`    //单次批量写入数
	BatchTimeout time2.Duration `json:"batchTimeout"` //单次批量写入超时时间
}

// NewBaseConfig 通过json配置获取基础文件流输出配置
func NewBaseConfig(conf *config.JSON) (*BaseConfig, error) {
	c := &BaseConfig{}
	if err := json.Unmarshal([]byte(conf.String()), c); err != nil {
		return nil, err
	}
	return c, nil
}

// GetBatchTimeout 单次批量写入超时时间
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
