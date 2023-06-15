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

package plugin

import (
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go/encoding"
)

//DefaultJobCollector 默认工作收集器
type DefaultJobCollector struct{}

//NewDefaultJobCollector 创建默认工作收集器
func NewDefaultJobCollector() plugin.JobCollector {
	return &DefaultJobCollector{}
}

//MessageMap 空方法
func (d *DefaultJobCollector) MessageMap() *encoding.JSON {
	return nil
}

//MessageByKey 空方法
func (d *DefaultJobCollector) MessageByKey(key string) *encoding.JSON {
	return nil
}
