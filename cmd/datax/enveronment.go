package main

import (
	"context"
	"io/ioutil"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax"
)

type enveronment struct {
	config *config.JSON
	engine *datax.Engine
	err    error
	ctx    context.Context
	cancel context.CancelFunc
}

func newEnveronment(filename string) (e *enveronment) {
	e = &enveronment{}
	var buf []byte
	buf, e.err = ioutil.ReadFile(filename)
	if e.err != nil {
		return e
	}
	e.config, e.err = config.NewJSONFromBytes(buf)
	if e.err != nil {
		return e
	}
	e.ctx, e.cancel = context.WithCancel(context.Background())
	return e
}

func (e *enveronment) build() error {
	if e.err != nil {
		return e.err
	}
	e.engine = datax.NewEngine(e.ctx, e.config)

	e.err = e.engine.Start()
	return e.err
}

func (e *enveronment) close() {
	if e.cancel != nil {
		e.cancel()
	}
}
