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
	"context"

	"github.com/Breeze0806/go-etl/storage/database"
)

//OffsetTracker 位移追踪器
type OffsetTracker interface {
	Read(ctx context.Context) error
	Offset() Offset
	SetOffset(offset Offset)
	Write(ctx context.Context) error
	Close() error
}

//PageParamTracker 页查询参数追踪器
type PageParamTracker interface {
	Read(ctx context.Context) ([]PageParam, error)
	Add(pageParam PageParam)
	Delete(pageParam PageParam)
	Write(ctx context.Context) error
	Close() error
}

//Tracker 追踪器
type Tracker interface {
	OffsetTracker(master database.Table) (OffsetTracker, error)
	PageParamTracker(master database.Table) (PageParamTracker, error)
}
