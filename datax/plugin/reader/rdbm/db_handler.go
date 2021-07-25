package rdbm

import (
	"database/sql"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/storage/database"
)

type DbHandler interface {
	Querier(name string, conf *config.JSON) (Querier, error)
	Config(conf *config.JSON) (Config, error)
	TableParam(config Config, querier Querier) database.Parameter
}

type BaseDbHandler struct {
	newQuerier func(name string, conf *config.JSON) (Querier, error)
	opts       *sql.TxOptions
}

func NewBaseDbHandler(newQuerier func(name string, conf *config.JSON) (Querier, error), opts *sql.TxOptions) *BaseDbHandler {
	return &BaseDbHandler{
		newQuerier: newQuerier,
		opts:       opts,
	}
}

func (d *BaseDbHandler) Querier(name string, conf *config.JSON) (Querier, error) {
	return d.newQuerier(name, conf)
}

func (d *BaseDbHandler) Config(conf *config.JSON) (Config, error) {
	return NewBaseConfig(conf)
}

func (d *BaseDbHandler) TableParam(config Config, querier Querier) database.Parameter {
	return NewTableParam(config, querier, d.opts)
}
