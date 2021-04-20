package loader

import (
	"fmt"
	"sync"

	"github.com/Breeze0806/go-etl/libra/common/plugin"
)

var _centor = &center{
	dbStorages:     make(map[string]plugin.DBStorage),
	comparables:    make(map[string]plugin.RecordComparable),
	differStorages: make(map[string]plugin.DifferStorage),
	tableNameMaps:  make(map[string]plugin.TableNameMap),
	trackers:       make(map[string]plugin.Tracker),
}

//RegisterDBStorage 通过存储名name 注册存储storage，name重复或者storage为空时则会panic
func RegisterDBStorage(name string, storage plugin.DBStorage) {
	if err := _centor.registerDBStorage(name, storage); err != nil {
		panic(err)
	}
}

//LoadDBStorage 通过存储名name加载存储, name不存在则会panic
func LoadDBStorage(name string) plugin.DBStorage {
	storage, err := _centor.dbStorage(name)
	if err != nil {
		panic(err)
	}
	return storage
}

//RegisterRecordComparable 通过存储名name 注册存储comparable，name重复或者comparable为空时则会panic
func RegisterRecordComparable(name string, comparable plugin.RecordComparable) {
	if err := _centor.registerRecordComparable(name, comparable); err != nil {
		panic(err)
	}
}

//LoadRecordComparabale 通过存储名name加载comparable, name不存在则会panic
func LoadRecordComparabale(name string) plugin.RecordComparable {
	comparabale, err := _centor.recordComparabale(name)
	if err != nil {
		panic(err)
	}
	return comparabale
}

//RegisterTableNameMap 通过存储名name 注册存储表名映射，name重复或者comparable为空时则会panic
func RegisterTableNameMap(name string, comparable plugin.RecordComparable) {
	if err := _centor.registerRecordComparable(name, comparable); err != nil {
		panic(err)
	}
}

//LoadTableNameMap 通过存储名name加载存储表名映射, name不存在则会panic
func LoadTableNameMap(name string) plugin.TableNameMap {
	tableMap, err := _centor.tableNameMap(name)
	if err != nil {
		panic(err)
	}
	return tableMap
}

//RegisterTracker 通过存储名name 注册加载追踪器，name重复或者comparable为空时则会panic
func RegisterTracker(name string, tracker plugin.Tracker) {
	if err := _centor.registerTracker(name, tracker); err != nil {
		panic(err)
	}
}

//LoadTracker 通过存储名name加载追踪器, name不存在则会panic
func LoadTracker(name string) plugin.Tracker {
	tracker, err := _centor.tracker(name)
	if err != nil {
		panic(err)
	}
	return tracker
}

type center struct {
	dbStoragesMu sync.Mutex
	dbStorages   map[string]plugin.DBStorage

	comparablesMu sync.Mutex
	comparables   map[string]plugin.RecordComparable

	differStoragesMu sync.Mutex
	differStorages   map[string]plugin.DifferStorage

	tableNameMapsMu sync.Mutex
	tableNameMaps   map[string]plugin.TableNameMap

	trackersMu sync.Mutex
	trackers   map[string]plugin.Tracker
}

func (c *center) registerDBStorage(name string, storage plugin.DBStorage) error {
	if storage == nil {
		return fmt.Errorf("libra: storage(%v) is nil", name)
	}

	c.dbStoragesMu.Lock()
	defer c.dbStoragesMu.Unlock()
	if _, ok := c.dbStorages[name]; ok {
		return fmt.Errorf("libra: storage(%v) duplicates", name)
	}

	c.dbStorages[name] = storage
	return nil
}

func (c *center) dbStorage(name string) (plugin.DBStorage, error) {
	c.dbStoragesMu.Lock()
	defer c.dbStoragesMu.Unlock()
	storage, ok := c.dbStorages[name]
	if !ok {
		return nil, fmt.Errorf("libra: storage(%v) does not exist", name)
	}

	return storage, nil
}

func (c *center) registerRecordComparable(name string,
	comparable plugin.RecordComparable) error {
	if comparable == nil {
		return fmt.Errorf("libra: recordComparbale(%v) is nil", name)
	}

	c.comparablesMu.Lock()
	defer c.comparablesMu.Unlock()
	if _, ok := c.comparables[name]; ok {
		return fmt.Errorf("libra: recordComparable(%v) duplicates", name)
	}

	c.comparables[name] = comparable
	return nil
}

func (c *center) recordComparabale(name string) (plugin.RecordComparable, error) {
	c.comparablesMu.Lock()
	defer c.comparablesMu.Unlock()
	comparable, ok := c.comparables[name]
	if !ok {
		return nil, fmt.Errorf("libra: recordComparable(%v) does not exist", name)
	}

	return comparable, nil
}

func (c *center) registerDifferStorage(name string,
	differStorage plugin.DifferStorage) error {
	if differStorage == nil {
		return fmt.Errorf("libra: differStorage(%v) is nil", name)
	}

	c.differStoragesMu.Lock()
	defer c.differStoragesMu.Unlock()
	if _, ok := c.differStorages[name]; ok {
		return fmt.Errorf("libra: differStorage(%v) duplicates", name)
	}

	c.differStorages[name] = differStorage
	return nil
}

func (c *center) recordDifferStorage(name string) (plugin.DifferStorage, error) {
	c.differStoragesMu.Lock()
	defer c.differStoragesMu.Unlock()
	differStorage, ok := c.differStorages[name]
	if !ok {
		return nil, fmt.Errorf("libra: differStorage(%v) does not exist", name)
	}

	return differStorage, nil
}

func (c *center) registerTableNameMap(name string,
	tableNameMap plugin.TableNameMap) error {
	if tableNameMap == nil {
		return fmt.Errorf("libra: tableNameMap(%v) is nil", name)
	}

	c.tableNameMapsMu.Lock()
	defer c.tableNameMapsMu.Unlock()
	if _, ok := c.tableNameMaps[name]; ok {
		return fmt.Errorf("libra: tableNameMap(%v) duplicates", name)
	}

	c.tableNameMaps[name] = tableNameMap
	return nil
}

func (c *center) tableNameMap(name string) (plugin.TableNameMap, error) {
	c.tableNameMapsMu.Lock()
	defer c.tableNameMapsMu.Unlock()
	tableNameMap, ok := c.tableNameMaps[name]
	if !ok {
		return nil, fmt.Errorf("libra: tableNameMap(%v) does not exist", name)
	}

	return tableNameMap, nil
}

func (c *center) registerTracker(name string,
	tracker plugin.Tracker) error {
	if tracker == nil {
		return fmt.Errorf("libra: tracker(%v) is nil", name)
	}

	c.trackersMu.Lock()
	defer c.trackersMu.Unlock()
	if _, ok := c.trackers[name]; ok {
		return fmt.Errorf("libra: tracker(%v) duplicates", name)
	}

	c.trackers[name] = tracker
	return nil
}

func (c *center) tracker(name string) (plugin.Tracker, error) {
	c.trackersMu.Lock()
	defer c.trackersMu.Unlock()
	tracker, ok := c.trackers[name]
	if !ok {
		return nil, fmt.Errorf("libra: tracker(%v) does not exist", name)
	}

	return tracker, nil
}
