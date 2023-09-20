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

package database

import (
	"sync"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/schedule"
	"github.com/pingcap/errors"
)

var dbMap = schedule.NewResourceMap()

// DBWrapper 数据库连接池包装，用于复用相关的数据库连接池(单元到实例：user)
type DBWrapper struct {
	*DB

	close sync.Once
}

// Open 通过数据库name和json配置conf 获取可以复用的数据库连接池包装,类似智能指针
func Open(name string, conf *config.JSON) (dw *DBWrapper, err error) {
	var source Source
	if source, err = NewSource(name, conf); err != nil {
		return
	}
	create := func() (r schedule.MappedResource, err error) {
		var db *DB
		db, err = NewDB(source)
		if err == nil {
			return db, nil
		}
		return
	}

	var resource schedule.MappedResource

	if resource, err = dbMap.Get(source.Key(), create); err != nil {
		return nil, errors.Wrapf(err, "Get fail")
	}
	return &DBWrapper{
		DB: resource.(*DB),
	}, nil
}

// Close 释放数据库连接池，如果有多个引用，则不会关闭该数据库连接池，没有引用时就直接关闭
func (d *DBWrapper) Close() (err error) {
	d.close.Do(func() {
		if d.DB != nil {
			err = dbMap.Release(d.DB)
		}
	})

	return errors.Wrapf(err, "Release fail")
}
