package spi

import "github.com/Breeze0806/go-etl/libra/common/plugin"

type Task interface {
	JobID() int64
	TaskGroupID() int64
	TaskID() int64
	ExtraData() error
	Compare() error
	Post() error
}

type PluginMaker interface {
	MasterTable() plugin.MasterTable
	SlaveTable() plugin.SlaveTable
}
