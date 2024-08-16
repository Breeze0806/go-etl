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

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax"
	"github.com/gorilla/handlers"
)

type enveronment struct {
	config *config.JSON
	engine *datax.Engine
	err    error
	ctx    context.Context
	cancel context.CancelFunc
	server *http.Server
	addr   string
}

func newEnveronment(filename string, addr string) (e *enveronment) {
	e = &enveronment{}
	var buf []byte
	buf, e.err = os.ReadFile(filename)
	if e.err != nil {
		return e
	}
	e.config, e.err = config.NewJSONFromBytes(buf)
	if e.err != nil {
		return e
	}
	e.ctx, e.cancel = context.WithCancel(context.Background())
	e.addr = addr
	return e
}

func (e *enveronment) build() error {
	return e.initEngine().initServer().startEngine().err
}

func (e *enveronment) initEngine() *enveronment {
	if e.err != nil {
		return e
	}
	e.engine = datax.NewEngine(e.ctx, e.config)

	return e
}

func (e *enveronment) initServer() *enveronment {
	if e.err != nil {
		return e
	}
	if e.addr != "" {
		r := http.NewServeMux()
		recoverHandler := handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))
		r.Handle("/metrics", newMetricHandler(e.engine))
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
		e.server = &http.Server{
			Addr:    e.addr,
			Handler: recoverHandler(r),
		}
		go func() {
			log.Debugf("listen begin: %v", e.addr)
			defer log.Debugf("listen end: %v", e.addr)
			if err := e.server.ListenAndServe(); err != nil {
				log.Errorf("ListenAndServe fail. addr: %v err: %v", e.addr, err)
			}
			log.Infof("ListenAndServe success. addr: %v", e.addr)
		}()
	}

	return e
}

func (e *enveronment) startEngine() *enveronment {
	if e.err != nil {
		return e
	}
	go func() {
		statsTimer := time.NewTicker(1 * time.Second)
		defer statsTimer.Stop()
		exit := false
		for {
			select {
			case <-statsTimer.C:
			case <-e.ctx.Done():
				exit = true
			}
			if e.engine.Container != nil {
				fmt.Printf("%v\r", e.engine.Metrics().JSON())
			}

			if exit {
				return
			}
		}
	}()
	e.err = e.engine.Start()

	return e
}

func (e *enveronment) close() {
	if e.server != nil {
		e.server.Shutdown(e.ctx)
	}

	if e.cancel != nil {
		e.cancel()
	}
}
