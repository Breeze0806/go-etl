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

import "github.com/Breeze0806/go/encoding"

// JobCollector 工作信息采集器，用于统计整个工作的进度，错误信息等
// toto 当前未实现监控模块，为此需要在后面来实现这个接口的结构体
type JobCollector interface {
	JSON() *encoding.JSON
	JSONByKey(key string) *encoding.JSON
}
