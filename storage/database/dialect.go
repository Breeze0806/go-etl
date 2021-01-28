package database

import (
	"fmt"
	"sync"
)

//Dialect 数据库方言
type Dialect interface {
	Source(*BaseSource) (Source, error) //数据源
}

var dialects = &dialectMap{
	dialects: make(map[string]Dialect),
}

//RegisterDialect 注册数据库方言，当注册名称相同或者dialect为空时会panic
func RegisterDialect(name string, dialect Dialect) {
	if err := dialects.register(name, dialect); err != nil {
		panic(err)
	}
}

//UnregisterAllDialects 注销所有的数据库方言
func UnregisterAllDialects() {
	dialects.unregisterAll()
}

type dialectMap struct {
	sync.RWMutex
	dialects map[string]Dialect
}

func (d *dialectMap) register(name string, dialect Dialect) error {
	if dialect == nil {
		return fmt.Errorf("dialect %v is nil", name)
	}

	d.Lock()
	defer d.Unlock()
	if _, ok := d.dialects[name]; ok {
		return fmt.Errorf("dialect %v exists", name)
	}

	d.dialects[name] = dialect
	return nil
}

func (d *dialectMap) dialect(name string) (dialect Dialect, ok bool) {
	d.RLock()
	defer d.RUnlock()
	dialect, ok = d.dialects[name]
	return
}

func (d *dialectMap) unregisterAll() {
	d.Lock()
	defer d.Unlock()
	d.dialects = make(map[string]Dialect)
}
