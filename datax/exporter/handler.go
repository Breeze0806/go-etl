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
	if err := json.Unmarshal([]byte(j.String()), jm); err != nil {
		log.Errorf("Unmarshal fail. err: %v, data: %v", err, j.String())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("Unmarshal fail. err: %v, data: %v", err, j.String()).Error()))
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
