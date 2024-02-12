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

// DBWrapper is a wrapper for database connection pools, used to reuse relevant database connection pools (from unit to instance: user)
type DBWrapper struct {
	*DB

	close sync.Once
}

// Open can acquire a reusable database connection pool wrapper through the database name and JSON configuration conf, similar to a smart pointer
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

// Close releases the database connection pool. If there are multiple references, the database connection pool will not be closed. When there are no references, it will be closed directly.
func (d *DBWrapper) Close() (err error) {
	d.close.Do(func() {
		if d.DB != nil {
			err = dbMap.Release(d.DB)
		}
	})

	return errors.Wrapf(err, "Release fail")
}
