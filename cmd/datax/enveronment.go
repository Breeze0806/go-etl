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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax"
	"github.com/Breeze0806/go-etl/datax/exporter"

	"github.com/fatih/color"
	"github.com/gorilla/handlers"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"golang.org/x/term"
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
		r.Handle("/metrics", exporter.NewHandler(e.engine))
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
		w, err := termWidth(os.Stdout)
		if err != nil {
			log.Errorf("termWidth fail. err: %v", err)
		}
		p := mpb.New(
			mpb.WithOutput(color.Output),
			mpb.WithWidth(w),
			mpb.WithAutoRefresh(),
		)

		barMap := make(map[string]*mpb.Bar)
		before := time.Now()
		barBuilder := func() mpb.BarFillerBuilder {
			s := mpb.SpinnerStyle("[    ]", "[   =]", "[  ==]", "[ ===]", "[====]", "[=== ]", "[==  ]", "[=   ]")
			return s.Meta(func(s string) string {
				return s
			})
		}

		barFunc := func(taskKey string, value int64) {
			var bar *mpb.Bar
			ok := false
			if bar, ok = barMap[taskKey]; !ok {
				bar = p.New(-1, barBuilder(),
					mpb.PrependDecorators(
						decor.Name(taskKey, decor.WC{C: decor.DindentRight | decor.DextraSpace}),
						decor.CurrentNoUnit("%v"),
					),
					mpb.AppendDecorators(
						decor.EwmaSpeed(nil, "% .1f/s", 30),
						decor.OnAbort(
							decor.EwmaETA(decor.ET_STYLE_GO, 0), "done\n",
						),
					))
				barMap[taskKey] = bar
			}

			bar.EwmaSetCurrent(value, time.Since(before))
		}

		barCurrentFunc := func(taskKey string, value int64) {
			var bar *mpb.Bar
			ok := false
			if bar, ok = barMap[taskKey]; !ok {
				bar = p.New(-1, barBuilder(),
					mpb.PrependDecorators(
						decor.Name(taskKey, decor.WC{C: decor.DindentRight | decor.DextraSpace}),
						decor.CurrentNoUnit("%v"),
					),
					mpb.AppendDecorators(
						decor.OnAbort(
							decor.EwmaETA(decor.ET_STYLE_GO, 0), "done\n",
						),
					))
				barMap[taskKey] = bar
			}
			bar.SetCurrent(value)
		}

		barByteFunc := func(taskKey string, value int64) {
			var bar *mpb.Bar
			ok := false
			if bar, ok = barMap[taskKey]; !ok {
				bar = p.New(-1, barBuilder(),
					mpb.PrependDecorators(
						decor.Name(taskKey, decor.WC{C: decor.DindentRight | decor.DextraSpace}),
						decor.Current(decor.SizeB1024(0), "% .2f"),
					),
					mpb.AppendDecorators(
						decor.EwmaSpeed(decor.SizeB1024(0), "% .2f", 30),
						decor.OnAbort(
							decor.EwmaETA(decor.ET_STYLE_GO, 0), "done\n",
						),
					))
				barMap[taskKey] = bar
			}
			bar.EwmaSetCurrent(value, time.Since(before))
		}

		barByteCurrentFunc := func(taskKey string, value int64) {
			var bar *mpb.Bar
			ok := false
			if bar, ok = barMap[taskKey]; !ok {
				bar = p.New(-1, barBuilder(),
					mpb.PrependDecorators(
						decor.Name(taskKey, decor.WC{C: decor.DindentRight | decor.DextraSpace}),
						decor.Current(decor.SizeB1024(0), "% .2f"),
					),
					mpb.AppendDecorators(
						decor.OnAbort(
							decor.EwmaETA(decor.ET_STYLE_GO, 0), "done\n",
						),
					))
				barMap[taskKey] = bar
			}
			bar.SetCurrent(value)
		}

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
				jm := &exporter.JobMetric{}
				j := e.engine.Metrics().JSON()
				if err := json.Unmarshal([]byte(j.String()), jm); err != nil {
					log.Errorf("Unmarshal fail. err: %v, data: %v", err, j.String())
					continue
				}

				for _, vi := range jm.Metrics {
					for _, vj := range vi.Metrics {
						taskName := fmt.Sprintf("job_id=%v,task_group_id=%v,task_id=%v",
							jm.JobID, vi.TaskGroupID, vj.TaskID)
						taskKey := fmt.Sprintf("datax_channel_total_byte(%v)", taskName)
						barByteFunc(taskKey, vj.Channel.TotalByte)
						taskKey = fmt.Sprintf("datax_channel_total_record(%v)", taskName)
						barFunc(taskKey, vj.Channel.TotalRecord)
						taskKey = fmt.Sprintf("datax_channel_byte(%v)", taskName)
						barByteCurrentFunc(taskKey, vj.Channel.Byte)
						taskKey = fmt.Sprintf("datax_channel_record(%v)", taskName)
						barCurrentFunc(taskKey, vj.Channel.Record)
					}
				}
				before = time.Now()
			}

			if exit {
				for _, bar := range barMap {
					bar.Abort(false)
				}
				p.Wait()
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

func termWidth(w io.Writer) (width int, err error) {
	if f, ok := w.(*os.File); ok {
		width, _, err = term.GetSize(int(f.Fd()))
		if err == nil {
			return width, nil
		}
	} else {
		err = errors.New("output is not a *os.File")
	}
	return 0, err
}
