package database

import (
	"sync"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/schedule"
)

var dbMap = schedule.NewResourceMap()

//DBWrapper 数据库连接池包装，用于复用相关的数据库连接池(单元到实例：user)
type DBWrapper struct {
	*DB

	close sync.Once
}

//Open 通过数据库name和json配置conf 获取可以复用的数据库连接池包装,类似智能指针
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
		return nil, err
	}
	return &DBWrapper{
		DB: resource.(*DB),
	}, nil
}

//Close 释放数据库连接池，如果有多个引用，则不会关闭该数据库连接池，没有引用时就直接关闭
func (d *DBWrapper) Close() (err error) {
	d.close.Do(func() {
		if d.DB != nil {
			err = dbMap.Release(d.DB)
		}
	})

	return
}
