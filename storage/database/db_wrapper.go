package database

import (
	"sync"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/schedule"
)

var dbMap = schedule.NewResourceMap()

type DBWrapper struct {
	*DB

	close sync.Once
}

func Open(name string, conf *config.Json) (dw *DBWrapper, err error) {
	var source Source
	if source, err = NewSource(name, conf); err != nil {
		return
	}
	loadOrNew := func() (r schedule.MappedResource, err error) {
		var db *DB
		db, err = NewDB(source)
		if err == nil {
			return db, nil
		}
		return
	}

	if r := schedule.NewLoadMappedResource(source.Key()); dbMap.UseCount(r) > 0 {
		loadOrNew = func() (schedule.MappedResource, error) {
			return r, nil
		}
	}

	var resource schedule.MappedResource

	if resource, err = dbMap.Get(loadOrNew); err != nil {
		return nil, err
	}
	return &DBWrapper{
		DB: resource.(*DB),
	}, nil
}

func (d *DBWrapper) Close() (err error) {
	d.close.Do(func() {
		if d.DB != nil {
			err = dbMap.Release(d.DB)
		}
	})

	return
}
