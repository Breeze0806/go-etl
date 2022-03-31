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
	"github.com/Breeze0806/go-etl/element"
)

//RecordSender 记录发送器
type RecordSender interface {
	CreateRecord() (element.Record, error)  //创建记录
	SendWriter(record element.Record) error //将记录发往写入器
	Flush() error                           //将记录刷新到记录发送器
	Terminate() error                       //终止发送信号
	Shutdown() error                        //关闭
}
