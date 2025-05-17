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

package exporter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Breeze0806/go-etl/datax"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	engine *datax.Engine
}

func NewHandler(engine *datax.Engine) *Handler {
	return &Handler{
		engine: engine,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jm := &JobMetric{}
	j := h.engine.Metrics().JSON()
	var err error
	if err = json.Unmarshal([]byte(j.String()), jm); err != nil {
		writeResult(w, fmt.Errorf("MarshalIndent fail. err: %v, data: %v", err, *jm))
		return
	}

	if r.URL.Query().Get("t") == "json" {
		var data []byte
		if data, err = json.MarshalIndent(jm, "", "    "); err != nil {
			writeResult(w, fmt.Errorf("MarshalIndent fail. err: %v, data: %v", err, *jm))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}

	jsonMetricCollector := JSONMetricCollector{
		Metric: jm,
	}

	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(jsonMetricCollector)
	ph := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	ph.ServeHTTP(w, r)
}

func writeResult(w http.ResponseWriter, err error) {
	log.Errorf("%v", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
